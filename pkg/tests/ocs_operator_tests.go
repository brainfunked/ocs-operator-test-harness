package tests

import (
	"github.com/brainfunked/ocs-operator-test-harness/pkg/metadata"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Prow Operator Tests", func() {
	defer ginkgo.GinkgoRecover()
	config, err := rest.InClusterConfig()

	if err != nil {
		panic(err)
	}

	ginkgo.It("Deployment ocs-operator exists in openshift-storage namespace", func() {
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		// Make sure the Deployment exists
		_, err = clientset.AppsV1().Deployments("openshift-storage").Get("ocs-operator", metav1.GetOptions{})

		if err != nil {
			metadata.Instance.FoundDeployment = false
		} else {
			metadata.Instance.FoundDeployment = true
		}

		Expect(err).NotTo(HaveOccurred())
	})
})
