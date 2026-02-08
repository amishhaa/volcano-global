package utils_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"volcano.sh/volcano-global/pkg/utils"
)

var _ = Describe("Cluster Ready E2E", func() {
	var (
		k8sClient client.Client
		ctx       context.Context
	)

	BeforeSuite(func() {
		ctx = context.Background()

		// Try in-cluster config first (for CI), fallback to local kubeconfig
		cfg, err := rest.InClusterConfig()
		if err != nil {
			cfg, err = rest.InClusterConfig()
			Expect(err).NotTo(HaveOccurred())
		}

		scheme := runtime.NewScheme()
		Expect(clusterv1alpha1.AddToScheme(scheme)).To(Succeed())

		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
		Expect(err).NotTo(HaveOccurred())
	})

	It("should report ready when a Karmada cluster is ready", func() {
		clusterName := "member1" // MUST exist in your cluster

		var cluster clusterv1alpha1.Cluster

		Eventually(func() bool {
			err := k8sClient.Get(ctx, client.ObjectKey{Name: clusterName}, &cluster)
			return err == nil
		}, 2*time.Minute, 5*time.Second).Should(BeTrue())

		ready, msg := utils.CheckClusterReady(&cluster)

		Expect(ready).To(BeTrue())
		Expect(msg).To(BeEmpty())
	})
})
