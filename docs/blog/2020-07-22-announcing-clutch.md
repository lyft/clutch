---
title: Announcing Clutch, the Open-source Platform for Infrastructure Tooling
author: Daniel Hochman
author_title: Lyft Engineer
author_url: https://github.com/danielhochman
author_image_url: https://user-images.githubusercontent.com/4712430/87979981-839a7900-ca98-11ea-9d35-07c01b4cec14.png
tags: [hello, world]
description: This is my first post.
image: https://user-images.githubusercontent.com/4712430/104760766-7a2c5980-5727-11eb-93f5-3296b23ba3a0.png
hide_table_of_contents: false
---

Today we are excited to announce the open-source availability of [Clutch](https://clutch.sh/), Lyft’s extensible UI and API platform for infrastructure tooling. Clutch empowers engineering teams to build, run, and maintain user-friendly workflows that also incorporate domain-specific safety mechanisms and access controls.

<!--truncate-->

---

**Note**: [This article](https://eng.lyft.com/announcing-clutch-the-open-source-platform-for-infrastructure-tooling-143d00de9713) was originally published at [eng.lyft.com](https://eng.lyft.com/).

**Nota**: *[Este artículo](https://medium.com/lyft-engineering-en-espa%C3%B1ol/anunciando-clutch-la-plataforma-de-c%C3%B3digo-abierto-para-administraci%C3%B3n-de-infraestructura-855dafe4380a) también está en español a [eng-espanol.lyft.com](http://eng-espanol.lyft.com/).*

---

By [Daniel Hochman](https://github.com/danielhochman) and [Derek Schaller](https://github.com/dschaller)


Clutch ships with several features for managing platforms such as AWS, Envoy, and Kubernetes with an emphasis on extensibility so it can host features for any component in the stack.

![Diagram showing Clutch connecting multiple systems](https://user-images.githubusercontent.com/4712430/104760766-7a2c5980-5727-11eb-93f5-3296b23ba3a0.png)

The dynamic nature of cloud computing has significantly reduced the adoption cost of new infrastructure. The [CNCF Landscape](https://landscape.cncf.io/) currently tracks over 300 open source projects and 1,000 more commercial offerings. Although organizations can rapidly adopt these projects and vendors, each new technology comes with its own set of configuration, tooling, logs, and metrics. Allowing developers to quickly and safely make changes throughout the stack requires significant upfront and ongoing investment in tooling, which most organizations fail to account for. Therefore, while new infrastructure is increasingly easy to *adopt*, it is difficult to scale the *management* of new components, especially as the complexity of the overall platform and the size of the engineering team grows. Clutch solves this problem by enabling infrastructure teams to deliver intuitive and safe interfaces for infrastructure management to their entire engineering organization.

Clutch is the result of a year-long development cycle to address deficiencies in Lyft’s developer experience and tooling. At its core, Clutch is made up of two main components. The **Go backend**  is designed to be an extensible infrastructure control plane, bringing together a patchwork collection of systems behind a single protobuf-driven API with common authorization, observability, and audit logging. API definitions for the backend automatically generate clients for the frontend via protobuf tooling. The **React frontend** is a pluggable and workflow-oriented UI allowing users and developers to create new features behind a single pane of glass with less code, less prior JavaScript knowledge, and less maintenance.

## Design and Architecture
What differentiates Clutch from other solutions in the developer tooling space? At the outset of the project, we did a thorough analysis of the existing tools before building our own.

Our main goals for the tooling we intended to adopt were:

- **Reduce the Mean Time To Repair (MTTR)** infrastructure when responding to an alert. Engineers were spending too long reading through runbooks and navigating complex tooling when responding to pages while on-call.
- **Eliminate accidental outages** when performing maintenance tasks. Significant outages have occurred as the result of a user missing a warning in a runbook or even modifying or deleting the wrong resource altogether, e.g. one they thought was unused but had significant traffic and usage.
- **Enforce granular permissions and audit all activity in a common format.** Some permissions are too broad because the vendor access controls do not support fine-grained control. Also, while we were collecting audit logs from various tools for security purposes, it was hard to distill that data into actionable insights on how we could improve our tools.
- **Provide a platform that greatly eases the development of future tools. At Lyft’s scale, projects with a large scope are rarely successful if they don’t account for contributions outside of the immediate team. We don’t have enough resources to build every feature that Lyft needs, much less support it.**

We started by looking at the shortcomings of the available vendor UIs. Vendor tools are slow (and in some cases dangerous) due to a lack of specialization. They require unnecessary steps to perform common tasks and present a superset of the necessary information. There are generally few guardrails beyond simple access controls. The result is that an operator may perform an action that seems innocuous but actually degrades the system. On the other hand, they may be unfamiliar with the tool such that it results in delayed remediation. Ideally, engineers are only on-call every four to six weeks. It’s easy to forget how to use a tool, especially considering the possibility of going multiple on-call cycles without needing to interact with a particular system.

The net outcome of relying on vendor tooling is high cognitive load due to fragmentation and information sprawl. In contrast, a vendor-agnostic tool like Clutch can unify disparate systems with a clear and consistent UX and offer specialized functionality to perform common tasks with as few clicks and as little training as possible.

Next, we turned to the open-source community. We found that the scope of open-source infrastructure management tooling is still usually limited to a single system and not designed for extensive customization. This does not address the problem of cognitive load and fragmentation. Also, while there are other frontend frameworks for building consoles, none of them incorporate an integrated backend framework with authentication, authorization, auditing, observability, API schemas, and a rich plugin model. There is a popular continuous delivery platform that addresses many of the same overarching issues as Clutch (e.g., lowering MTTR, user-friendly UI). However, it requires significant investment in running many microservices and migrating applications to a structure different from our own. Clutch’s backend simplifies feature development by integrating the core functions listed above for free on every API endpoint. It also deploys as a single binary requiring little operational investment.

Finally, we wanted a platform that we could invest in as an organization, thus requiring it to be easy for other internal teams to understand and build on. Clutch offers an integrated and guided development model that makes feature development a straightforward process. In addition to the first-class backend features, Clutch’s frontend offers unique abstractions for state management and multi-step forms that make frontend development easier for infrastructure teams without a lot of JavaScript experience.

For a detailed analysis on how other tools measure up to Clutch, see [Comparisons to Other Tools](/docs/about/comparison).

## Features

### The 'Control Plane' Model
[Envoy Proxy](https://www.envoyproxy.io/) was created at Lyft. Today, it’s one of the most popular proxies, deployed at many large internet companies and continuously advancing the standard for cloud networking. Our team has learned a lot from maintaining Envoy alongside the larger community. One of the most popular topics discussed among Envoy users is the [state of control plane development](https://mattklein123.dev/2020/03/15/on-the-state-of-envoy-proxy-control-planes/), specifically how to systematically integrate a wide range of disparate components such that Envoy can effectively route and report on network traffic. This is directly analogous to Clutch, which integrates disparate infrastructure systems behind a uniform API.

Clutch adopts many of Envoy Proxy’s core patterns that emerged from years of work on network control planes. Like Envoy, Clutch is [configuration-driven](/docs/configuration), [schema-driven](/docs/development/api), and leverages a modular [extension-based architecture](/docs/about/architecture) to make it work for a wide variety of use cases without compromising maintainability. Extending Clutch does not require forks or rewrites, and custom code can easily be compiled into the application from a custom public or private external repository. These same patterns that enable organizations large and small with unique tech stacks to converge on a single proxy solution will hopefully enable similarly unique organizations to converge on Clutch as an infrastructure control plane.

### Safety and Security
Additionally, Clutch ships with [authentication and authorization components](/docs/advanced/auth). OpenID Connect (OIDC) authentication flows for single-sign on, resource-level role-based access control (RBAC) via static mapping, and automatic auditing of all actions with the ability to run additional sinks for output, e.g., a Slackbot.

![Slack message showing the action performed by a user via Clutch](https://user-images.githubusercontent.com/4712430/104761542-ae544a00-5728-11eb-9206-fad9fa91aec5.png)

Clutch also has features to mitigate the potential for accidents. Guardrails and heuristics normally documented in runbooks can be implemented programmatically. For example, we never allow a user to scale down a cluster more than 50% at a time since this behavior has historically led to accidental outages during normal maintenance. In the future, we plan to fetch CPU and other usage data to display alongside the cluster information, even going as far as to limit the lower bounds for scale down if we determine that it is likely to cause an outage. By implementing guardrails and heuristics directly into the tool, we avoid the need to rely solely on training and runbooks to prevent accidents.

### Deployment and Onboarding
Clutch ships as a single binary that contains the frontend and backend, making it trivial to deploy. Many changes can be achieved via configuration rather than recompiling a new binary.

Other systems that offer infrastructure lifecycle tooling require a complicated set of microservices or migration to an opinionated way of managing and deploying applications. Clutch is meant to complement existing systems rather than replace them.

### Frameworks and Components
Clutch is powered by a Go backend and React frontend. It provides full-featured frameworks for backend and frontend development. All components in Clutch are composable, allowing for partial use of the framework offerings or completely custom features.

This component-and-workflow architecture allows a developer with limited frontend experience to replace clunky tooling or command-line scripts with a clear and easy-to-use step-by-step UI in [under an hour of development time](/docs/development/feature).

Clutch’s frontend packages offer components to easily build step-by-step workflows with a consistent and seamless user experience, including:

- [DataLayout](/docs/about/architecture/#data-layout), a workflow-local state management component that handles user input and data from API calls.
- [Wizard](/docs/about/architecture#wizard), for presenting step-by-step forms to users.
- Custom Material UI components, for displaying rich information with minimal code in a consistent manner across workflows.

Clutch’s backend relies heavily on generated code from protobuf API definitions. Protobuf tooling also generates frontend clients which keeps the backend and frontend in sync as APIs evolve. Components on the backend include:

- [Modules](/docs/about/architecture/#modules), implementations of the code generated API stubs
- [Services](/docs/about/architecture/#services), for interacting with external sources of data
- [Middleware](/docs/about/architecture/#middleware), for inspecting request and response data and applying auditing, authorization, etc.
- [Resolvers](/docs/about/architecture/#resolvers), a common interface to locate resources based on free-form text search or structured queries

Resolvers are one Clutch abstraction we hope will make a big impact on the way features can be abstracted to multiple organizations. Resolvers are easily extended with custom resource location code, allowing operators to locate resources (such as K8s pods or EC2 instances) by the common name(s) that the organization is accustomed to rather than the terse canonical identifier. For example, if developers call their application `myService-staging`, it’s easy to add code that will interpret such a query as the structured elements `${application_name}-${environment}`. Furthermore, the frontend automatically generates user input forms from the backend definitions.

With one line of code on the frontend:

```jsx
<Resolver type="clutch.aws.ec2.v1.Instance" />
```

The following form is rendered:
![Form with multiple input fields for finding an instance such as ID and region](https://user-images.githubusercontent.com/4712430/104761860-1e62d000-5729-11eb-9e57-d84ee5c06c7f.png)

Configuring additional search dimensions on the backend will automatically reflect in the rendered form on the frontend.

## Clutch at Lyft
![Image showing an engineer having to look at many different tools](https://user-images.githubusercontent.com/4712430/104762004-481bf700-5729-11eb-9767-f7779d87f797.png)

Before Clutch, Lyft engineers relied on a hodgepodge of command line tools, web interfaces, and runbooks to perform simple tasks. The most common alerts at Lyft required as many as six different sources of information to resolve: the alert, other service dashboards, the runbook, other sources of documentation, the vendor console or scripts, and configuration settings. As Lyft scaled in terms of team, product, and stack, we realized that the tools did not keep pace. We had no path forward to solve these problems with the existing frameworks. This led to the development of the first iteration of Clutch.

![Image showing an engineer using Clutch to quickly respond to alerts](https://user-images.githubusercontent.com/4712430/104762068-6255d500-5729-11eb-8346-b9377602797a.png)

Over the past year, Clutch has enjoyed an incredible rate of internal adoption, both in usage and development. Thousands of otherwise risky operations related to infrastructure management have gone through Clutch, each one representing the potential for an accident or delayed incident mitigation causing loss of trust from our drivers and riders. Seven internal engineering teams are already planning to add new features by the end of 2020 at the time of this writing, with at least half of those targeting open source. Engineers (including our wonderful interns) have been able to develop meaningful functionality with limited guidance. Most importantly, we’re finally able to see a path forward for delivering our internal platform through a single pane of glass, making Lyft infrastructure a product that meets our customers’ needs rather than a patchwork collection of systems and tools.

We have received a lot of positive feedback internally, for example:

> “I’m so happy that this exists because otherwise I would still be waiting on the tab to load in the cloud provider’s console.”

More specifics on Clutch at Lyft can be found in the [Lyft Case Study](/docs/about/lyft-case-study) article.

## Roadmap

Throughout the journey of building Clutch, the product has evolved and our external and internal roadmaps now encompass the entirety of the developer experience at Lyft. Our long-term vision involves building a context-aware developer portal. Not only providing a set of tools to developers, but presenting the most relevant tools and information as soon as the user lands in the portal.

Upcoming features include:
- **Envoy UI**, giving users a real-time dashboard to monitor their distributed application’s network behavior and configuration.
- **Chaos testing**, integrating with Envoy to allow scheduled fault injection and squeeze testing with automatic abort criteria.
- **Auto remediation**, responding to alerts automatically with the appropriate Clutch action.
- **Safety enhancements**, including escalation capabilities, challenge modals, and two-phase approvals.
- **Additional infrastructure lifecycle management capabilities**, viewing the state of a cluster to find outliers, perform long-running maintenance tasks.
- **Service health dashboard**, providing feedback to developers on the state of their service (e.g. code coverage, cost, active incidents) using configurable reporting mechanisms.
- **Generalized configuration management**, allowing users to manage complex configuration via a guided UI or otherwise reflect changes in infrastructure as code declarations.
- **Topology mapping**, associating users with the services they own and showing them relevant data and tools on the landing page.

For more details on upcoming projects visit our [roadmap](/docs/about/roadmap).

## Community
Clutch has had a significant impact on Lyft’s developer experience, allowing the infrastructure and other engineering teams to deliver tooling as a polished product rather than an afterthought. Teams are super excited about contributing new features internally, and engineers at Lyft love using the platform:

> “The usability is super awesome and improves the overall quality of life for engineers at Lyft.”

**Join us!** We think Clutch has the same outstanding potential for other engineering teams both large and small. We look forward to working with the community on additional plugins and core features in order to build a rich ecosystem. Our goal is to give every team and every tech stack access to first-class cloud-native tooling and reduced cognitive load. All contributions are welcome, from ideas to implementations, and we’re happy to help you get off the ground with your first feature.

To learn more about Clutch, contribute, or follow along:
- Visit our website and documentation at https://clutch.sh.
- Check out the code at https://github.com/lyft/clutch.
- Join us on Slack [@lyftoss/clutch](https://join.slack.com/t/lyftoss/shared_invite/zt-casz6lz4-G7gOx1OhHfeMsZKFe1emSA).
- Follow us on Twitter [@clutchdotsh](https://twitter.com/clutchdotsh).
- Lyft is hiring engineers to work on Clutch. [Apply here](https://grnh.se/65fd04572us).

## Thanks
Clutch would not have been possible without the contributions and hard work from many engineers at Lyft including Sindhura Tokala, Matthew Gumport, Scarlett Perry, Ryan Cox, Shawna Monero, Bill Gallagher, Gastón Kleiman, Paul Dittamo, Mike Cutalo, Rafal Augustyniak, Alan Chiu, Kathan Shah, Ansu Kar, Tony Allen, Amy Baer, Stephan Mercatoris, Susan Li, Jose Nino, and Matt Klein. Special thanks to Patrick Sunday, Martin Conte Mac Donnell, Polly Peterson, Michael Rebello, and Pete Morelli. Thanks again to all of those directly involved in the project and those who have provided us with guidance and support throughout this process.

**We can’t wait to see where the open-source community will take the project. Onward!**