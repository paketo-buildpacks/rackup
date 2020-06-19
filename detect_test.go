package rackup_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-community/rackup"
	"github.com/paketo-community/rackup/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir        string
		detect            packit.DetectFunc
		gemfileLockParser *fakes.GemParser
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workingDir, "Gemfile"), []byte{}, 0644)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workingDir, "config.ru"), []byte{}, 0644)
		Expect(err).NotTo(HaveOccurred())

		gemfileLockParser = &fakes.GemParser{}

		detect = rackup.Detect(gemfileLockParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the Gemfile.lock specifies rack", func() {
		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(gemfileLockParser.ParseCall.Receives.Path).To(Equal(workingDir))
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "gems",
						Metadata: rackup.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "bundler",
						Metadata: rackup.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "mri",
						Metadata: rackup.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			}))
		})
	})

	context("when there is no Gemfile.lock and a config.ru file exists", func() {
		it.Before(func() {
			gemfileLockParser.ParseCall.Returns.RackFound = false
		})

		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(gemfileLockParser.ParseCall.Receives.Path).To(Equal(workingDir))
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "gems",
						Metadata: rackup.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "bundler",
						Metadata: rackup.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "mri",
						Metadata: rackup.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			}))
		})
	})

	context("when the workingDir does not have a config.ru and rack is not somehow specified", func() {
		it.Before(func() {
			gemfileLockParser.ParseCall.Returns.RackFound = false
			Expect(os.Remove(filepath.Join(workingDir, "config.ru"))).To(Succeed())
		})

		it("detect should fail with error", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail))
			Expect(gemfileLockParser.ParseCall.Receives.Path).To(Equal(workingDir))
		})
	})

	context("failure cases", func() {
		context("when unable to stat config.ru", func() {
			it.Before(func() {
				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to stat config.ru")))
			})
		})

		context("when parsing Gemfile.lock fails", func() {
			it.Before(func() {
				gemfileLockParser.ParseCall.Returns.Err = errors.New("failed to parse Gemfile.lock")
			})

			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError("failed to parse Gemfile.lock"))
			})
		})
	})
}
