---
title: Custom Gateway
{{ .EditURL }}
---

A custom gateway makes it possible to run custom or private extensions without forking or rewrites.

### Scaffolding Tool

The scaffolding tool will create a custom gateway with an example plugin.

From the `clutch` base directory, run:

```bash
make scaffold-gateway
```

The scaffolding tool will attempt to guess the correct path for Clutch based on the current shell's username and `GOPATH`. It will prompt the user to verify the path and customize it if needed.

Once the user has verified the path, the scaffolding tool will create the new frontend and backend files for a custom gateway from templates. Finally, it will automatically generate the Go and Javascript API code from protobuf using the `make api` target so the custom gateway is ready to build with one command.

After the scaffolding process has completed successfully, navigate to the path and run `make` to build the combined frontend and backend binary.

For more information on building Clutch see the [Local Build](/docs/getting-started/local-build) reference.

:::info
The scaffolding tool does not currently support the creation of individual components such as workflows.
:::

### Pushing to GitHub

After the custom gateway is created it only takes a few steps to push it to a new remote repository.

1. Create a new repository on GitHub, e.g. `example/clutch-custom-gateway`, private if desired.
1. `cd` to the directory that the scaffolding tool placed the custom gateway.
1. `git remote add origin ssh://git@github.com/example/clutch-custom-gateway`
1. `git add .`
1. `git commit -m "Initial commit"`
1. `git push origin master`
