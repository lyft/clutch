---
title: Local Build
{{ .EditURL }}
---

import useBaseUrl from '@docusaurus/useBaseUrl';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Requirements

In order to build Clutch, the following tools are required:

- Go ([golang.org](https://golang.org/doc/install))
- Node.js ([nodejs.org](https://nodejs.org/tr/download/package-manager/))
- Yarn ([yarnpkg.com](https://classic.yarnpkg.com/en/docs/install))

:::info
If you are building on OSX you'll need to install coreutils. This is easiest with [homebrew](https://brew.sh/).
```bash
brew install coreutils
```
:::

[Homebrew package manager](http://brew.sh/) is recommended for macOS users to manage these dependencies.

:::info
If running Clutch in a Docker container is preferred, see the [Docker](/docs/getting-started/docker) docs.
:::

## Building Clutch

#### 1. Clone
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

#### 2. Build

Run `make` to build a combined frontend and backend binary. The frontend is configured at build time by [clutch.config.js](https://github.com/lyft/clutch/blob/main/frontend/packages/app/src/clutch.config.js).

```bash
make
```

#### 3. Run
Launch Clutch with back-end configuration [clutch-config.yaml](https://github.com/lyft/clutch/blob/main/backend/clutch-config.yaml).

```bash
./build/clutch -c backend/clutch-config.yaml
```

#### 4. Use
:tada: Clutch should now be accessible from `localhost:8080` in the browser.

<img alt="Clutch Landing Page Screenshot" src={useBaseUrl('img/docs/landing-page.png')} width="50%" />

:::info
Clutch may have external dependencies, to run Clutch with mocked dependencies see [Mock Gateway](/docs/getting-started/mock-gateway).
:::


## Additional Build Targets

Clutch includes a comprehensive [`Makefile`](https://github.com/lyft/clutch/blob/main/Makefile) to simplify the execution of commands related to compiling, testing, and executing tools in the project.

The default target, e.g. running `make` in the root of Clutch, builds the frontend first and then packages it into the backend using [vfsgen](https://github.com/shurcooL/vfsgen).

The most commonly used `make` targets that every Clutch developer should know:

| Command | Description |
| --- | --- |
| `make api` | Re-generate frontend and backend API objects from changes to `.proto` files in `api/`.  |
| `make dev` | Compile and start the backend on `localhost:8080` and a hot-reloading frontend on `localhost:3000`. |
| `make` | Build a unified frontend and backend binary and place it in `build/clutch`. |
| `make lint-fix` | Fix any fixable linting errors. There are also targets for each of API, backend, and frontend, e.g. `make frontend-lint-fix`. |
| `make test` | Run tests. There are also targets for each of backend and frontend, e.g. `make backend-test`. |
