---
title: Community
{{ .EditURL }}
---

import CommunityCard from '@site/src/components/CommunityCard';

Clutch was built to be an ecosystem, and community is an important part of that. We value all contributions and ideas, big or small.

## Discussion
<CommunityCard icon="slack" to="https://join.slack.com/t/lyftoss/shared_invite/zt-casz6lz4-G7gOx1OhHfeMsZKFe1emSA">

Join the **Lyft OSS Slack** and come see us in channel **`#clutch`**!

</CommunityCard>

<CommunityCard icon="twitter" to="https://twitter.com/clutchdotsh">

Follow us on **Twitter** for news and updates `@clutchdotsh`.

</CommunityCard>

<CommunityCard icon="github" to="https://github.com/lyft/clutch/issues">

File a GitHub issue in **`lyft/clutch`**.

</CommunityCard>

<CommunityCard icon="calendar" to="https://calendly.com/clutchsh/office-hours">

Schedule time during our office hours with **Calendly**.

*Available in 15-minute slots on Tuesdays from 11am - 12pm Pacific.*

</CommunityCard>

## Contributing

Please view the [Code of Conduct](https://github.com/lyft/clutch/blob/main/CODE_OF_CONDUCT.md) if you're interested in contributing.

To have pull requests merged, contributors must sign the [Contributor License Agreement](https://oss.lyft.com/cla/clas/1.0) using their GitHub credentials at the [Lyft CLA Service](https://oss.lyft.com/cla).

### Core

The core project, hosted in [`lyft/clutch`](https://github.com/lyft/clutch), is designed to be universal and broadly applicable to a large variety of organizations. Code that is overly-specific to a single organization will not be accepted into the core. However, Clutch was designed to allow custom code without having to fork. See [Custom Extensions](/docs/usage/extensions) for more details.

It's also possible to contribute features incrementally to the core. For example, you can contribute an API definition but keep the implementation private.

In the event that one of Clutch's APIs does not work for your organization and you can't figure out how to open an agnostic change to the core, let us know! We will work with you to figure out a path forward.

If new to the project and looking for a way to contribute, check out issues marked as [`good first issue`](https://github.com/lyft/clutch/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) on GitHub.

### Non-core

We encourage you to open source your extensions, whether modules, resolvers, services, or middleware, even if they're not the right fit for the core project. Others may do the work for you to generalize them for use in core if there's enough interest.

Please tag your repo with the [`#clutch-plugin`](https://github.com/topics/clutch-plugin) topic! For more details see [Github Docs: Topics](https://help.github.com/en/github/administering-a-repository/classifying-your-repository-with-topics).
