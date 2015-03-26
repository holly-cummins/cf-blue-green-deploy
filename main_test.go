package main_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	plugin "hub.jazz.net/git/bluemixgarage/cf-blue-green-deploy"

	"github.com/cloudfoundry/cli/plugin/fakes"
)

var _ = Describe("Main", func() {
	Describe("Plugin", func() {
		Describe("blue-green-deploy", func() {
			It("exists", func() {
				connection := &fakes.FakeCliConnection{}
				p := plugin.BlueGreenDeploymentPlugin{}
				p.Run(connection, []string{})
			})

			Describe("OldAppVersionList", func() {
				It("returns error", func() {
					connection := &fakes.FakeCliConnection{}
					connection.CliCommandWithoutTerminalOutputStub = func(args ...string) ([]string, error) {
						return nil, errors.New("Failed retrieving app names")
					}
					p := plugin.BlueGreenDeploymentPlugin{Connection: connection}
					_, err := p.OldAppVersionList("app-name")
					Expect(err).To(HaveOccurred())
				})

				It("returns list of application names", func() {
					connection := &fakes.FakeCliConnection{}
					connection.CliCommandWithoutTerminalOutputStub = func(args ...string) ([]string, error) {
						return []string{
								"Getting apps in org garage@uk.ibm.com / space dev as garage@uk.ibm.com...",
								"OK",
								"",
								"name                  					requested state   instances   memory   disk   urls",
								"app-name-20150326120000    		started           1/1         32M      1G",
								"app-name-20150326110000-old    started           1/1         32M      1G",
							},
							nil
					}
					p := plugin.BlueGreenDeploymentPlugin{Connection: connection}
					appList, _ := p.OldAppVersionList("app-name")

					Expect(appList).To(Equal([]string{"app-name-20150326110000-old"}))
				})
			})
		})

		Describe("DeleteApps", func() {
			It("deletes all apps and mapped routes in list", func() {
				connection := &fakes.FakeCliConnection{}
				p := plugin.BlueGreenDeploymentPlugin{Connection: connection}
				p.DeleteApps([]string{"app1", "app2"})

				Expect(strings.Join(connection.CliCommandArgsForCall(0), " ")).To(Equal("delete app1 -f -r"))
				Expect(strings.Join(connection.CliCommandArgsForCall(1), " ")).To(Equal("delete app2 -f -r"))
			})
		})

		Describe("DeleteOldAppVersions", func() {
			Context("when getting old app versions fails", func() {
				It("returns error", func() {
					connection := &fakes.FakeCliConnection{}
					connection.CliCommandWithoutTerminalOutputStub = func(args ...string) ([]string, error) {
						return nil, errors.New("Failed retrieving app names")
					}
					p := plugin.BlueGreenDeploymentPlugin{Connection: connection}
					Expect(p.DeleteOldAppVersions("app-name")).To(MatchError("Failed retrieving app names"))
				})
			})

			Context("when getting old app versions succeeds", func() {
				It("deletes all old app versions", func() {
					connection := &fakes.FakeCliConnection{}
					connection.CliCommandWithoutTerminalOutputStub = func(args ...string) ([]string, error) {
						return []string{
								"Getting apps in org garage@uk.ibm.com / space dev as garage@uk.ibm.com...",
								"OK",
								"",
								"name                  					requested state   instances   memory   disk   urls",
								"app-name-20150326120000    		started           1/1         32M      1G",
								"app-name-20150326110000-old    started           1/1         32M      1G",
							},
							nil
					}
					p := plugin.BlueGreenDeploymentPlugin{Connection: connection}
					p.DeleteOldAppVersions("app-name")
					Expect(strings.Join(connection.CliCommandArgsForCall(0), " ")).To(Equal("delete app-name-20150326110000-old -f -r"))
				})
			})
		})
	})
})
