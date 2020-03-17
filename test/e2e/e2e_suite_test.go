package e2e_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	dockerconfig "github.com/docker/distribution/configuration"
	dockerregistry "github.com/docker/distribution/registry"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	imageRegistry     *dockerregistry.Registry
	imageRegistryPort int
	imageRegistryHost string

	rootDockerDir = "testdata"
)

func init() {
	var err error
	// imageRegistryPort, err = freeport.GetFreePort()
	// if err != nil {
	// 	panic(err)
	// }
	imageRegistryPort = 5000
	imageRegistryHost = fmt.Sprintf("localhost:%d", imageRegistryPort)

	config := &dockerconfig.Configuration{}
	config.HTTP.Addr = fmt.Sprintf(":%d", imageRegistryPort)
	config.HTTP.DrainTimeout = time.Duration(10) * time.Second
	config.Log.Level = "debug"
	config.Storage = map[string]dockerconfig.Parameters{"filesystem": map[string]interface{}{
		"rootdirectory": rootDockerDir,
	}}
	config.HTTP.Net = "unix"
	// config.Storage = map[string]dockerconfig.Parameters{"inmemory": map[string]interface{}{}}

	imageRegistry, err = dockerregistry.NewRegistry(context.Background(), config)
	// Expect(err).NotTo(HaveOccurred())
	if err != nil {
		panic(err)
	}

	if err := imageRegistry.ListenAndServe(); err != nil {
		panic(err)
	}
}

func TestOperatorRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

// var _ = BeforeSuite(func() {
// 	config := &dockerconfig.Configuration{}
// 	config.HTTP.Addr = fmt.Sprintf(":%d", imageRegistryPort)
// 	config.HTTP.DrainTimeout = time.Duration(10) * time.Second
// 	config.Log.Level = "debug"
// 	// config.Storage = map[string]dockerconfig.Parameters{"filesystem": map[string]interface{}{
// 	// 	"rootdirectory": rootDockerDir,
// 	// }}
// 	config.Storage = map[string]dockerconfig.Parameters{"inmemory": map[string]interface{}{}}

// 	var err error
// 	imageRegistry, err = dockerregistry.NewRegistry(context.Background(), config)
// 	Expect(err).NotTo(HaveOccurred())

// 	go func() {
// 		if err := imageRegistry.ListenAndServe(); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	time.Sleep(10 * time.Minute)
// })

var _ = AfterSuite(func() {

})
