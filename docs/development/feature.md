---
title: Feature Development
{{ .EditURL }}
---

import useBaseUrl from '@docusaurus/useBaseUrl';
import Image from '@site/src/components/Image';

This is a step-by-step walkthrough of developing a feature from scratch, taking advantage of the most commonly used Clutch capabilities across the frontend and backend. Familiarity with concepts from the [Architecture Reference](/docs/about/architecture) will help, so it is recommended to skim there first, or refer back to it throughout this tutorial.

⏲️ **Estimated time to complete:** ~30 minutes

:::info
This guide should work for a [custom gateway](/docs/development/custom-gateway) or Clutch core as the directory structure is the same for both. In general, starting with a custom gateway is a good idea and recommended for this guide.
:::

If you get lost at any point all of the files referenced in this tutorial can be found in the [examples](https://github.com/lyft/clutch/tree/main/examples/amiibo).

## The New Feature

The goal of this tutorial is to build a visual search on top of [AmiiboAPI](https://www.amiiboapi.com).

Users will start by searching for the name of the item they are interested in. Users can search with a simple text query.

The search will return a list of results, after which the user can pick a single result to see full details on.

## Design and Implementation

### 1. Draw a Wireframe

Start by creating a low fidelity wireframe of the new feature. This helps to gather requirements and feedback from stakeholders as well as focus the implementation.

Get out a pen and paper or use a drawing tool. Popular options include [draw.io](https://app.diagrams.net/) and [Google Drawings](https://docs.google.com/drawings/).

<img alt="Feature Wireframe" src={useBaseUrl('img/docs/feature-development/feature-mockup.png')} width="50%" />

From the diagram, it's clear that there is only one API endpoint to define. The user provides input on the first step and sees a list of results on the second step.

### 2. API: Define the Schema

Now it's time to determine how information will be presented to the frontend from the backend.

In this example, it would be easy enough to call the API directly from the frontend. However, in order to illustrate Clutch concepts, all calls will be proxied through the Clutch backend. This would normally provide authentication, authorization, auditing, logging, stats, additional input validation, API token handling, etc, none of which are strictly needed for this example. Following the guide should reinforce Clutch concepts albeit with a toy example.

For the reasons stated above, the API defined below will largely look identical to the provider's API.

#### Boilerplate

Clutch uses protobuf for interface definitions. An empty definition is provided below.

```protobuf title="api/amiibo/v1/amiibo.proto"
syntax = "proto3";

package clutch.amiibo.v1;

option go_package = "amiibov1";

import "google/api/annotations.proto";

service AmiiboAPI {}
```

- The path of the API file should correspond to the path of the API itself. e.g. `/v1/amiibo` maps to `api/amiibo/v1/`.
- The header of the file will always starts with the `syntax` specifier. Clutch always uses `proto3`.
- The `package` directive will also map directly to the filename.
- `go_package` option make importing multiple packages a little less confusing by giving a default import name other than `v1`.
- `service` is the group of APIs that share the same underlying `struct` and defines the gRPC service.

#### Endpoints

Next, add the endpoint definitions, including the [Google API annotation](https://github.com/googleapis/googleapis/blob/d7c66c92df10a9822fa1380d88b73286651b4f9f/google/api/http.proto#L46). Added lines are highlighted below.

```protobuf title="api/amiibo/v1/amiibo.proto"
syntax = "proto3";

package clutch.amiibo.v1;

option go_package = "amiibov1";

import "google/api/annotations.proto";

service AmiiboAPI {
// highlight-start
  rpc GetAmiibo (GetAmiiboRequest) returns (GetAmiiboResponse) {
    option (google.api.http) = {
      post : "/v1/amiibo/getAmiibo"
      body : "*"
    };
  }
// highlight-end
}

// highlight-next-line
message GetAmiiboRequest {}

// highlight-next-line
message GetAmiiboResponse {}
```

- Clutch always uses `post` for the method to simplify the mapping between gRPC and HTTP endpoint names.

#### Request and Response

Now fill in the request and response. [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate) message options are used for input validation. Added lines are highlighted below.

```protobuf title="api/amiibo/v1/amiibo.proto"
syntax = "proto3";

package clutch.amiibo.v1;

option go_package = "amiibov1";

import "google/api/annotations.proto";
// highlight-next-line
import "validate/validate.proto";

service AmiiboAPI {
  rpc GetAmiibo (GetAmiiboRequest) returns (GetAmiiboResponse) {
    option (google.api.http) = {
      post : "/v1/amiibo/getAmiibo"
      body : "*"
    };
  }
}

message GetAmiiboRequest {
  // highlight-next-line
  string name = 1 [ (validate.rules).string = {min_bytes : 1} ];
}

message GetAmiiboResponse {
  // highlight-next-line
  repeated Amiibo amiibo = 1;
}

// highlight-start
message Amiibo {
  string character = 1;
  string name = 2;
  string amiibo_series = 3;
  string image_url = 4;

  enum Type {
    UNSPECIFIED = 0;
    CARD = 1;
    FIGURE = 2;
    YARN = 3;
  }
  Type type = 5;
}
// highlight-end
```

- Only use input validation on input objects. Input validation is superfluous on response objects.

#### Generate the Code

With the API definition complete, the code can be generated for the frontend and backend. From the Clutch root, run:

```bash
make api
```

The resulting generated code will be in:
- `frontend/api/src/index.js`, the unified bundle of all frontend generated objects.
- `backend/api/amiibo/v1/`, the `amiibov1` package.

### 3. Backend: Implement the APIs

On the backend, a `service` is needed to interact with the API and a `module` is needed to service the API endpoints defined in protobuf.

The module is the implementation of the API endpoints defined in the protobuf.

:::caution On Complexity
Services are reuseable clients for external data or services. In that way, this example is a bit contrived for the purposes of illustrating Clutch concepts.

**It is acceptable to call the third-party service directly from the module to simplify the code** if there is no need for substitution via interface or there are not multiple modules calling into a service. Refactoring from a combined module is not that hard if the need does arise in the future.
:::

#### Service

The service code is fairly straightforward, especially if you're familiar with Go. Of note are the `New` method and the `Name`. The signature of `New` is the same for every service, and makes the service usable in the Clutch configuration.

Also note that the `Client` interface freely makes use of the generated objects from the API.

```go title="backend/service/amiibo/amiibo.go"
package amiibo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	amiibov1 "github.com/lyft/clutch/backend/api/amiibo/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.amiibo"

const apiHost = "https://www.amiiboapi.com"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	return &client{http: &http.Client{}}, nil
}

type Client interface {
	GetAmiibo(ctx context.Context, name string) ([]*amiibov1.Amiibo, error)
}

type client struct {
	http *http.Client
}

type RawResponse struct {
	Amiibo []*RawAmiibo `json:"amiibo"`
}

type RawAmiibo struct {
	Character    string `json:"character"`
	AmiiboSeries string `json:"amiiboSeries"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Type         string `json:"type"`
}

func (r RawAmiibo) toProto() *amiibov1.Amiibo {
	t := strings.ToUpper(r.Type)
	return &amiibov1.Amiibo{
		Name:         r.Name,
		AmiiboSeries: r.AmiiboSeries,
		ImageUrl:     r.Image,
		Character:    r.Character,
		Type:         amiibov1.Amiibo_Type(amiibov1.Amiibo_Type_value[t]),
	}
}

func charactersFromJSON(data []byte) ([]*amiibov1.Amiibo, error) {
	raw := &RawResponse{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	ret := make([]*amiibov1.Amiibo, len(raw.Amiibo))
	for i, a := range raw.Amiibo {
		ret[i] = a.toProto()
	}
	return ret, nil
}

func (c *client) GetAmiibo(ctx context.Context, name string) ([]*amiibov1.Amiibo, error) {
	url := fmt.Sprintf("%s/api/amiibo?character=%s", apiHost, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, status.Error(service.CodeFromHTTPStatus(resp.StatusCode), string(body))
	}
	return charactersFromJSON(body)
}
```

#### Module
The module is instantiated after the service when the gateway starts. This allows the module to fetch a reference to any services it needs, which are then stored on the module instance for use during requests.

```go title="backend/module/amiibo/amiibo.go"
package amiibo

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	amiibov1 "github.com/lyft/clutch/backend/api/amiibo/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	amiiboservice "github.com/lyft/clutch/backend/service/amiibo"
)

const Name = "clutch.module.amiibo"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	svc, ok := service.Registry["clutch.service.amiibo"]
	if !ok {
		return nil, errors.New("no amiibo service was registered")
	}

	client, ok := svc.(amiiboservice.Client)
	if !ok {
		return nil, errors.New("amiibo service in registry was the wrong type")
	}

	return &mod{client: client}, nil
}

type mod struct {
	client amiiboservice.Client
}

func (m *mod) Register(r module.Registrar) error {
	amiibov1.RegisterAmiiboAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(amiibov1.RegisterAmiiboAPIHandler)
}

func (m *mod) GetAmiibo(ctx context.Context, request *amiibov1.GetAmiiboRequest) (*amiibov1.GetAmiiboResponse, error) {
	a, err := m.client.GetAmiibo(ctx, request.Name)
	if err != nil {
		return nil, err
	}
	return &amiibov1.GetAmiiboResponse{Amiibo: a}, nil
}
```

#### Gateway Registration

##### Components

The module and service must be in the component list at compile time in order for the binary to be able to read instantiate them from config based on their respective names.

```go title="backend/main.go"
package main

import (
    "github.com/lyft/clutch/backend/cmd/assets"
    "github.com/lyft/clutch/backend/gateway"
    // highlight-start
    amiibomod "github.com/lyft/clutch/backend/module/amiibo"
    amiiboservice "github.com/lyft/clutch/backend/service/amiibo"
    // highlight-end
)

func main() {
	flags := gateway.ParseFlags()
	components := gateway.CoreComponentFactory
	// highlight-start
    components.Modules[amiibomod.Name] = amiibomod.New
    components.Services[amiiboservice.Name] = amiiboservice.New
    // highlight-end

	gateway.Run(flags, components, assets.VirtualFS)
}
```

##### Config
Add the components to the config.
```yaml title="backend/clutch-config.yaml"
...
services:
  ...
  // highlight-next-line
  - name: clutch.service.amiibo
modules:
  ...
  // highlight-next-line
  - name: clutch.module.amiibo
...
```

##### Test!

Run the gateway with the new components:
```bash
make backend-dev
```

With the gateway still running, try the following in a separate shell:
```bash
curl -X POST localhost:8080/v1/amiibo/getAmiibo -d '{"name": "peach"}'
```

A number of results should appear in the terminal. Also try passing an empty name to see how input validation is automatically applied from the protobuf annotation of the field.

### 4. Frontend: User Interface

The frontend consists of a few pieces that work together to display your workflow to users. Workflows are exposed via a registration function containing default configuration values and workflow specific properties. One of these properties is the component that should be rendered for this workflow. Once we have that component we need to register the new workflow on the Clutch application.

#### Scaffolding

To simplify the creation of workflows you can run a scaffolding tool. It will prompt you for some information and produce a new workflow with some templates.

To get started run the generator and provide the details about our Amiibo workflow.
```
> make scaffold-workflow
*** Is the destination okay:
> clutch/frontend/workflows
Is this okay? [Y/n]: Y
Enter the name of this workflow [Hello World]: Amiibo
Enter a description of the workflow [Greet the world]: Lookup all Amiibo by name
Enter the developer's name [dschaller]: Lyft
Enter the developer's email [derek@lyft.com]: hello@example.com

*** Generating...
*** All done!
```

#### Building out the Workflow

You should now have a new Amiibo workflow in the destination directory outlined in the CLI. However, this scaffolding left behind some default values that we should update.

##### Component
Let's update the component.

Define a functional component for the amiibo lookup by adding the highlighted lines.

```tsx title="frontend/workflows/amiibo/src/hello-world.tsx"
// highlight-start
import React, { ChangeEvent } from "react";
import {
  Button,
  ButtonGroup,
  useWizardContext,
  TextField,
} from "@clutch-sh/core";

import { useDataLayout } from "@clutch-sh/data-layout";
// highlight-end
import { Wizard, WizardStep } from "@clutch-sh/wizard";
...
// highlight-start
const AmiiboLookup: React.FC<WizardChild> = () => {
  const { onSubmit } = useWizardContext();
  const userInput = useDataLayout("userInput");

  const onChange = ((event: ChangeEvent<{value: string}>) => {
    userInput.assign({name: event.target.value});
  });

  return (
    <>
      <TextField onChange={onChange} onReturn={onSubmit}/>
      <ButtonGroup>
        <Button text="Search" onClick={onSubmit}/>
      </ButtonGroup>
    </>
  );
};
// highlight-end
```

This will present a text field which updates a data layout called `userInput` on changes and a button for users to click when they are ready to search.

<Image alt="Amiibo Lookup Panel" src={useBaseUrl('img/docs/feature-development/lookup-panel.png ')} width="75%" variant="centered"/>

Now let's build a way to display the details panel.

```tsx title="frontend/workflows/amiibo/src/hello-world.tsx"
...
import React, { ChangeEvent } from "react";
// highlight-next-line
import _ from "lodash";
import {
	...
	useWizardContext,
	// highlight-next-line
	Table,
	// highlight-next-line
	TableRow,
	TextField,
} from "@clutch-sh/core";
...
// highlight-start
const AmiiboDetails: React.FC<WizardChild> = () => {
  const amiiboData = useDataLayout("amiiboData");
  let amiiboResults = amiiboData.displayValue();
  if (_.isEmpty(amiiboResults)) {
    amiiboResults = [];
  }

  return (
    <WizardStep error={amiiboData.error} isLoading={amiiboData.isLoading}>
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
    </WizardStep>
  );
};
// highlight-end
...
```

Once the `amiiboData` data layout has been hydrated a table will be displayed with all amiibo matching the input criteria.

<Image alt="Amiibo Lookup Panel" src={useBaseUrl('img/docs/feature-development/details-panel.png ')} width="75%" variant="centered" />

Let's tie these two components together.

```tsx title="frontend/workflows/amiibo/src/hello-world.tsx"
...
import {
	Button,
	// highlight-next-line
	client,
	useWizardContext,
	...
} from "@clutch-sh/core";
...
// highlight-start
const Amiibo: React.FC<WorkflowProps> = ({ heading }) => {
  const dataLayout = {
    userInput: {},
    amiiboData: {
      deps: ["userInput"],
      hydrator: (userInput: { name: string }) => {
        return client
          .post("/v1/amiibo/getAmiibo", {
            name: userInput.name,
          })
          .then(response => {
            return response?.data?.amiibo || [];
          });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <AmiiboLookup name="Lookup" />
      <AmiiboDetails name="Details" />
    </Wizard>
  );
};

export default Amiibo;
// highlight-end
```

Here we define the data layout structure and pass it to the wizard along with the functional components we built above as children.

##### Configuration

The final step is registering this new Amiibo workflow with the Clutch app. First update the component name by replacing the scaffolded default (`HelloWorld`) in our workflow's registration function.

```tsx title="frontend/workflows/amiibo/src/index.tsx"
import type { BaseWorkflowProps, WorkflowConfiguration } from "@clutch-sh/core";
// highlight-next-line
import Amiibo from "./hello-world";

export interface WorkflowProps extends BaseWorkflowProps {}

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
  // highlight-next-line
        path: "/lookup",
        description: "Lookup all Amiibo by name.",
	// highlight-next-line
        component: Amiibo,
      },
    },
  };
};

export default register;
```

Next, open your `clutch.config.js` file and add the following:

```js title="clutch.config.js"
module.exports = {
  ...
  // highlight-start
  "@clutch-sh/amiibo": {
    landing: {
      trending: true,
    },
  },
  // highlight-end
}
```

If everything is in order you should see an Amiibo card on the homepage that you can use to access your workflow!

<Image alt="Amiibo Lookup Panel" src={useBaseUrl('img/docs/feature-development/landing-page.png ')} width="50%" variant="centered" />

That's it! You should be able to remove all the remaining generated code.

#### Test!

Run the frontend with the new workflow:
```bash
make frontend-dev
```

You should now be able to navigate to `http://localhost:3000/amiibo/lookup` to see your new workflow.

Try entering `Peach` into the input field and clicking search.

### 5. Additional Tutorials

Though not available yet, a future guide will cover how to incorporate the resolver pattern.
