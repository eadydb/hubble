package testutil

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// SetupFakeKubernetesContext replaces the current Kubernetes configuration
// file to setup a fixed current context.
func (t *T) SetupFakeKubernetesContext(config api.Config) {
	kubeConfig := t.TempFile("config", []byte{})

	if err := clientcmd.WriteToFile(config, kubeConfig); err != nil {
		t.Fatalf("writing temp kubeconfig")
	}

	t.Setenv("KUBECONFIG", kubeConfig)
}
