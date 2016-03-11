package main_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var redirectorPath string
var redirectorPort int
var redirectorAddress string

var redirectorRunner ifrit.Runner
var redirectorProcess ifrit.Process
var logger lager.Logger

func TestRedirector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redirector Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	redirector, err := gexec.Build("github.com/ematpl/redirector", "-race")
	Expect(err).NotTo(HaveOccurred())

	return []byte(redirector)
}, func(payload []byte) {
	redirectorPath = string(payload)
	redirectorPort = 6000 + GinkgoParallelNode()
	redirectorAddress = fmt.Sprintf("127.0.0.1:%d", uint16(redirectorPort))

	logger = lagertest.NewTestLogger("test")
})

var _ = BeforeEach(func() {
	redirectorRunner = newRedirectorRunner(redirectorPath, redirectorAddress)
	redirectorProcess = ifrit.Invoke(redirectorRunner)
})

var _ = AfterEach(func() {
	ginkgomon.Kill(redirectorProcess)
})

func newRedirectorRunner(bin, listenAddr string) ifrit.Runner {
	return ginkgomon.New(ginkgomon.Config{
		Name: "redirector",
		Command: exec.Command(
			bin,
			"-listenAddr", listenAddr,
		),
		StartCheck: "redirector.started",
	})
}
