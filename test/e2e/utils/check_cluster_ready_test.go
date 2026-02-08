package utils

import (
	"testing"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"volcano.sh/volcano-global/pkg/utils"
)

func TestCheckClusterReady(t *testing.T) {
	tests := []struct {
		name        string
		cluster     *clusterv1alpha1.Cluster
		wantReady   bool
		wantMessage string
	}{
		{
			name: "cluster ready",
			cluster: &clusterv1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-1",
				},
				Status: clusterv1alpha1.ClusterStatus{
					Conditions: []metav1.Condition{
						{
							Type:   clusterv1alpha1.ClusterConditionReady,
							Status: metav1.ConditionTrue,
						},
					},
				},
			},
			wantReady:   true,
			wantMessage: "",
		},
		{
			name: "cluster not ready",
			cluster: &clusterv1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-2",
				},
				Status: clusterv1alpha1.ClusterStatus{
					Conditions: []metav1.Condition{
						{
							Type:    clusterv1alpha1.ClusterConditionReady,
							Status:  metav1.ConditionFalse,
							Reason:  "ConnectionFailed",
							Message: "cannot reach apiserver",
						},
					},
				},
			},
			wantReady:   false,
			wantMessage: "Cluster <cluster-2> is not ready, reason: ConnectionFailed, message: cannot reach apiserver",
		},
		{
			name: "cluster has no ready condition",
			cluster: &clusterv1alpha1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster-3",
				},
				Status: clusterv1alpha1.ClusterStatus{
					Conditions: []metav1.Condition{
						{
							Type:   "SomeOtherCondition",
							Status: metav1.ConditionTrue,
						},
					},
				},
			},
			wantReady:   false,
			wantMessage: "Cluster<cluster-3> has not Ready Condition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ready, msg := utils.CheckClusterReady(tt.cluster)

			if ready != tt.wantReady {
				t.Fatalf("expected ready=%v, got %v", tt.wantReady, ready)
			}

			if msg != tt.wantMessage {
				t.Fatalf("expected message=%q, got %q", tt.wantMessage, msg)
			}
		})
	}
}
