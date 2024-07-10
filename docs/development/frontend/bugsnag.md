---
title: BugSnag
{{ .EditURL }}
---

The frontend code in Clutch allows a user to utilize BugSnag for their bug catching abilities. The setup when using a Gateway or the normal clutch app is quite simple as it just requires the usage of environment variables.

An example is below

## Structure
```
frontend
├─ package.json
├─ .env.production
├─ ...
```

### Example Script in package.json
```json
"upload-sourcemaps": "yarn workspace @clutch-sh/tools uploadSourcemaps $PWD build/static .env.production --"
```

### Example .env.production
```bash
VITE_APP_SERVICE_NAME=<app>
VITE_APP_BUGSNAG_API_TOKEN=....
VITE_APP_BASE_URL=https://<app>.net
APPLICATION_ENV=production
```

And thats it, after setting up these few items Clutch will wrap the application in a BugSnag error boundary and report errors as it catches them.