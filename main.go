package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)
	parser := NewGemfileLockParser()

	detect := Detect(parser)
	build := Build(logger)

	packit.Run(detect, build)
}
