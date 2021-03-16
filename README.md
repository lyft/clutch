<h1 align="center">
  <a href="http://www.clutch.sh/">
    <img src="https://user-images.githubusercontent.com/1004789/86156525-f1b3d780-baba-11ea-88a3-51a7391cd310.png" alt="Clutch" width="100%">
  </a>
</h1>

<div align="center">
  <a href="http://www.clutch.sh/">
    <img src="https://user-images.githubusercontent.com/1004789/86159195-d054ea80-babe-11ea-997f-a309b1b43040.png" alt="Clutch" width="300">
  </a>
</div>
<h4 align="center">An extensible platform for infrastructure management</h4>

<p align="center">
  <a href="/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-Apache%202.0-blue" alt="License">
  </a>
  <a href="https://join.slack.com/t/lyftoss/shared_invite/zt-casz6lz4-G7gOx1OhHfeMsZKFe1emSA">
    <img alt="Slack" src="https://img.shields.io/badge/Slack-@lyftoss/clutch-lightgrey?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABHNCSVQICAgIfAhkiAAAAllJREFUOI2Fkk9IVFEUh79z732jZePEhJagi4pMhAor21Q7KSpcJBjVItpEywpCKBGKalMRtIwgokXtIsn+kFK0ECGhViEIUW5a+I8ZU57z3pt7WkyFPgMv3M09H989h98RUmf7x9nTucDeVWVqdH9uF8CO131jAjXFxfDsRNfd90t5lxbkAtsj1jUYYxp2Dy8+C2d78iZwLaiSraq+ASwTmLQgKfPBZKoR6xB8HZJ4vAdVRJA0v0Lw+WDu0kibFR+FzRP713QYWR/8Kwqa5l3Lx0JXNpAeVA2IAqiqiWM/No2c2Ujvil9bBq7ccWIPzvm5XlfreCjW5VG/DAoyVXv3DCeTYaFHIVt51ErHa9bVXPZxwvo498AgFEQESF2F2Ngx76NPElgksJRVvzW/6T2gXhEERecr7Y2Gb0F8W1TKGOCLC9QbO0h71W0ABvpuIbaeY9fOsUnh0bUhoMTR68ccwGT7gbBMJruJkQ4BJml8aPk5veHPOHF072do6vZF/a1DhlbNxMVI1C/UAPKj6dR4xrhtToT6709kovHk8+ogcxyFhYmou/bV6M61LtsXJcsDsFYolZLPTkTqPYrBVEZvkpbIl5FK6OeNo6CApgJUr4DkXVJcuEhtzYVfyeLrStbyLw5pUMkf/tY909/6UsVu/lsTwCc+lqRw020t9j+myON01gAkKAob5GsnK5cQ+M8mLsV0ysazb5uHwne7NRxs06kX266uKlC0LCLYyt7OAz4ue6LYI1LduarAz4f3E/VzJZ+Mb+HpifyR8UNxrN/j2M8Ynb2R5n8DJJL6gIHi/LUAAAAASUVORK5CYII%3D&style=social">
  </a>
  <a href="https://reactjs.org/">
    <img alt="react" src="https://img.shields.io/badge/React-grey?logo=react">
  </a>
  <a href="https://golang.org/">
    <img alt="golang" src="https://img.shields.io/badge/Go-grey?logo=go">
  </a>
</p>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#getting-started">Getting Started</a> •
  <a href="#extending-clutch">Extending Clutch</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#contributing">Contributing</a>
</p>

<div align="center">
  <img src="https://user-images.githubusercontent.com/39421794/104223887-39c09900-5412-11eb-9791-af10afdc6bbb.gif" width="75%" />
</div>

## Key Features

**Clutch** provides everything you need to simplify operations and in turn improve your developer experience and operational capabilities. It comes with several out-of-the-box features for managing cloud-native infrastructure, but is designed to be org-agnostic and easily taught how to find or interact with whatever you run, wherever you run it.

- 🔌 **Highly extensible.**
  - Extension points exist throughout the stack to allow custom integrations without rewrites.
  - Clutch is configuration-driven so it can be deployed and reconfigured for varied environments with ease.
  - Private extensions can be plugged-in without maintaining a fork.
- 🔎 **Built for discovery.**
  - Resources have many common names. Clutch's Resolver pattern makes it easier than ever to locate resources.
  - The Resolver provides server-generated forms with one-line of frontend code, ensuring the API and frontend are always in sync.
- 💻 **Easy to develop, run, and maintain.**
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

## Getting Started

So you want to run Clutch? That's great to hear! There are several supported methods of running Clutch, all
of which are outlined in our [Getting Started documentation](https://clutch.sh/docs/getting-started/build-guides) to learn how to run Clutch in Docker or build it locally.

Clutch also has a [mock server](https://clutch.sh/docs/getting-started/mock-gateway) for testing and developing features in isolation from the systems they depend on.

## Extending Clutch

Clutch ships with a default configuration and some out of the box workflows to make on-boarding as easy
as possible. However, there are lots of use cases for Clutch, some of which may not be written yet and others which are not broadly applicable.

To get started developing new features or functionality within Clutch check out the
[development guides](https://clutch.sh/docs/development/guide) on how to
develop each of the different pieces. While you're there, take a few additional minutes to read through the [configuration documentation](https://clutch.sh/docs/configuration). This allows you to override the default configuration that ships out of the box with Clutch.

## Documentation

Clutch has extensive documentation that can be found on our site [clutch.sh](https://clutch.sh/docs/).

If you're looking for the source of the hosted documentation both the content and the code
for the website live within the [docs/](docs/) directory.

## Contributing

Thinking of contributing back to Clutch? Awesome! We love and welcome all contributions.

First things first, please read over our [Code of Conduct](./CODE_OF_CONDUCT) and our
[guidelines](https://clutch.sh/docs/community#contributing) before opening a pull
request.

Want to contribute but not sure where to start? Check out the
[`good first issue`](https://github.com/lyft/clutch/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22) for tasks specifically scoped for those with less
familiarity.

All contributions require that the author has signed the [Lyft CLA](https://oss.lyft.com/cla/clas/1.0). Login to the [Lyft CLA Service](https://oss.lyft.com/cla) with GitHub to review and sign the CLA.
If you are contributing on behalf of an organization please reach out to clutch@lyft.com to have
your company added.
