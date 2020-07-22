# @clutch-sh/tools

This package exists to install development dependencies at the root of
clutch packages to be shared amongst all other packages.

## Dependencies

The dependencies defined in this projects `package.json` are listed as
production dependencies as opposed to development (even though they are
only used in development) to allow yarn to pull them upward.

Since this package is meant only for development this shouldn't be an
issue.

For more information see the [documentation](https://clutch.sh/docs/development/frontend#clutch-shtools).