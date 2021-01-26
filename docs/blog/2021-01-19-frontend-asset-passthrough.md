---
title: Frontend Asset Passthrough
authors:
  - name: Mike Cutalo
    url: https://github.com/mcutalo88
    avatar: https://avatars1.githubusercontent.com/u/2250844?s=460&u=24deb32096e9f892cc91a6ff1ca1af50193b1fbd&v=4
description: Frontend Asset Passthrough
image: https://user-images.githubusercontent.com/4712430/104760766-7a2c5980-5727-11eb-93f5-3296b23ba3a0.png
hide_table_of_contents: false
---

import useBaseUrl from '@docusaurus/useBaseUrl';

How Clutch handles serving frontend assets without a CDN in a distributed multi version deployment.
Will touch on how Lyft deploys Clutch and some of the early problems we faced when rolling it out.

<!--truncate-->

## Problem

When deploying Clutch at Lyft we noticed very early on when deploying a new version a subset of users would fail to load some frontend assets.
Clutch bundles all frontend assets into a single binary, which allows us to serve all frontend assets from disk easily.
However the frontend assets they were requesting did not exist on the hosts that were serving their request.

This problem was exacerbated when deploying to canary or waiting for the full production rollout to complete.
As there were at least 2 different versions of the application take production traffic.

On the surface it seems like the obvious solution would be to add a CDN and have that serve our static assets.
However our decision making processing of adding new dependencies to the Clutch architecture is simple, unless absolutely necessary find a different way to solve the problem.
We want Clutch to be extremely easy to deploy and configure with as few external dependencies as possible which are not tied to a specific cloud provider.

## Solution

Which leads us to our solution, simply called Frontend Asset Passthrough.
When a request for a static asset does not reside on disk Clutch will fallback to a configured provider to look for the asset.

Lets first take a look at the [proto configration](https://github.com/lyft/clutch/blob/890245e7d2a1bf91623a9e74b39f1083dbd5ea2c/api/config/gateway/v1/gateway.proto#L105-L119) and then will go into more detail.

```protobuf
message Assets {
  message S3Provider {
    string region = 1;
    string bucket = 2;
    string key = 3;
  }

  oneof provider {
    S3Provider s3 = 1;
  }
}
```

The `assets` configuration resides in the `gateway` configuration.

At Lyft we utilized AWS S3 as the provider of choice which we implemented first,
given this model its easy to add additional providers for any given deployment.

For S3 the configuration is simple, specify a `region`, `bucket` and the `key` where the assets live.
This does require you to also configure the `clutch.service.aws` service,
which enables Clutch fetch these assets with S3 API's.

```yaml
gateway:
  assets:
    s3:
      region: us-east-1
      bucket: static
      key: clutch

services:
  - name: clutch.service.aws
    typed_config:
      "@type": types.google.com/clutch.config.service.aws.v1.Config
      regions:
        - us-east-1
```

With this configured if a user requests an asset that does not reside on disk the configured passthrough provider will be utilized.
If the provider has the correct asset it will be served to the user.

Lets look at the diagram below to demonstrate this.
Our users makes a request, it goes through our load balancer in our case Envoy.
Envoy is load balancing all versions of Clutch that are currently deployed, in our example we have two Clutch versions deployed **v1** and **v2**.
Regardless of which version of the frontend a user might have all versions of the frontend assets live in S3,
if a clutch host does not have what the user is requesting it will check S3.

<img alt="Deploying" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-s3-logical.png')} width="75%"/>


## Deploying

However there is still an additional problem to solve, which will be unique to each deployment.
Will go over how at Lyft we manage uploading new frontend assets to S3 per deploy.

Below is a simplified version of our deployment process.
Early on in the deployment pipeline assets are uploaded to S3,
simply utilizing the aws cli to `aws s3 sync` the new assets to the target bucket.

Internally our CI/CD infrastructure utilize container images at most stages in the pipeline,
allowing us to easily access newly built assets to upload to S3.

<img alt="Deploying" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-s3-upload.png')} width="75%"/>

## Conclusion

Frontend asset passthrough enables distributed Clutch deployments to easily provide a fallback mechanism
when serving frontend assets.
With minimal configuration

## Contributing

If their is a provider that you need for your deployment, please contribute!
