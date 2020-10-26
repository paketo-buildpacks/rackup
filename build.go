package rackup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

func Build(logger scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Writing start command")

		_, err := os.Stat(filepath.Join(context.WorkingDir, "config.ru"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to stat config.ru: %w", err)
		}

		// Use RACK_ENV=production since Rack v1.6.0+ defaults the host to local host in development mode (default)
		baseCommand := `bundle exec rackup --env RACK_ENV=production`

		// The order of precedence in setting the port is:
		// 1. the $PORT variable if it is set
		// 2. A port listed in config.ru with the -p flag
		// 3. 1 and 2 are not met, and the default of 9292 is used by running the start command with `config.ru`.
		command := fmt.Sprintf(`if [[ -z "${PORT}" ]]; then %s config.ru; else %s -p "${PORT}"; fi`, baseCommand, baseCommand)
		logger.Subprocess(command)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: command,
					},
				},
			},
		}, nil
	}
}
