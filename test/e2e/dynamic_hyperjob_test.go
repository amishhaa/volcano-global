package e2e

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Volcano Global Smoke Test", func() {
	It("Should create and delete a HyperJob via Dynamic Client", func() {
		home, _ := os.UserHomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		Expect(err).NotTo(HaveOccurred())

		client, err := dynamic.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		resource := schema.GroupVersionResource{
			Group:    "training.volcano.sh", 
			Version:  "v1alpha1", 
			Resource: "hyperjobs",
		}

		hj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "training.volcano.sh/v1alpha1",
				"kind":       "HyperJob",
				"metadata": map[string]interface{}{
					"name": "smoke-test-hj",
				},
				"spec": map[string]interface{}{},
			},
		}

		By("Step 1: Creating the HyperJob in kind-host")
		_, err = client.Resource(resource).Namespace("default").Create(context.TODO(), hj, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		By("Step 2: Verifying the HyperJob exists")
		Eventually(func() error {
			_, err := client.Resource(resource).Namespace("default").Get(context.TODO(), "smoke-test-hj", metav1.GetOptions{})
			return err
		}, "20s", "2s").Should(Succeed())

		By("Step 3: Cleaning up (Deleting the HyperJob)")
		err = client.Resource(resource).Namespace("default").Delete(context.TODO(), "smoke-test-hj", metav1.DeleteOptions{})
		Expect(err).NotTo(HaveOccurred())
	})
})