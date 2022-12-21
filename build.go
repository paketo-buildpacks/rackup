package rackup

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Writing start command")

		logger.Debug.Process("Checking the config.ru file")
		_, err := os.Stat(filepath.Join(context.WorkingDir, "config.ru"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to stat config.ru: %w", err)
		}

		// check if the config.ru file specifies a port
		configru, err := os.ReadFile(filepath.Join(context.WorkingDir, "config.ru"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to read config.ru: %w", err)
		}

		port := "9292"
		matches := regexp.MustCompile(`(?m)[\r\n]*^.*#\\.*?(-p|--port) (\d+).*$`).FindStringSubmatch(string(configru))
		if len(matches) == 3 {
			port = matches[2]
			logger.Debug.Subprocess("config.ru specifies a port: %s", port)

			content := strings.Replace(string(configru), matches[0], "", 1)
			err = os.WriteFile(filepath.Join(context.WorkingDir, "config.ru"), []byte(content), 0600)
			if err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to rewrite config.ru: %w", err)
			}
		}
		logger.Debug.Break()

		// Hardcode `--env production` since Rack v1.6.0+ defaults the host to local host in development mode (default)
		// The order of precedence in setting the port is:
		// 1. the $PORT variable if it is set
		// 2. A port listed in config.ru with the -p or --port flag
		// 3. 1 and 2 are not met, and the fallback port is set to the default of 9292.
		args := fmt.Sprintf(`bundle exec rackup --env %s -p "${PORT:-%s}"`, "production", port)
		processes := []packit.Process{
			{
				Type:    "web",
				Command: "bash",
				Args:    []string{"-c", args},
				Default: true,
				Direct:  true,
			},
		}
		logger.LaunchProcesses(processes)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
