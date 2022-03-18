# Rackup Cloud Native Buildpack

## `gcr.io/paketo-buildpacks/rackup`

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

## Logging Configurations

To configure the level of log output from the **buildpack itself**, set the
`$BP_LOG_LEVEL` environment variable at build time either directly (ex. `pack
build my-app --env BP_LOG_LEVEL=DEBUG`) or through a [`project.toml`
file](https://github.com/buildpacks/spec/blob/main/extensions/project-descriptor.md)
If no value is set, the default value of `INFO` will be used.

The options for this setting are:
- `INFO`: (Default) log information about the progress of the build process
- `DEBUG`: log debugging information about the progress of the build process

```shell
$BP_LOG_LEVEL="DEBUG"
```
