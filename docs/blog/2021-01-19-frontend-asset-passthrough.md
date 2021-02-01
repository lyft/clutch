---
title: Deploying Hashed Frontend Builds Without Interruption
authors:
  - name: Mike Cutalo
    url: https://github.com/mcutalo88
    avatar: https://avatars1.githubusercontent.com/u/2250844?s=460&u=24deb32096e9f892cc91a6ff1ca1af50193b1fbd&v=4
description: How caching and hashed frontend builds can lead to a blank screen; and how to fix it.
image: https://user-images.githubusercontent.com/2250844/106201558-7956e700-616d-11eb-887d-28410b67d558.png
hide_table_of_contents: false
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Ever wonder how Clutch handles serving frontend assets without a CDN in a distributed multi-version deployment?
This article will touch on how Lyft deploys Clutch and some of the early problems we faced when rolling it out.

<!--truncate-->

## Problem

During our deploys of Clutch at Lyft we noticed very early on that when a new version is deployed,
a subset of users would fail to load the webpage (only seeing a blank page) as some of the frontend assets could not be resolved.

<img style={ {border: "1px solid black"} } alt="Unable to Load" src="https://user-images.githubusercontent.com/2250844/106497808-5ed58400-6473-11eb-9fbe-fa57a0e18969.png" />

```text
Uncaught SyntaxError: Unexpected token '<'
```

Clutch embeds frontend assets into its binary, simplifying our architecture by omitting the need for a CDN or a separate frontend server.
However, during deploys, if deploying to a canary environment or one of multiple availability zones, there will briefly be two different versions of the application taking traffic.
In this intermediary state, the frontend assets requested by the user may not yet exist on the hosts serving that request while the rollout is progressing.

Typically, organizations solve this problem by adding a CDN to their service architecture and serving static assets from there.
However, we decided to go a different route, as the Clutch architecture and design philosophy values simplicity and avoiding adding new dependencies unless absolutely necessary.
We want to keep Clutch easy to deploy and configure, and cloud provider agnostic, with as few external dependencies as possible.
In addition to the design considerations, not using a CDN does provide other advantages.
While CDNs are inherently public facing, asset passthrough is private, making it easier to secure your deployment.

Before we go through the illustration of the problem below we first need to understand how Webpack builds frontend bundles by default when using `create-react-app`.
Clutch uses webpack as the build system for the frontend,
when building a new release [webpack templates the output filename](https://webpack.js.org/guides/caching/#output-filenames) to include a content hash eg: `main.[contenthash].chunk.js`.
This uniqueness of this content hash allows Clutch to cache bust what the browser has locally allowing the new version to be requested.

<img alt="Problem Diagram" src="https://user-images.githubusercontent.com/2250844/106201546-765bf680-616d-11eb-83d3-c70cf93ba252.png" />

The client makes a request to Clutch and a canary host responds with an `index.html` page with a script tag asking for `main.a3762de8.chunk.js`.
Illustrated by the green and red arrows, there is only one host which has the correct asset the client is asking for, the canary.
If the client's request for `main.a3762de8.chunk.js` goes anywhere else the page will fail to load.

This is not a problem if you are using a CDN, as traditionally there is a backing object storage such as S3 where all of the applications assets are stored.
If a CDN does not have an asset it checks its backing object storage, warms the CDN cache, and serves it to the client.
Alternatively, sticky sessions or session affinity could be used to remediate this problem,
as the requesting client would always be routed to a specific version of the application.

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
however, Clutch is extensible and can easily be extended with support for additional blob storage providers.

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

With this configured, if a user requests an asset that does not reside on disk, Clutch will attempt to serve the file from the configured passthrough provider.
If the provider has the correct asset it will be served to the user.

Let's look at the diagram below to demonstrate this.
When a user makes a request, it goes through a load balancer, in this example, Envoy.
Envoy is load balancing to all instances of Clutch using the round-robin algorithm.
In this example we have two Clutch versions deployed: `v1` and `v2`.
Regardless of which version of the frontend a user might request, all versions of the frontend assets live in S3,
so if a Clutch host does not have what the user is requesting, it can check S3.

<img alt="Logical Architecture" src="https://user-images.githubusercontent.com/2250844/106201558-7956e700-616d-11eb-887d-28410b67d558.png" />


## Deploying

Now that we have configured a passthrough asset provider, there is still an additional problem to solve, which will be unique to each Clutch deployment.
This problem, namely, is how exactly the assets end up in the fallback provider in the first place.
As an example, I will go over how we solved this problem at Lyft.

Below is a simplified version of our deployment process.
Early on in the deployment pipeline assets are uploaded to S3,
utilizing the [AWS CLI](https://docs.aws.amazon.com/cli/latest/reference/s3/sync.html) to `aws s3 sync` the new assets to the target bucket.

<img alt="Deploying" src="https://user-images.githubusercontent.com/2250844/106201560-7956e700-616d-11eb-9f9d-a4b1345bcf41.png" />

## Conclusion

Through Frontend Asset Passthrough, a distributed Clutch deployment is able to provide a fallback mechanism
when serving frontend assets with minimal configuration and without the additional complexity of a CDN,
ensuring the correct assets are served regardless of the divergence in deployed versions across a Clutch fleet.

## Contributing

If there is a provider that you need for your deployment, please open an [issue](https://github.com/lyft/clutch/issues) or consider [contributing](https://github.com/lyft/clutch#contributing)!
If you have any questions or would like to chat with the team, join us in our [Slack](/docs/community)!
