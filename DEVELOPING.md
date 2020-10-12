# Developing Clutch

## Backend

There are many different ways clutch can be configured for backend development.
Depending on your use case you can enable / disable features in the `clutch-config.yaml`
to start up clutch with a minimal set of dependencies.
Here we're going to cover some of the common use cases.

### Prerequisites

* [docker](https://docs.docker.com/engine/install/ubuntu/) (for *nix) or [docker-for-mac](https://docs.docker.com/docker-for-mac/install/)
  > NOTE: If you prefer to use a different docker environment such as `docker-machine`,
    you must ensure port-forwarding is configured properly so `clutch` can reach
    its dependencies such as the datastore.
    You can find a list of ports to expose in the `docker-compose.yaml` under
    the `expose:` list for each container.

* [docker-compose](https://docs.docker.com/compose/install/)
  > NOTE: if your using docker-for-mac, docker-compose is already included.

* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
  > NOTE: Only required for the kubernetes usecase

### The Basics

To start the backend simply run, `make backend-dev` which will watch for
changes and reload the server for you.

### If you need a database

First add the postgres service to your `clutch-config.yaml`.

```yaml
services:
  - name: clutch.service.db.postgres
    typed_config:
      "@type": types.google.com/clutch.config.service.db.postgres.v1.Config
      connection:
        host: 0.0.0.0
        port: 5432
        user: clutch
        ssl_mode: DISABLE
        dbname: clutch
        password: clutch
```

1. From the root of the clutch project run the docker-compose command

    ```sh
    # This will start the postgres database
    docker-compose up -d
      ```

2. Run the database migration script against your local datastore,
instuction can be found [here](./backend/cmd/migrate/README.md).

3. Finally start up clutch and develop

    ```sh
    # Runs the clutch backend
    make backend-dev
    ```

### If you need a kubernetes cluster

First add the Kubernetes configuration to your `clutch-config.yaml`.

```yaml
modules:
  - name: clutch.module.k8s
...
services:
  - name: clutch.service.k8s
    typed_config:
      "@type": types.google.com/clutch.config.service.k8s.v1.Config
...
resolvers:
  - name: clutch.resolver.k8s
```

Running through the commands below will spin up a local Kubernetes cluster in docker.
This will also create a few Kubernetes resources so you can immediately start testing against them.
Envoy `deployments` & `HPAs` will be created in a `envoy-staging` and `envoy-production` namespace.

```sh
# This will start a local kubernetes cluster that runs as a single docker container.
# The cluster will be seeded with a few resources so you can start testing immediately.
make k8s-start

# This exports a `KUBECONFIG` which is read in by clutch as well as kubectl
eval $(make k8s-env)

# Runs the clutch backend
make backend-dev
```
