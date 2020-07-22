---
title: What is Clutch?
{{ .EditURL }}
---

import Highlight from '@site/src/components/Highlight';

**<Highlight>Clutch</Highlight>** is an open source web UI and API platform designed to simplify, accelerate, and derisk common debugging, maintenance, and operational tasks.

Clutch provides everything you need to improve your developers' experience and operational capabilities. It comes with several out-of-the-box features for managing cloud-native infrastructure, but is easily configured or extended to interact with whatever you run, wherever you run it.

Stop putting your team through an endless stream of high-friction tools and user interfaces. Go beyond infrastructure as a platform to infrastructure as a product with happy engineers as your customer. Adopt Clutch!

<img alt="Clutch landing page" src={require('@docusaurus/useBaseUrl').default('img/docs/landing-page.png')} />

## Goals
- 🧰 **Simplify operations.**
  - Prevent accidents during rollout and maintenance by programatically catching mistakes and enforcing that configuration changes are safe to apply.
  - Lower mean-time-to-resolution (MTTR) during incidents that require manual intervention and information gathering.
  - Reduce training requirements for operators and documentation requirements for systems.
- 🕹️ **Optimize and integrate the developer experience.**
  - Deliver a straightforward and consistent user-experience so operators know exactly what is going to happen when they click a button.
  - YAML and command-line interfaces can only go so far. Clutch provides a straightforward web UI to simplify complex or rarely needed tasks.
  - Every company is unique, so Clutch is extensible to suit all needs and focuses on presenting only the required information for the task at hand.
  - Build a single pane of glass and minimize information sprawl across several tools, reducing cognitive load.
- 🔧 **Complement infrastructure-as-code.**
  - Live, runtime changes are inevitable in any system, give them a home with Clutch.
  - Provide guided workflows for complex configuration.

## Features
- 🔌 **Highly extensible.**
  - Extension points exist throughout the stack to allow custom integrations without rewrites.
  - Clutch is configuration-driven so it can be deployed and reconfigured for varied environments with ease.
  - Private extensions can be plugged-in without maintaining a fork.
- 🔍 **Built for discovery.**
  - Resources have many common names. Clutch's Resolver pattern makes it easier than ever to locate resources.
  - The Resolver provides server-generated forms with one-line of frontend code, ensuring the API and frontend are always in sync.
- ⚛️ **Easy to develop, run, and maintain.** 
  - Developed with Go and Typescript, plus Protobuf for generated interfaces throughout.
  - Back-end abstractions ensure loose coupling and put feature development on rails.
  - Frontend components make it simple for developers with limited frontend experience to ship features.
  - Deployable as a single binary containing both backend and frontend resources.
  - Basic auditing, authorization, stats, and logging come for free with every endpoint.
- 🔒 **Secure and observable.** 
  - Single sign-on support.
  - Role-based access control (RBAC) engine for granular access control beyond what vendor IAM policies support.
  - Built-in auditing with sinks for Slack and more.
  - Extensive logging and stats capabilities.

## Vision
- 💡 **Get smart(er).** 
  - Present relevant information on the homepage based on context such as user identity, ownership, and ongoing incidents.
  - Additional core integrations with cloud-native infrastructure.
  - Fully-customizable heuristic support. Safety is based on the current context and design of the overall system.
  - Auto-remediation and alert annotation. Consume events from a monitoring system and append relevant info or automatically solve the problem.
- 🎛 **More interfaces.**
  - Clutch is effectively an infrastructure control plane. Take advantage of the safety and integrations in Clutch beyond the core GUI.
  - Chat gateway, e.g. a Slackbot, for non-destructive tasks and status checks.
  - Command-line gateway, e.g. `clutchctl`, for those times when `grep` is just faster than anything else.
  - Code-generated SDK for polyglot automation support.

For more concrete items, see the [Roadmap](/docs/roadmap).

## Why Clutch?

Clutch was designed and built after careful study of the other available options for solving infrastructure user experience and developer experience issues.

Other tools that solve similar problems suffered from one or more of the following issues:
- Not easily extensible.
- Require extensive migration to an opinionated ecosystem.
- Difficult to deploy or maintain.
- Too limited in scope and not designed for extensibility.
- Have no concept of safety and security.

Clutch has patterned itself off of the successful model used by Envoy Proxy, also developed at Lyft, offering a rich, configuration-driven alternative to legacy tools. Deploying Clutch is also dead simple, with a single binary for both frontend and backend.
For a direct comparison between Clutch and other tools, see [Comparison to Other Tools](/docs/about/comparison).

## FAQ

This section will be updated as-needed with answers to common questions from the community. To reach out to the team directly, please see the available contact options in [Community](/docs/community).
