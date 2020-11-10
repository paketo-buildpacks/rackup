package rackup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

		// check if the config.ru file specifies a port
		configru, err := ioutil.ReadFile(filepath.Join(context.WorkingDir, "config.ru"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to read config.ru: %w", err)
		}
		configPort, err := regexp.MatchString(`#\\.*?(-p|--port) \d+`, string(configru))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to parse contents of config.ru: %w", err)
		}

		fallbackPort := "9292"

		// if config.ru specifies a port, just parse out the port number
		if configPort {
			reg, _ := regexp.Compile(`(-p|--port) \d+`)
			portString := reg.FindString(string(configru))
			fallbackPort = strings.TrimLeft(portString, "(--port|-p) ")
		}

		// Use RACK_ENV=production since Rack v1.6.0+ defaults the host to local host in development mode (default)
		// The order of precedence in setting the port is:
		// 1. the $PORT variable if it is set
		// 2. A port listed in config.ru with the -p or --port flag
		// 3. 1 and 2 are not met, and the fallback port is set to the default of 9292.
		command := fmt.Sprintf(`bundle exec rackup --env RACK_ENV=production -p "${PORT:-%s}"`, fallbackPort)
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
