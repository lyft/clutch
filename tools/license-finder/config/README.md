# Using License Finder

[License Finder](https://github.com/pivotal/LicenseFinder) scans third-party dependencies used by Clutch. If new dependencies are added in `./backend/go.mod`, the license scanner test may fail, outputing the unapproved dependency and its associated license. A project owner will need to approve the license or dependency for use.

# Approving licenses or dependencies

> Note: Commands below assume you're in the Clutch root directory and have docker installed.

```sh
# List dependencies that are not approved
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan/tools/license-finder && license_finder"

# Approve a license (any dependency with this license will be approved)
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan/tools/license-finder && license_finder permitted_licenses add '<license_to_add>'"

# Approve a dependency (approves just the specific dependency)
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan/tools/license-finder && license_finder approvals add '<dependency_to_add>'"
```

When a license or dependecy is approved, the changes will be automatically added to [`license_dependency_decisions.yml`](./license_dependency_decisions.yml). Commit the file changes to the pull request.

# Handling unknown licenses (preferred method)
When License Finder reports that a dependency's license is 'unknown', the license should be manually researched. Then record it with:

```sh
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan/tools/license-finder && license_finder licenses add '<dependency>' '<license>'"
```

Then you can procced with the command above to approve a dependency.

For additional usages not mentioned, please see [License Finder docs](https://github.com/pivotal/LicenseFinder).