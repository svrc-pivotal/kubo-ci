package oss_only_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubectl", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand(
			"create", "namespace", runner.Namespace()).Wait("60s")
	})

	AfterEach(func() {
		runner.RunKubectlCommand(
			"delete", "namespace", runner.Namespace()).Wait("60s")
	})

	It("Should have a ulimit of 65536", func() {
		podName := test_helpers.GenerateRandomUUID()
		output, err := runner.GetOutput("run", podName, "--image", "pcfkubo/ulimit", "--restart=Never", "-i", "--rm")
		Expect(err).NotTo(HaveOccurred())
		Expect(output[0]).To(Equal("65536"))
	})

})
