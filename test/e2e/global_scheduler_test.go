package e2e

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = BeforeSuite(func() {
	err := SetupCluster()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	TeardownCluster()
})

var _ = Describe("Volcano Global Scheduler", func() {

	It("should create pods successfully", func() {
		client := GetClient()

		ns := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "volcano-global-e2e",
			},
		}
		_, err := client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			Fail(err.Error())
		}

		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: ns.Name,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "pause",
						Image: "k8s.gcr.io/pause:3.9",
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		}

		_, err = client.CoreV1().Pods(ns.Name).Create(context.TODO(), pod, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() v1.PodPhase {
			p, _ := client.CoreV1().Pods(ns.Name).Get(context.TODO(), pod.Name, metav1.GetOptions{})
			return p.Status.Phase
		}, 60*time.Second, 2*time.Second).Should(Equal(v1.PodRunning))
	})
})
