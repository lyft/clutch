# Using License Finder

We run a [license scanner](https://github.com/pivotal/LicenseFinder) for third-party dependencies used by Clutch. If new dependencies are added in `./backend/go.mod`, the license scanner test may fail, outputing the unapproved dependency and its associated license. A project owner will need to approve the license or dependency for use ([instructions](#common-usage) below).

# Approving licenses or dependencies
When a license or dependecy is approved, the changes will be automatically added to [`license_dependency_decisions.yml`](./license_dependency_decisions.yml). The file changes will need to be commited to the pull request.

# Common Usage
> Note: Commands below assume you're in the Clutch repo root folder and have docker installed.
```sh
# List dependencies that are not approved
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan && license_finder"

# Approve a license (any dependency with this license will be approved)
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan && license_finder permitted_licenses add '<license_to_add>'"

# Approve a dependency (approves just the specific dependency)
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan && license_finder approvals add '<dependency_to_add>'"
```

# Handling unknown licenses (preferred method)
When `license_finder` reports that a dependency's license is 'unknown', the actual license should be manually researched. When the real license has been established, record it with:

```sh
$ docker run -v $PWD:/scan -it licensefinder/license_finder /bin/bash -lc "cd /scan && license_finder licenses add '<dependency>' '<license>'"
```

Then you can procced with commands above to approve the dependency.

For additional usages not mentioned, please see [License Finder docs](https://github.com/pivotal/LicenseFinder).