# Documentation

This directory stores the Clutch docs along with the code for the microsite.

This main components within the directory are:
- The microsite, which also hosts documentation in [_website](./_website).
- The table of contents in [sidebars.json](./sidebars.json).
- All of the markdown for documentation.

### Format and Rendering

Documentation is rendered using a combination of **MDX** and **Go templates**.

To learn more about Docusaurus-flavored **MDX**, see https://v2.docusaurus.io/docs/markdown-features/.
Pay particular attention to the [section on headers](https://v2.docusaurus.io/docs/markdown-features/#markdown-headers).

To learn about the capabilities of **Go templates**, see docs for [`text/template`](https://pkg.go.dev/text/template?tab=doc#hdr-Actions).
Templating is performed by [_website/generator/generate.go](./_website/generator/generate.go).

In the future the generator will pull documentation from code and API definitions dynamically. For now,
the documentation generator simply evaluates templates using the Go templating language with an empty
context. If you want to add additional context or callable functions to templates, you can do so in the
Go file linked above.


### Microsite

The microsite lives in [_website](./_website) and is built on [Docusaurus](https://v2.docusaurus.io/docs/introduction).

### Development

`make docs-dev` will start a live-reload server for documentation. Once the server is up, you will need to
run `make docs-generate` to refresh the documentation as you edit it. Any changes resulting from
generation should immediately appear on the dev server.
