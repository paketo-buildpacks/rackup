package rackup

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

type BuildPlanMetadata struct {
	Launch bool `toml:"launch"`
}

//go:generate faux --interface GemParser --output fakes/gem_parser.go
type GemParser interface {
	Parse(path string) (rackFound bool, err error)
}

func Detect(parser GemParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		rackFound, err := parser.Parse(context.WorkingDir)
		if err != nil {
			return packit.DetectResult{}, err
		}

		if !rackFound {
			exists, err := fs.Exists(filepath.Join(context.WorkingDir, "config.ru"))
			if err != nil {
				return packit.DetectResult{}, fmt.Errorf("failed to stat 'config.ru': %w", err)
			}

			if !exists {
				return packit.DetectResult{}, packit.Fail.WithMessage("no 'config.ru' file found")
			}
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "gems",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "bundler",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "mri",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			},
		}, nil
	}
}
