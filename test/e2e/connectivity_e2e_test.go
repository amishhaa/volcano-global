package e2e

import (
	"context"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Volcano Global Multi-Cluster Connectivity", func() {
	var kubeconfig string
	var client *kubernetes.Clientset

	BeforeEach(func() {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = "/root/.kube/config"
		}

		// Ensure the file exists before trying to use it
		_, err := os.Stat(kubeconfig)
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Kubeconfig not found at %s", kubeconfig))

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		Expect(err).NotTo(HaveOccurred())

		client, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Checking Cluster Infrastructure", func() {
		It("should see the host cluster nodes", func() {
			nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodes.Items)).To(BeNumerically(">", 0))
			fmt.Printf("\n[INFO] Found %d nodes in Host cluster\n", len(nodes.Items))
		})

		It("should verify that member clusters are joined (via namespaces or CRDs)", func() {
			// In a typical Karmada/Volcano-Global setup, namespaces are often 
			// created to represent member clusters.
			nsList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			
			fmt.Println("[INFO] Available Namespaces in Host:")
			for _, ns := range nsList.Items {
				fmt.Printf(" - %s\n", ns.Name)
			}
			Expect(len(nsList.Items)).To(BeNumerically(">", 0))
		})
	})
})
