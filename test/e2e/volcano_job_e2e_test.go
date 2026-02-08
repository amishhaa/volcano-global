package e2e

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	batchv1alpha1 "volcano.sh/apis/pkg/apis/batch/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Volcano Job E2E", func() {
	var (
		ctx        context.Context
		dynClient dynamic.Interface
		kube      *kubernetes.Clientset
		namespace string
	)

	BeforeEach(func() {
		ctx = context.Background()
		namespace = "volcano-e2e"

		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		Expect(err).NotTo(HaveOccurred())

		dynClient, err = dynamic.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		kube, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		// Create namespace if not exists
		_, err = kube.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		_ = kube.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	})

	It("should run a Volcano Job successfully", func() {
		jobGVR := schema.GroupVersionResource{
			Group:    batchv1alpha1.SchemeGroupVersion.Group,
			Version: batchv1alpha1.SchemeGroupVersion.Version,
			Resource: "jobs",
		}

		jobRes := dynClient.Resource(jobGVR).Namespace(namespace)

		job := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "batch.volcano.sh/v1alpha1",
				"kind":       "Job",
				"metadata": map[string]interface{}{
					"name": "pi-job",
				},
				"spec": map[string]interface{}{
					"minAvailable": int64(1),
					"schedulerName": "volcano",
					"tasks": []interface{}{
						map[string]interface{}{
							"name":     "pi",
							"replicas": int64(1),
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"restartPolicy": "Never",
									"containers": []interface{}{
										map[string]interface{}{
											"name":  "pi",
											"image": "perl",
											"command": []interface{}{
												"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err := jobRes.Create(ctx, job, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		By("waiting for job to complete")

		Eventually(func() string {
			obj, err := jobRes.Get(ctx, "pi-job", metav1.GetOptions{})
			if err != nil {
				return ""
			}

			phase, found, _ := unstructured.NestedString(obj.Object, "status", "state", "phase")
			if !found {
				return ""
			}
			return phase
		}, 5*time.Minute, 5*time.Second).Should(Equal("Completed"))
	})
})
