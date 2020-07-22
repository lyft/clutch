---
title: Roadmap
{{ .EditURL }}
---

import RoadmapItem from '@site/src/components/RoadmapItem';

Clutch ships with an extensible frontend and backend that makes building maintainable workflows and features easy. We're not done yet though! Some of our ideas for the near future of Clutch are listed here.

## Platform Improvements

<RoadmapItem title="Asynchronous Tasks" description="Long-running or asynchronous tasks are a normal part of safely performing infrastructure maintenance. Clutch will track, execute, and report on jobs originating from Clutch workflows." />

<RoadmapItem title="Remote Execution" description="Safely run diagnostic commands across a cluster or set of resources.">

Sometimes an operator just needs a shell to diagnose a problem. Observability systems do not ingest granular data because it's expensive even though it's readily available from command line tools.

Clutch workflows may also require remote command execution in the event that an API does not otherwise exist for remote commands.

`fab` is common command-line tool for this purpose. However it's slow and the output is extremely difficult to follow when dealing with large numbers of resources.

For more reading on this topic, see [Netflix Bolt](https://netflixtechblog.com/introducing-bolt-on-instance-diagnostic-and-remediation-platform-176651b55505?gi=59a3aad4070a).

</RoadmapItem>

<RoadmapItem title="Topology Cache" description="Cache resource lists and add typeahead to forms.">

While adapting any topology to Clutch is possible with an extension, large infrastructure topologies can be slow to search if based on vendor APIs. To make autocomplete possible and allow for a more responsive user experience, Clutch will periodically gather resource IDs and store them in its own online database for faster access.

</RoadmapItem>

<RoadmapItem title="Additional Gateways" description="Access Clutch APIs from other interfaces, e.g. the command-line or a Slackbot.">

Logic and safeguards that are implemented in Clutch should also be usable from other interfaces in addition to the React frontend. For example, informational Clutch APIs should be accessible from a standardized Slackbot.

</RoadmapItem>

<RoadmapItem title="UI Design System" description="Put UI development on rails and ensure a consistent user experience.">

UI components should be available for study and development in isolation and used consistently across the frontend.

For more reading on this topic, see [Storybook makes building stunning UIs organized and efficient](https://storybook.js.org/).

</RoadmapItem>

<RoadmapItem title="API / CLI Eject" description="Show corresponding CLI command and/or API call to user for replay and documentation purposes.">

Display the corresponding CLI command or raw API calls associated with an action before and after execution. This would allow operators to better understand what is happening throughout a workflow and seek additional documentation if necessary. Furthermore, if Clutch is partially unavailable or otherwise unable to perform an action, an operator could have the CLI command or API call executed in an environment with the proper permissions and access.

</RoadmapItem>

<RoadmapItem title="Feature Flags" description="Dynamically toggle feature availability to targeted users in Clutch." />


<RoadmapItem title="Additional Integrations" description="Clutch is a platform for integrating the tools engineers use on a consistent basis.">

We're building Clutch to be an everyday tool, consolidating disparate tools and sources information into one place that developers can visit to get things done and get back to developing new features. Integrating additional providers and technologies and combining them in intuitive and delightful ways will allow Clutch to further this mission.

</RoadmapItem>

## Safety Features

<RoadmapItem title="Two-phase Approval" description="Sensitive actions should require a '+1' from another operator before executing to ensure correctness." />

<RoadmapItem title="Rate Limiting" description="Ensure destructive actions are not performed too frequently." />

<RoadmapItem title="Challenge Modals" description="Require additional confirmation for potentially risky actions." />

## Alert Handling

<RoadmapItem title="Annotation" description="Append links to workflows or metadata from Clutch integrations onto alerts." />

<RoadmapItem title="Auto-remediation" description="Automatically execute a scripted runbook in response to an alert.">

Many workflows are to resolve an immediate problem based on an alert or page. If the pages result in the same mechanical action every time, the page can also be remediated using the same action programatically being instead of requiring a human operator.

For more reading on this topic, see [Netflix Winston](https://netflixtechblog.com/introducing-winston-event-driven-diagnostic-and-remediation-platform-46ce39aa81cc).

</RoadmapItem>

## Envoy Capabilities

Clutch :heart: Envoy. The project was modeled after Envoy and inspired by its success. We are investing significantly in Envoy functionality in the near future.

<RoadmapItem title="Advanced Fault Injection" description="Target and inject faults between other Envoys or third-party networks." />
<RoadmapItem title="Config Dump Diff" description="Select two Envoys and diff their configuration remotely." />
<RoadmapItem title="Real-time Stats Viewer" description="Drop a probe on an Envoy for per-second stats visualization." />
<RoadmapItem title="Tap Interface" description="Tap an Envoy with the filter enabled so that it dumps all requests and response information to your screen." />
<RoadmapItem title="Runtime Manager" description="Envoy has a runtime system that enables dynamic reconfiguration. It needs a UI." />
<RoadmapItem title="Route Manager" description="Envoy routing tables can be extremely complex. " />
<RoadmapItem title="Config Generator" description="Need SNI configuration for a remote host? Generate it and other complex configurations with a guided interface." />

## Suggestions?

We're excited to hear from you. Visit the [Community](/docs/community) page for ways to contact the Clutch team and share your ideas.
