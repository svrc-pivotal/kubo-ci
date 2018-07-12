package cidrs_test

import (
	"net"
	"tests/config"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	client_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var _ = Describe("Custom CIDRs", func() {
	var (
		k8s           kubernetes.Interface
		testconfig    *config.Config
		testNamespace string
		err           error
		svcController client_v1.ServiceInterface
	)

	BeforeEach(func() {
		k8s, err = NewKubeClient()
		Expect(err).NotTo(HaveOccurred())

		testNamespace = "test-" + GenerateRandomUUID()
		_, err = k8s.CoreV1().Namespaces().Create(&v1.Namespace{
			ObjectMeta: meta_v1.ObjectMeta{Name: testNamespace},
		})
		Expect(err).NotTo(HaveOccurred())

		svcController = k8s.CoreV1().Services(testNamespace)

		testconfig, err = config.InitConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		k8s.CoreV1().Namespaces().Delete(testNamespace, &meta_v1.DeleteOptions{})
	})

	Context("Services", func() {
		It("configures Kubernetes API server to the provided IP", func() {
			service, err := k8s.CoreV1().Services("default").Get("kubernetes", meta_v1.GetOptions{})

			Expect(err).NotTo(HaveOccurred())
			Expect(service.Spec.ClusterIP).To(Equal(testconfig.Kubernetes.KubernetesServiceIP))
		})

		It("creates service in the specified CIDR", func() {
			svcSpec := v1.Service{
				ObjectMeta: meta_v1.ObjectMeta{Name: "foo"},
				Spec:       v1.ServiceSpec{Ports: []v1.ServicePort{{Protocol: v1.ProtocolTCP, Port: 80}}},
			}
			svc, err := svcController.Create(&svcSpec)
			defer svcController.Delete("foo", &meta_v1.DeleteOptions{})

			Expect(err).NotTo(HaveOccurred())
			_, subnet, _ := net.ParseCIDR(testconfig.Kubernetes.ClusterIPRange)
			Expect(subnet.Contains(net.ParseIP(svc.Spec.ClusterIP))).To(BeTrue())
		})

		It("configures Kube-DNS to the provided IP", func() {
			service, err := k8s.CoreV1().Services("kube-system").Get("kube-dns", meta_v1.GetOptions{})

			Expect(err).NotTo(HaveOccurred())
			Expect(service.Spec.ClusterIP).To(Equal(testconfig.Kubernetes.KubeDNSIP))
		})
	})
})
