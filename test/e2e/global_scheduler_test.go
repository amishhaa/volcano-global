package e2e

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Volcano Global Scheduler", func() {
	var (
		clientset *kubernetes.Clientset
		ctx       context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		Expect(err).NotTo(HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should create namespace successfully", func() {
		nsName := "volcano-global-e2e"

		_, err := clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: nsName,
			},
		}, metav1.CreateOptions{})

		if err != nil && !apierrors.IsAlreadyExists(err) {
			Fail(err.Error())
		}
	})
})
