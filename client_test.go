package client_test

import (
    "errors"
    "net/http"

    "github.com/pivotal-cf/app-automator-cf-client"
    "github.com/pivotal-cf/app-automator-cf-client/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/ginkgo/extensions/table"
    . "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
    Describe("Scale()", func() {
        It("uses TryWithRefresh", func() {
            cache := &mockAppGuidCache{}
            c := client.Client{
                Oauth: &mockOauth{},
                Capi: &mockCapi{
                    apps: []models.App{{Guid: "app-guid", Name: "app-name"}},
                },
                AppGuidCache: cache,
            }
            Expect(c.Scale("app-name", 1)).To(Succeed())
            Expect(cache.called).To(BeTrue())
        })

        DescribeTable("errors", func(modify func(*mockCapi, *mockAppGuidCache)) {
            capi := &mockCapi{
                apps: []models.App{{Guid: "app-guid"}},
            }
            cache := &mockAppGuidCache{}
            modify(capi, cache)

            c := client.Client{
                Oauth:        &mockOauth{},
                Capi:         capi,
                AppGuidCache: cache,
            }

            Expect(c.Scale("lemons", 2)).ToNot(Succeed())
        },
            Entry("TryWithRefresh returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                cache.tryErr = errors.New("expected")
            }),
            Entry("scale returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                capi.scaleErr = errors.New("expected")
            }),
        )
    })

    Describe("Process()", func() {
        It("uses TryWithRefresh", func() {
            cache := &mockAppGuidCache{}
            c := client.Client{
                Oauth: &mockOauth{},
                Capi: &mockCapi{
                    apps: []models.App{{Guid: "app-guid", Name: "app-name"}},
                },
                AppGuidCache: cache,
            }
            _, err := c.Process("app-name", "web")
            Expect(err).ToNot(HaveOccurred())
            Expect(cache.called).To(BeTrue())
        })

        DescribeTable("errors", func(modify func(*mockCapi, *mockAppGuidCache)) {
            capi := &mockCapi{
                apps:    []models.App{{Guid: "app-guid"}},
                process: models.Process{Instances: 2},
            }
            cache := &mockAppGuidCache{}
            modify(capi, cache)

            c := client.Client{
                Oauth:        &mockOauth{},
                Capi:         capi,
                AppGuidCache: cache,
            }

            _, err := c.Process("lemons", "web")
            Expect(err).To(HaveOccurred())
        },
            Entry("TryWithRefresh returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                cache.tryErr = errors.New("expected")
            }),
            Entry("process returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                capi.processErr = errors.New("expected")
            }),
        )
    })

    Describe("CreateTask()", func() {
        It("uses TryWithRefresh", func() {
            cache := &mockAppGuidCache{}
            c := client.Client{
                Oauth: &mockOauth{},
                Capi: &mockCapi{
                    apps: []models.App{{Guid: "app-guid", Name: "app-name"}},
                },
                AppGuidCache: cache,
            }
            task, err := c.CreateTask("app-name", "echo test", models.TaskConfig{})
            Expect(err).ToNot(HaveOccurred())
            Expect(cache.called).To(BeTrue())
            Expect(task.Guid).To(Equal("task-guid"))
        })

        It("sets a reasonable task name if not provided", func() {
            capi := &mockCapi{
                apps: []models.App{{Guid: "app-guid"}},
            }
            c := client.Client{
                Oauth:        &mockOauth{},
                Capi:         capi,
                AppGuidCache: &mockAppGuidCache{},
            }

            _, err := c.CreateTask("app-guid", "echo test", models.TaskConfig{
                DiskInMB:    10,
                MemoryInMB:  10,
                DropletGUID: "droplet-guid",
            })
            Expect(err).ToNot(HaveOccurred())
            Expect(capi.taskCfg.Name).To(Equal("echo test"))
        })

        It("passes on header options", func() {
            capi := &mockCapi{
                apps: []models.App{{Guid: "app-guid"}},
            }
            c := client.Client{
                Oauth:        &mockOauth{},
                Capi:         capi,
                AppGuidCache: &mockAppGuidCache{},
            }

            var headerOptionUsed bool
            opt := func(header *http.Header) {
                headerOptionUsed = true
            }
            _, err := c.CreateTask("app-guid", "echo test", models.TaskConfig{}, opt)
            Expect(err).ToNot(HaveOccurred())
            Expect(headerOptionUsed).To(BeTrue())
        })

        DescribeTable("errors", func(modify func(*mockCapi, *mockAppGuidCache)) {
            capi := &mockCapi{
                apps:    []models.App{{Guid: "app-guid"}},
                process: models.Process{Instances: 2},
            }
            cache := &mockAppGuidCache{}
            modify(capi, cache)

            c := client.Client{
                Oauth:        &mockOauth{},
                Capi:         capi,
                AppGuidCache: cache,
            }

            _, err := c.CreateTask("lemons", "command", models.TaskConfig{})
            Expect(err).To(HaveOccurred())
        },
            Entry("TryWithRefresh returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                cache.tryErr = errors.New("expected")
            }),
            Entry("create task returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                capi.taskErr = errors.New("expected")
            }),
        )
    })

    Describe("Stop()", func() {
        It("stops the app", func() {
            cache := &mockAppGuidCache{}
            c := client.Client{
                Oauth: &mockOauth{},
                Capi: &mockCapi{
                    apps: []models.App{{Guid: "app-guid", Name: "app-name"}},
                },
                AppGuidCache: cache,
            }
            Expect(c.Stop("app-name")).To(Succeed())
            Expect(cache.called).To(BeTrue())
        })

        DescribeTable("errors", func(modify func(*mockCapi, *mockAppGuidCache)) {
            capi := &mockCapi{
                apps:    []models.App{{Guid: "app-guid"}},
                process: models.Process{Instances: 2},
            }
            cache := &mockAppGuidCache{}
            modify(capi, cache)

            c := client.Client{
                Oauth:        &mockOauth{},
                Capi:         capi,
                AppGuidCache: cache,
            }

            err := c.Stop("lemons")
            Expect(err).To(HaveOccurred())
        },
            Entry("TryWithRefresh returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                cache.tryErr = errors.New("expected")
            }),
            Entry("stop returns an error", func(capi *mockCapi, cache *mockAppGuidCache) {
                capi.stopErr = errors.New("expected")
            }),
        )
    })
})

