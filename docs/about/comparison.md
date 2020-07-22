---
title: Comparison to Other Tools
sidebar_label: Comparison
{{ .EditURL }}
---

Clutch was designed and built after careful study of the other available options for solving infrastructure user and developer experience issues. Clutch patterns itself off of the successful model used by Envoy Proxy, also developed at Lyft, offering a rich, configuration-driven alternative to legacy tools.

### Clutch Differentiators
- :rocket: **Ease of Deployment**
  - A single binary contains the React frontend and Go backend. No complicated set of microservices required to run.
- :nut_and_bolt: **Complementary**
  - Clutch is agnostic to the environment's tech stack via rich configuration and extension patterns. Other systems require migration to a reference architecture.
- :electric_plug: **Extensible**
  - Forks are never required to add new features.
  - Plugins are configurable, allowing them to cater to a variety of different environments and systems.
- :computer: **Developer-friendly**
  - Interfaces and abstractions are straightforward for infrastructure teams, particularly those with little to no frontend experience.
  - API schemas minimize frontend and backend boilerplate with generated code. Schema changes automatically apply to generated objects and interfaces.
  - The backend includes a resource location framework, first-class auditing and authorization support, integrated API schemas, and gRPC and JSON support.
  - The frontend offers abstractions for state management, multi-step workflows, server-driven resource location, and data visualization.
- :sparkles: **User-friendly**
  - Workflows offer clear and concise visual cues so users are comfortable clicking buttons.
  - To prevent accidents, operational safety is a key consideration in Clutch. Clutch is not intended primarily for superusers like most other tools.

## Backstage
[Backstage](https://github.com/spotify/backstage), released in March 2020, is the result of Spotify's efforts to release an open source version of their internal platform that powers the day-to-day developer experience at Spotify. The internal version relies on a comprehensive set of microservices that Spotify is currently unable to release alongside the public version of the platform. Backstage open sourced with a small set of functionality and has been building additional features in the open at a steady pace.

Backstage is billed as an extensible frontend platform. It consists of scaffolding tools and libraries to create a new frontend deployment of the Backstage console and custom plugins for the console. The developer experience in Backstage is tailored towards experienced frontend engineers who are looking for a lightweight framework to integrate their code into a single application. At the time of this writing, Backstage's backend is in very early stages. The scaffolding tools that simplify frontend deployment and creating plugins have no support for the backend.

We're excited to see what the Backstage team releases going forward, especially since Spotify's internal version of Backstage is an example of the developer console done right. Clutch currently offers a suite of tools for developing maintainable workflows with richer frontend and backend patterns.

## Spinnaker

[Spinnaker](https://www.spinnaker.io/) bills itself as "an open source, multi-cloud continuous delivery platform". It comes with an impressive host of features and support for various cloud providers. It has a large community with an impressive roster of users including Airbnb, Box, Cisco, Google, and Netflix. 

Spinnaker is architecturally complex, comprised of [more than 10 microservices](https://www.spinnaker.io/reference/architecture/#system-dependencies). This may be difficult to maintain for smaller organizations with smaller teams. Even doing a test deployment of Spinnaker for evaluation would require significant investment. Spinnaker is also not necessarily a good option for organizations that already have large deployments. In order to take advantage of Spinnaker, an organization has to completely migrate to its concepts of clusters, applications, and services. Finally, many users of Spinnaker have built additional tooling for Spinnaker deployments outside of Spinnaker itself. Clutch could serve a similar purpose for users of Spinnaker but with a much lower bar for entry.

For organizations that do not have the appetite to run or migrate to Spinnaker, Clutch allows building tooling that is complementary to any existing architecture.

## Jenkins

Jenkins is a popular platform used by many organizations for continuous integration and delivery. However, most with experience deploying and running it at scale will caution against relying on it as a platform for a growing organization. It is possible to use it as a workflow engine of sorts with parameterized jobs. This also gives a developers a UI to execute, view logs, and audit usage of a job. However, Jenkins is not considered maintainable and user-friendly. Parameterized jobs don't offer any kind of built-in input validation or autocompletion. While it is possible to develop custom plugins for Jenkins, the API is not friendly and plugins do not integrate well with the primitives available in Jenkins.

## Fabric (`fab`)

CLI frameworks like [Fabric](https://www.fabfile.org/) are effectively a collection of commands with some abstractions for targeting resources. Relying on local execution comes with all kinds of environment and permission headaches.

Furthermore, CLIs are not always the best tool for the job. Actions that are performed infrequently or above a certain level of complexity have high cognitive load on a CLI. For more reading on this topic see [IBM Design: The value of web UIs for CLI-oriented users](https://medium.com/design-ibm/real-developers-dont-use-uis-daea7404fb4e).

## Vendor Tools

Vendor tooling and application-specific tooling offer many lessons for Clutch. Clutch is also a platform to rebundle core, commonly-used functionality from vendor tools.   

### AWS

#### Console
Most developers are comfortable with the console, and it presents resources and actions in a mostly organized and discoverable fashion. However, the amount of data presented to users can be overwhelming. The console interface presents users with as much information as possible, regardless of whether it is relevant to the task at hand.

For organizations with large deployments (1,000+ resources), the console is extremely slow, and many pages will fail to load. If operating across regions, each region and the resources in it must be viewed separately in the UI. The console is not sensitive to specific IAM permissions and displays an error instead of redacting inaccessible pages. In many cases, forms only display a permissions error when clicking 'Apply'. This is frustrating if it took time to locate the proper form and consider the correct values. AWS is slowly addressing some of these pain points, such as loading performance, but this work is still in its early stages.

#### `aws-cli`
The AWS CLI mirrors the console in its complexity through comprehension. Its various sub and sub-sub commands can leave newcomers at a loss when searching for the appropriate actions. Operators often create "cheat sheets", aliases, or wrap the SDK to simplify on-call. However, these can quickly become out of date and provide no auditing or visibility into the side-effects of a command.

Guardrails to avoid using the CLI and/or console in an unsafe manner are non-existent because AWS has no concept of what "correct" means for a given deployment. Accidents can result in expensive incidents and are not uncommon when dealing with the terse CLI. The possibility of making a mistake can also yield slower response to an ongoing incident.

### Kubernetes Tooling  
Organizations that run multiple clusters cannot make use of the K8s dashboard for unified tooling since it only supports a single cluster (also known as a context). The CLI tool is similar in that it is designed around single-cluster operations.

Some interesting projects exist address the shortcomings of the core K8s tools, such as [hjacobs/kube-ops-view](https://github.com/hjacobs/kube-ops-view) and [vmware-tanzu/octant](https://github.com/vmware-tanzu/octant).
