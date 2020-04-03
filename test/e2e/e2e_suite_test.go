package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	libimage "github.com/operator-framework/operator-registry/pkg/lib/image"
)

var (
	dockerUsername = os.Getenv("DOCKER_USERNAME")
	dockerPassword = os.Getenv("DOCKER_PASSWORD")

	registryHost      string
	bundleImage       string
	indexImage        string
	stopLocalRegistry func()
	loginToRegistry   func(containerTool string)
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	if dockerUsername == "" && dockerPassword == "" {
		By("using a local registry when credentials are not available")
		ctx, cancel := context.WithCancel(context.Background())
		stopLocalRegistry = cancel

		var err error
		registryHost, err = libimage.RunDockerRegistry(ctx, "")
		Expect(err).ToNot(HaveOccurred())
	} else {
		By("using quay.io when registry credentials are available")
		loginToRegistry = func(containerTool string) {
			registryHost = "quay.io"
			dockerlogin := exec.Command(containerTool, "login", "-u", dockerUsername, "-p", dockerPassword, registryHost)
			Expect(dockerlogin.Run()).To(Succeed(), "Error logging into %s", registryHost)
		}
	}

	bundleImage = registryHost + "/olmtest/e2e-bundle"
	indexImage = registryHost + "/olmtest/e2e-index:" + indexTag
})

var _ = AfterSuite(func() {
	if stopLocalRegistry != nil {
		By("cleaning up local registry")
		stopLocalRegistry()
	}
})
