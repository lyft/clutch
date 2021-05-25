---
title: Overview
{{ .EditURL }}
---

The frontend code in Clutch is separated into two distinct pieces; the packages that power the platform and the out of the box workflows that utilize these packages.

The goal for developing the core Clutch code was to make it as easy as possible for those with little to no domain knowledge or experience. To do this we leverage [Lerna](https://github.com/lerna/lerna) and [Yarn workspaces](https://classic.yarnpkg.com/en/docs/workspaces/) to run commands against all core packages and installed workflows.

## Structure

```
frontend
├─ api             # generated code from proto
├─ packages
│  ├─ app          # main entry point, workflow and layout manager
│  ├─ core         # re-useable UI components
│  ├─ data-layout  # state management library
│  ├─ tools        # reusable lint configuration and scripts
│  └─ wizard       # step-by-step flow and form management
└─ workflows       # core workflows
```

## Packages

The packages within Clutch are separated by functionality, with the intent being that each package provides value to workflows while also not being required.

An interactive playground of our component library is available at [storybook.clutch.sh](https://storybook.clutch.sh).

### @clutch-sh/core

The Core package consists of various reusable components that are shared between workflows and/or other Clutch packages. This can span from things like the application provider, which renders the Clutch app, and resolver component ([see above](./###Resolver)) to something as simple as a centralized button component.

It’s important to note that all components which reside in this package should be developed in a reusable manner and should have a wide application. Workflow-specific code should not be contributed to the core package.

### @clutch-sh/data-layout

The data-layout package provides both layouts for workflows to store data within and a manager for those layouts.

Data layouts follow a pattern outlined below.

1. Define a layout object, assiging each layout to a key. These keys will be how you reference the layouts throughout the workflow.

```jsx
const dataLayout = {
    resourceData: {},
    terminationData: {
      deps: ["resourceData"],
      hydrator: resourceData => {
        return client
          .post("/v1/aws/ec2/terminateInstance", {
            instance_id: resourceData.instanceId,
            region: resourceData.region,
          })
          .then(resp => {
            if (resp.data?.entries?.[0].status !== null) {
              const { status } = resp.data.entries[0];
              throw new Error(status);
            }
            return resp;
          });
      },
    },
  };
```

Layouts can be completely empty as showing with `resourceData` above but they can also override the default values for a layout to implement custom behavior.

In the example above `terminationData` defines values for both `deps` and `hydrator`. Let's dig into these keys more.

  - `deps` is powerful in that it allows layouts to specify their dependency on other layouts. If a layout is accessed before its dependencies have data, an error will be displayed to the user. This is a failsafe way to ensure that users provide required data. In the example above, `terminationData` has a dependency on `resourceData` containing some data.

  - `hydrator` is a callback fired when a layout is first used by a workflow with all of the defined `deps` data as arguments. This is particularly useful when data is needed from the user before dispatching a request. In the example above the `terminationData` layout posts to `/v1/aws/ec2/terminateInstance` with data from it's dependency `resourceData` and analyzes the response for possible errors before returning the result.

2. Use data layout(s) in a workflow.

```jsx
const terminationData = useDataLayout("terminationData");
```

This allows the workflow to access any data the layout contains from either a prior step or as a result of the hydrator. As mentioned in the prior step, if this is the first time the `terminationData` layout has been used, this will invoke the hydrator callback.

### @clutch-sh/tools

The Tools package is meant to be used exclusively in development and contains dependencies and configuration files for tasks such as linting, testing, etc.

### @clutch-sh/wizard

The Wizard package works hand-in-hand with the data-layout package, leveraging it for workflow progression and state while providing its steps with the ability to easily access data layouts.

There are two pieces to effectively utilizing the wizard component, a `Wizard` and a `WizardStep` components where the Wizard has one or more WizardStep children.

```jsx
const Step1 = () => {
  const userData = useDataLayout("userData");

  return (
    <WizardStep error={userData.error} isLoading={userData.isLoading}>
      <div>Hello, {userData.value.firstName}!</div>
    </WizardStep>
  );
};

const Workflow = () => (
  <Wizard>
    <Step1>
  </Wizard>
);
```

In the above example you can see that the `WizardStep` utilizes a data layout's `error` and `isLoading` props to display helpful feedback to the user.

To move between wizard steps the `@clutch-sh/core` package exposes hooks for moving forward and backward. This allows workflows to take additional action before submitting a step or disable backward movement altogether.

```jsx
import { useWizardContext } from "@clutch-sh/core";

...

const { onSubmit, onBack } = useWizardContext();
```

## Workflows

Workflows are essentially a combination of the packages listed above with configration values for routing and display purposes. Bear in mind that while, at the time of writing this, the out of the box workflows are all linear and serve a single purpose not all workflows need to follow that design. A workflow could be anything from a dashboard of information to a single form for submission.

The most important piece of a workflow, even before the implementation is the configuration. A workflows configuration allows it to be registered on the Clutch app, allows for users to override default values to suit their specific needs, and provides ownership details for others.

```jsx
import RemoteTriage from "./remote-triage";

const register = () => {
  return {
    id: "clutch-envoy",
    developer: {
      name: "Lyft",
      contactUrl: "mailto:hello@clutch.sh",
    },
    path: "envoy",
    group: "Envoy",
    displayName: "Envoy",
    routes: {
      remoteTriage: {
        path: "triage",
        component: RemoteTriage,
        displayName: "Remote Triage",
        description: "Triage Envoy configurations.",
        trending: false,
      },
    },
  };
};

export default register;
```

Let's break down the example above.

Every workflow package has a single exported callback function that is used by Clutch to register it. This callback function must return an object with the following:

  - `developer`: the name and contact url for the party developing this workflow package. This is used if a workflow fails to load with an error message giving users the proper channel to report failures to.
  - `path`: the top level routing path for this collection of workflows.
  - `group`: the group to nest these workflows under within the navigation drawer.
  - `displayName`: The display name for these workflows. This is used throughout the Clutch application in various spots.
  - `routes`: A collection of the workflows and their metadata contained within this package.

Each route is a key/value pair where the key is used by others to register this workflow on their Clutch application. The keys of a route are similar to a package with a few new ones:

  - `path`: the path of the component relative to the top level routing path.
  - `component`: the component to render for this route.
  - `description`: a short description of what this component does.
  - `trending`: this controls if a workflow is shown on the landing page of Clutch. It is generally advised to omit this and override it's value in the application config.
  - `props`: this is a list of keys that users must specify in `clutch.config.js` when adding this workflow. If these are not present in the config the workflow will not be loaded.

### Using the Resolver Component

The purpose and structure of the Resolver component is explained in the [Architecture & Components section](/docs/about/architecture#resolvers), if you haven’t had a chance to read that yet this would be a good chance to as its content will be helpful in understanding the usage.

Using the Resolver component in your workflow should be mostly seamless requiring only a few properties that are specific to your workflow. To help with our explanation let’s look at the following example:

```jsx
const InstanceIdentifier = ({ resolverType }) => {
  const { onSubmit } = useWizardContext();
  const resolvedResourceData = useDataLayout("resourceData");

  const onResolve = ({ results }) => {
    // Decide how to process results.
   // Thanks to the defined search limit we know the results will be limited to 1.
    resolvedResourceData.assign(results[0]);
    onSubmit();
  };

  return <Resolver type={resolverType} searchLimit={1} onResolve={onResolve} />;
};
```

Notice that the resolver specified only three properties: a type, search limit, and a callback on resolution.

  - `type`: this property denotes the type of resource the resolver will prompt the user for and pass back to the onResolve callback. For example, `clutch.aws.ec2.v1.Instance`. In the above example this value is being passed in as a prop from its parent which has it injected from the [application configuration](/docs/configuration).
  - `searchLimit`: allows developers to specify a threshold for the number of resources to search for. This is particularly useful when searching across multiple environments for a small number of resources as the server will respond as soon as this threshold is satisfied. This also denotes the expected number of environments to search, allowing the resolver component to warn when searches fail in specific environments.
  - `onResolve`: when the resolver completes a lookup and casts the response to a specified protobuf object it will need a way to pass the results back to a workflow. This property allows users to specify a callback that takes the results and use them in the appropriate way for that workflow.
