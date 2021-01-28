---
title: Frontend Asset Passthrough
authors:
  - name: Mike Cutalo
    url: https://github.com/mcutalo88
    avatar: https://avatars1.githubusercontent.com/u/2250844?s=460&u=24deb32096e9f892cc91a6ff1ca1af50193b1fbd&v=4
description: Frontend Asset Passthrough
image: https://user-images.githubusercontent.com/2250844/106030368-f8bbbc00-6082-11eb-819f-fd3798a53489.png
hide_table_of_contents: false
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Ever wonder how Clutch handles serving frontend assets without a CDN in a distributed multi-version deployment?
This article will touch on how Lyft deploys Clutch and some of the early problems we faced when rolling it out.

<!--truncate-->

## Problem

During our deploys of Clutch at Lyft we noticed very early on that when a new version is deployed,
a subset of users would fail to load the webpage (only seeing a blank page) as some of the frontend assets could not be resolved.

<img alt="Unable to Load" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-unable-to-load.png')} />

Clutch bundles all frontend assets into a single binary, which allows us to serve them from disk, simplifying our architecture by omitting the need for a CDN.
However, this means that during deploys, particularly when deploying to a canary environment or multiple availability zones, there will briefly be two different versions of the application taking traffic.
In this intermediary state, the frontend assets requested by the user may not yet exist on the hosts serving that request while the rollout is progressing.

Typically, organizations solve this problem by simply adding a CDN to their service architecture and serving static assets from there.
However, we decided to go a different route, as the Clutch architecture and design philosophy values simplicity and avoiding adding new dependencies unless absolutely necessary.
We want to keep Clutch easy to deploy and configure, and cloud provider agnostic, with as few external dependencies as possible.
In addition to the design considerations, not using a CDN does provide other advantages.
As CDN are inherently public facing, asset passthrough is private making it easier to secure your deployment.

To illustrate the problem lets go through the diagrams below.

<img alt="Problem Diagram" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-problem-diagram.png')} />

The client makes a request to Clutch and a canary host responses, an `index.html` page with script tag asking for `main.a3762de8.chunk.js`.
Illustrated by the green and red arrows there is only one host which has the correct asset the client is asking for, the canary.
If the clients request for `main.a3762de8.chunk.js` goes anywhere else the page will fail to load.

## Solution

This leads us to our solution, which we call Frontend Asset Passthrough.
How this works is simple: when a request for a static asset does not reside on disk,
Clutch will fallback to a configured provider to look for the asset.

Let's first take a look at the [proto configuration](https://github.com/lyft/clutch/blob/890245e7d2a1bf91623a9e74b39f1083dbd5ea2c/api/config/gateway/v1/gateway.proto#L105-L119) and then will go into more detail.

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

The configuration can be found in the [gatewayOptions](/docs/configuration#gatewayoptions).

At Lyft we utilize AWS S3 as the provider of choice and have implemented this setup as an example;
however, the Clutch protobuf model is extensible and can easily be extended to add additional providers for any given deployment.

For S3 the configuration is simple: specify a `region`, `bucket`, and the `key` where the assets live.
This does require you to also configure the `clutch.service.aws` service,
which allows Clutch to fetch these assets via S3 APIs.

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

With this configured, if a user requests an asset that does not reside on disk the configured passthrough provider will be utilized.
If the provider has the correct asset it will be served to the user.

Let's look at the diagram below to demonstrate this.
A user makes a request, it goes through a load balancer, in this example is Envoy.
Envoy is load balancing all versions of Clutch that are currently deployed.
In this example we have two Clutch versions deployed: `v1` and `v2`.
Regardless of which version of the frontend a user might have, all versions of the frontend assets live in S3,
so if a Clutch host does not have what the user is requesting it can simply check S3.

<img alt="Logical Architecture" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-s3-logical.png')} />


## Deploying

Now that we have configured a passthrough asset provider, there is still an additional problem to solve, which will be unique to each Clutch deployment.
This problem, namely, is how exactly the assets end up in the fallback provider in the first place.
As an example, I will go over how we solved this problem at Lyft.

Below is a simplified version of our deployment process.
Early on in the deployment pipeline assets are uploaded to S3,
simply utilizing the [AWS CLI](https://docs.aws.amazon.com/cli/latest/reference/s3/sync.html) to `aws s3 sync` the new assets to the target bucket.

<img alt="Deploying" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-s3-upload.png')} />

## Conclusion

Through Frontend Asset Passthrough, a distributed Clutch deployment is able to provide a fallback mechanism
when serving frontend assets with minimal configuration and without the additional complexity of a CDN,
ensuring the correct assets are served regardless of the divergence in deployed versions across a Clutch fleet.

## Contributing

If there is a provider that you need for your deployment, please consider contributing!
