# Rackup Cloud Native Buildpack

## `gcr.io/paketo-community/rackup`

The Rackup CNB sets the start command for a given rack-compliant ruby application.

## Integration

This CNB writes a start command, so there's currently no scenario we can
imagine that you would need to require it as dependency. If a user likes to
include some other functionality, it can be done independent of the Rackup CNB
without requiring a dependency of it.

To package this buildpack for consumption:
```
$ ./scripts/package.sh
```
This builds the buildpack's source using GOOS=linux by default. You can supply another value as the first argument to package.sh.

## `buildpack.yml` Configurations

There are no extra configurations for this buildpack based on `buildpack.yml`.
