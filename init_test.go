package main_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitRackup(t *testing.T) {
	suite := spec.New("rackup", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("GemfileParser", testGemfileParser)
	suite.Run(t)
}
