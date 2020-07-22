---
title: Workflows and Components
{{ .EditURL }}
---

import FrontendWorkflow from '@site/src/components/FrontendWorkflow';
import BackendComponent from '@site/src/components/BackendComponent';

For details on how each type of component integrates with Clutch please visit the [Architecture reference](/docs/about/architecture).

## Frontend

Frontend workflows are written using React and registered and configured in the [frontend config file](https://clutch.sh/docs/configuration#frontend) at build time.

### Workflow Packages
{{ range $p := .WorkflowPackages }}
<FrontendWorkflow packageName="{{ $p.PackageName }}" to="{{ $p.URL}}" workflows={ {{ toJson $p.Workflows | }} } />
{{ end }}

## Backend

Backend components are written in Go and registered and configured in the [backend config file](https://clutch.sh/docs/configuration#backend) at runtime.

### Services
{{ range $c := .Services }}
<BackendComponent name="{{ $c.Name }}" to="{{ $c.URL }}" desc="{{ $c.ClutchDoc.Description }}" />
{{ end }}

### Modules
{{ range $c := .Modules }}
<BackendComponent name="{{ $c.Name }}" to="{{ $c.URL }}" desc="{{ $c.ClutchDoc.Description }}" />
{{ end }}

### Resolvers
{{ range $c := .Resolvers }}
<BackendComponent name="{{ $c.Name }}" to="{{ $c.URL }}" desc="{{ $c.ClutchDoc.Description }}" />
{{ end }}

### Middleware
{{ range $c := .Middleware }}
<BackendComponent name="{{ $c.Name }}" to="{{ $c.URL }}" desc="{{ $c.ClutchDoc.Description }}" />
{{ end }}
