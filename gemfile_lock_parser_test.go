package rackup_test

import (
	"os"
	"testing"

	"github.com/paketo-buildpacks/rackup"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGemfileLockParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path   string
		parser rackup.GemfileLockParser
	)

	it.Before(func() {
		file, err := os.CreateTemp("", "Gemfile.lock")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		_, err = file.WriteString(`GEM
  remote: https://rubygems.org/
  specs:
    rack (1.5.2)
    rack-protection (1.5.2)
      rack
    sinatra (1.4.4)
      rack (~> 1.4)
      rack-protection (~> 1.4)
      tilt (~> 1.3, >= 1.3.4)
    tilt (1.4.1)

PLATFORMS
  ruby

DEPENDENCIES
	sinatra`)
		Expect(err).NotTo(HaveOccurred())

		path = file.Name()

		parser = rackup.NewGemfileLockParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	context("Parse", func() {
		it("parses the Gemfile.lock file to check for rack gem", func() {
			hasRack, err := parser.Parse(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(hasRack).To(Equal(true))
		})

		context("when the Gemfile.lock file does not exist", func() {
			it.Before(func() {
				Expect(os.Remove(path)).To(Succeed())
			})

			it("returns an ErrNotExist error", func() {
				_, err := parser.Parse(path)
				Expect(os.IsNotExist(err)).To(Equal(true))
			})
		})

		context("failure cases", func() {
			context("when the Gemfile.lock cannot be opened", func() {
				it.Before(func() {
					Expect(os.Chmod(path, 0000)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.Parse(path)
					Expect(err).To(MatchError(ContainSubstring("failed to parse Gemfile.lock:")))
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})
}
