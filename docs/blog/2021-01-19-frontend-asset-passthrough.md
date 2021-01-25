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

Handling frontend asset management with a distribtued deployment.

<!--truncate-->

## Problem

Clutch serves all frontend assets from disk as they are bundled as a singnal binary.
When deploying Clutch at Lyft, we noticed very early on that when we deployed a new version some of our users would fail to load some web assets.
As the frontend assets they were requesting did not exist on the hosts that were servering their request.
This was most obvious when we were in canary or as prod hosts rolled out the newer version of the application.

## Solution

On the surface it seems like it would be obvious to add a CDN to serve our static assets.
However our design philophicy for Clutch is simplicity above all else.
We want Clutch to be extreamly easy to deploy with as little external dependenices as possible.

Which leads us to our solution, simply called Frontend Assest Passthrough.
When a request for a static asset does not reside on disk Clutch will fallback to a configured provider to look for the asset.

Lets take a look at the [configration](https://github.com/lyft/clutch/blob/890245e7d2a1bf91623a9e74b39f1083dbd5ea2c/api/config/gateway/v1/gateway.proto#L105-L119) and then will go into more detail.

```protobuf
// Assets configuration provide a passthrough host for frontend static assets.
// This is useful if you dont have the ability to enable sticky sessions or a blue/green deploy system in place.
message Assets {
  // To use the S3Provider you must have the AWS service configured
  message S3Provider {
    string region = 1;
    string bucket = 2;
    // key is the path to clutchs frontend static assets in the configured bucket
    string key = 3;
  }

  oneof provider {
    S3Provider s3 = 1;
  }
}
```

At Lyft we utlized AWS S3 as the provider of choice which we implmeneted first,
given this model its easy to add additonal providers for any given deployment.

For S3 the configuration is simple, specify a `region`, `bucket` and the `key` where the assets live.

However there is still an addtional problem to solve, which will be unique to each deployment.
Will go over how at Lyft we manage uploading new frontend assets to S3 per release.

## Deploying

Below is a simplifed version of our deployment process.
Early on in the deployment pipeline assets are uploaded to S3,
utilizing the aws cli to `aws s3 sync` the new assets to the target bucket.

Internally our CI/CD infrastructure utilize container images at most stages in the pipeline,
allowing us to easily access newly built assets to upload to S3.

<img alt="Deploying" src={useBaseUrl('img/docs/blog/fe-asset-passthrough-s3-upload.png')} width="75%" />

## Conclusion

Frontend asset passthrough enables 

## Contributing
