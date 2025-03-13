---
title: Docker
{{ .EditURL }}
---

import useBaseUrl from '@docusaurus/useBaseUrl';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

The instructions below will get Clutch up and running in a Docker container.

:::info Building Locally
If building the binary outside of a container is preferred, see the [Local Build](/docs/getting-started/local-build) docs.
:::

## Prerequisites
Docker is required to use the examples below, see [Get Docker](https://docs.docker.com/get-docker/) for information on installing Docker itself.

## Using the Docker Container

Clutch provides a [`Dockerfile`](https://github.com/lyft/clutch/blob/main/Dockerfile) that builds and runs a version of Clutch with all [core components](/docs/components) compiled in.

### Building the Container From Scratch

#### Cloning the Repository

Start by cloning the Clutch repository and entering into the source directory.

<Tabs
  defaultValue="https"
  values={[
    {label: 'HTTPS', value: 'https'},
    {label: 'SSH', value: 'ssh'},
  ]}>

<TabItem value="https">

```bash
git clone https://github.com/lyft/clutch
cd clutch
```

</TabItem>
<TabItem value="ssh">

```bash
git clone git@github.com:lyft/clutch
cd clutch
```

</TabItem>
</Tabs>

### Building the Image

Build the Docker image locally. The frontend will build followed by the backend. The frontend is copied into the backend and compiled into a single binary.

```bash
docker build -t clutch:latest .
```

:::info
If you encounter this error:
```
 => ERROR [gobuild 7/7] RUN make backend-with-assets                                                                            0.3s
------
 > [gobuild 7/7] RUN make backend-with-assets:
0.196 /go/src/github.com/lyft/clutch
0.196 Running pre-flight checks...
0.212 Pre-flight checks satisfied!
0.213 cd backend && go run cmd/assets/generate.go ../frontend/packages/app/build && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=0.0.0"
0.218 go: errors parsing go.mod:
0.218 /go/src/github.com/lyft/clutch/backend/go.mod:3: invalid go version '1.23.0': must match format 1.23
0.219 make: *** [Makefile:47: backend-with-assets] Error 1
------
Dockerfile:19
--------------------
  17 |     COPY --from=nodebuild ./frontend/packages/app/build ./frontend/packages/app/build/
  18 |     
  19 | >>> RUN make backend-with-assets
  20 |     
  21 |     # Copy binary to final image.
--------------------
ERROR: failed to solve: process "/bin/sh -c make backend-with-assets" did not complete successfully: exit code: 2
```
As mentioned in this [Github Issue](https://github.com/lyft/clutch/issues/3173), if clutch won't build after running that command you'll need to modify the golang version to `1.23.0-bookworm` and the grc.io image version to `distroless/base-debian12`
:::


## Running the Image

On every commit to the main branch, a Docker image is built and published to [Docker Hub](https://hub.docker.com/r/lyft/clutch) tagged with `latest` and the commit SHA in the form of `sha-<abcdef>`.

The default configuration in [`backend/clutch-config.yaml`](https://github.com/lyft/clutch/blob/main/backend/clutch-config.yaml) is used.
If desired, use a custom configuration with the Docker image by mapping it into the container.



<Tabs
  defaultValue="default"
  values={[
    {label: 'Default Configuration', value: 'default'},
    {label: 'Custom Configuration', value: 'custom'},
    {label: 'With Environment Variables', value: 'env'},
  ]}>
<TabItem value="default">

```bash
docker run --rm -p 8080:8080 -it clutch:latest
```

</TabItem>
<TabItem value="custom">

```bash
docker run --rm -p 8080:8080 \
-v /host/absolute/path/to/config.yaml:/clutch-config.yaml:ro \
-it clutch:latest
```

</TabItem>
<TabItem value="env">

```bash
docker run --rm -p 8080:8080 \
-e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
-e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
-it clutch:latest
```

</TabItem>
</Tabs>

:::info Configuration
For more information on configuring Clutch, see the [Configuration Reference](/docs/configuration).
:::

### Accessing Clutch
:tada: Clutch should now be accessible from `localhost:8080` in the browser.

<img alt="Clutch Landing Page Screenshot" src={useBaseUrl('img/docs/landing-page.png')} width="50%" />


## Next Steps

- For more information on core components, see the [Components](/docs/components) reference.
- To better understand how custom components fit into Clutch, visit the [Architecture](/docs/about/architecture) reference.
- For documentation on developing custom components, check the [Development](/docs/development) docs.
