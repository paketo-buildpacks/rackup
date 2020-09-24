package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/paketo-buildpacks/rackup"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)
	parser := rackup.NewGemfileLockParser()

	packit.Run(
		rackup.Detect(parser),
		rackup.Build(logger),
	)
}
