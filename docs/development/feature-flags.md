---
title: Feature Flags
{{ .EditURL }}
---

Feature flags allow developers to roll out new features or workflows safely.

For sake of ease, the examples throughout this document use the amiibo code developed within the [Feature documentation](/docs/development/feature).
Be sure to read through that if you're unfamiliar with how features are developed within Clutch.

## Gating a Workflow

### Backend

The feature flag module on the backend powers the storing and serving of feature flag configuration for the frontend.

At the moment the module only supports simple flags with boolean types.

```yaml  title="backend/clutch-config.yaml"
modules:
  - name: clutch.module.assets
  - name: clutch.module.healthcheck
  ...
  // highlight-start
  - name: clutch.module.featureflag
    typed_config:
      "@type": types.google.com/clutch.config.module.featureflag.v1.Config
      simple:
        flags:
          amiiboLookupEnabled: false
  // highlight-end
```

In the example above the feature flag configuration has a simple flag `amiiboLookupEnabled` with a value of `false`. When adding new simple feature flags, generally speaking, they should start with a value of false to allow the deployment of new code without having it be enabled in production. Now that the server is configured to serve our new feature flag we can consume it on the frontend.

### Frontend

All workflow routes can be gated behind a feature flag via their route configuration.
If a route depends on a simple feature flag and it is either not found or set to false the route and component will not be registered on the application.

:::info
Workflow route feature flags must be set as simple feature flags on the server with a boolean type.
:::

The frontend will fetch feature flags on the initial rendering of Clutch as well as at a determined interval. This interval
defaults to 5 minutes but can be overridden by setting the `REACT_APP_FF_POLL` environment variable. The polling behavior
ensures that even if users do not reload the application there is a max TTL of user's local feature flag values of N, where N is the
afformentioned polling interval.

```tsx title="frontend/workflows/amiibo/src/index.tsx"
const register = (): WorkflowConfiguration => {
  return {
    developer: {
      name: "Name McName",
      contactUrl: "mailto:foo@foo-email.com",
    },
    path: "amiibo",
    group: "Amiibo",
    displayName: "Amiibo",
    routes: {
      landing: {
        path: "/lookup",
        description: "Lookup all Amiibo by name.",
        component: Amiibo,
        // highlight-next-line
        featureFlag: "amiiboLookupEnabled",
      },
    },
  };
};
```

### Enabling the Workflow

Once the above backend and frontend code has been deployed the configuration file can be updated, potentially in only a single environment, to enable the new flag.

```yaml title="backend/clutch-config.yaml"
- name: clutch.module.featureflag
    typed_config:
      "@type": types.google.com/clutch.config.module.featureflag.v1.Config
      simple:
        flags:
          // highlight-next-line
          amiiboLookupEnabled: true
```

Above the `amiiboLookupEnabled` flag has been set to `true`. Once that is deployed the `/amiibo/lookup` route will now be accessible.


## Gating Components

### Backend

The backend setup in this use case is similar to that of gating a workflow but the flag name should be specific to the feature being gated.

```yaml title="backend/clutch-config.yaml"
- name: clutch.module.featureflag
    typed_config:
      "@type": types.google.com/clutch.config.module.featureflag.v1.Config
      simple:
        flags:
          // highlight-next-line
          amiiboImageEnabled: true
```

### Frontend

Individual components can also be gated behind feature flags by wrapping them in the `SimpleFeatureFlag` component.

```tsx title="frontend/workflows/amiibo/src/hello-world.tsx"
const WorkflowStep: React.FC<WizardChild> = () => {
  const amiiboData = useDataLayout("amiiboData");
  let amiiboResults = amiiboData.displayValue();
  if (_.isEmpty(amiiboResults)) {
    amiiboResults = [];
  }

  return (
    <WizardStep error={amiiboData.error} isLoading={amiiboData.isLoading}>
      <SimpleFeatureFlag feature="amiiboImageEnabled">
        <FeatureOff>
          <Table headings={["Name", "Series", "Type"]}>
            {amiiboResults.map((amiibo, index: number) => (
              <TableRow key={index}>
                {amiibo.name}
                {amiibo.amiiboSeries}
                {amiibo.type}
              </TableRow>
            ))}
          </Table>
        </FeatureOff>
        <FeatureOn>
          <Table headings={["Name", "Image", "Series", "Type"]}>
            {amiiboResults.map((amiibo, index: number) => (
              <TableRow key={index}>
                {amiibo.name}
                <img src={amiibo.imageUrl} height="75px"/>
                {amiibo.amiiboSeries}
                {amiibo.type}
              </TableRow>
            ))}
          </Table>
        </FeatureOn>
      </SimpleFeatureFlag>
    </WizardStep>
  );
};
```

In the above example a new column containing the amiibo image has been added, gated behind the `amiiboImageEnabled` flag. If this flag is `true`
users will be shown amiibo images, otherwise they won't be rendered on the page at all.

While the above example demonstrates using both the `FeatureOff` and `FeatureOn` components these do not need to be used in unison and can be used to add new components or remove old components accordingly.