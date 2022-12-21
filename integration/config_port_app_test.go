package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testConfigPortApp(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when building a simple app", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		context("a container port is specified via config.ru file only", func() {
			it("creates a working OCI image with a rackup start command and the port set in config.ru", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "config_port_app"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(
						settings.Buildpacks.MRI.Online,
						settings.Buildpacks.Bundler.Online,
						settings.Buildpacks.BundleInstall.Online,
						settings.Buildpacks.Rackup.Online,
					).
					WithEnv(map[string]string{"BP_LOG_LEVEL": "DEBUG"}).
					WithPullPolicy("never").
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.
					WithPublish("3000").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(BeAvailable())
				Eventually(container).Should(Serve(ContainSubstring("Hello world!")).OnPort(3000))

				Expect(logs).To(ContainLines(
					MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, settings.Buildpack.Name)),
					"  Writing start command",
					"  Checking the config.ru file",
					"    config.ru specifies a port: 3000",
					"",
					"  Assigning launch processes:",
					`    web (default): bash -c bundle exec rackup --env production -p "${PORT:-3000}"`,
				))

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(
					MatchRegexp(`INFO  WEBrick::HTTPServer#start: pid=\d+ port=3000`),
				)
			})
		})

		context("a container port is specified via $PORT environment variable AND config.ru file", func() {
			it("creates a working OCI image with a rackup start command using $PORT as the port", func() {
				var err error
				source, err = occam.Source(filepath.Join("testdata", "config_port_app"))
				Expect(err).NotTo(HaveOccurred())

				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(
						settings.Buildpacks.MRI.Online,
						settings.Buildpacks.Bundler.Online,
						settings.Buildpacks.BundleInstall.Online,
						settings.Buildpacks.Rackup.Online,
					).
					WithPullPolicy("never").
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8088"}).
					WithPublish("8088").
					WithPublish("3000").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(BeAvailable())

				_, exists := container.Ports["8088"]
				Expect(exists).To(BeTrue())
				Eventually(container).Should(Serve(ContainSubstring("Hello world!")).OnPort(8088))

				Expect(logs).To(ContainLines(
					MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, settings.Buildpack.Name)),
					"  Writing start command",
					"  Assigning launch processes:",
					`    web (default): bash -c bundle exec rackup --env production -p "${PORT:-3000}"`,
				))

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(
					MatchRegexp(`INFO  WEBrick::HTTPServer#start: pid=\d+ port=8088`),
				)
			})
		})
	})
}
