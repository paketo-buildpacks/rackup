package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/rackup"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	parser := rackup.NewGemfileLockParser()

	packit.Run(
		rackup.Detect(parser),
		rackup.Build(logger),
	)
}
