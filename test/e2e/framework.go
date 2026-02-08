package e2e

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string
var clientset *kubernetes.Clientset

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig")
}

func SetupCluster() error {
	if kubeconfig == "" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err = kubernetes.NewForConfig(cfg)
	return err
}

func TeardownCluster() {
	// cleanup if needed later
}

func GetClient() *kubernetes.Clientset {
	return clientset
}