type mockOauth struct {
    err error
}

func (o *mockOauth) Token() (string, error) {
    return "bearer token", o.err
}

type mockCapi struct {
    apps    []models.App
    appsErr error

    process    models.Process
    processErr error

    scaleErr error
    stopErr  error
    taskErr  error

    taskCfg models.TaskConfig
}

func (c *mockCapi) Apps(query map[string]string) ([]models.App, error) {
    return c.apps, c.appsErr
}

func (c *mockCapi) Process(appGuid, processType string) (models.Process, error) {
    return c.process, c.processErr
}

func (c *mockCapi) Scale(appGuid, processType string, instanceCount uint) error {
    return c.scaleErr
}

func (c *mockCapi) CreateTask(appGuid, command string, cfg models.TaskConfig, opts ...models.HeaderOption) (models.Task, error) {
    for _, o := range opts {
        o(&http.Header{})
    }
    c.taskCfg = cfg
    return models.Task{Guid: "task-guid"}, c.taskErr
}

func (c *mockCapi) Stop(appGuid string) error {
    return c.stopErr
}

type mockAppGuidCache struct {
    called bool
    tryErr error
}

func (c *mockAppGuidCache) TryWithRefresh(appName string, f func(appGuid string) error) error {
    c.called = true
    err := f("app-guid")
    if c.tryErr != nil {
        return c.tryErr
    }
    return err
}
