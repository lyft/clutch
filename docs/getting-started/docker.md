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


### Running the Image

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
docker run --rm -p 8080:8080 -it lyft/clutch:latest
```
 
</TabItem>
<TabItem value="custom">

```bash
docker run --rm -p 8080:8080 \
-v /host/absolute/path/to/config.yaml:/clutch-config.yaml:ro \
-it lyft/clutch:latest
```
 
</TabItem>
<TabItem value="env">

```bash
docker run --rm -p 8080:8080 \
-e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
-e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
-it lyft/clutch:latest
```
 
</TabItem>
</Tabs>

To update the image on subsequent uses of the container:
```bash
docker pull lyft/clutch:latest
```

:::info Configuration
For more information on configuring Clutch, see the [Configuration Reference](/docs/configuration).
:::

### Accessing Clutch
:tada: Clutch should now be accessible from `localhost:8080` in the browser.

<img style={ {border: "1px solid black"} } alt="Clutch Landing Page Screenshot" src={useBaseUrl('img/docs/screenshot-landing.png')} width="50%" />

## Building the Container From Scratch

### Cloning the Repository
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
docker build -t clutch .
```

### Running the Local Image

Use the commands from the earlier step [Running the Image](#running-the-image), replacing `lyft/clutch:latest` with `clutch`, e.g. `docker run --rm -p 8080:8080 -it clutch` and access Clutch in the browser.

## Next Steps

- For more information on core components, see the [Components](/docs/components) reference.
- To better understand how custom components fit into Clutch, visit the [Architecture](/docs/about/architecture) reference.
- For documentation on developing custom components, check the [Development](/docs/development/guide) docs.

