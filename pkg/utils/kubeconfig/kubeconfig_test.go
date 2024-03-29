package kubeconfig

import (
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGetRecommendedKubeconfigPath(t *testing.T) {
	gotKubeconfigPath := GetRecommendedKubeconfigPath()
	wantKubeconfigPath := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	if gotKubeconfigPath != wantKubeconfigPath {
		t.Errorf("got %q, want %q", gotKubeconfigPath, wantKubeconfigPath)
	}
}

var testKubeconfig = `apiVersion: v1
clusters:
- cluster:
    server: http://127.0.0.1
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: ""
  name: test-cluster
current-context: test-cluster
kind: Config
preferences: {}
users: null
`

func TestAddContext(t *testing.T) {
	kubeconfigPath := "./test/kubeconfig"
	defer func() {
		_ = os.Remove(kubeconfigPath)
	}()
	clusterName := "test-cluster"
	err := AddContext(kubeconfigPath, clusterName, &Config{
		Cluster: &clientcmdapi.Cluster{
			Server: "http://127.0.0.1",
		},
		Context: &clientcmdapi.Context{
			Cluster: clusterName,
		},
	})
	if err != nil {
		t.Errorf("got %v, want nil", err)
	}

	want := testKubeconfig
	got, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		t.Errorf("failed to read kubeconfig file: %v", err)
	}
	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}

func TestRemoveContext(t *testing.T) {
	kubeconfigPath := "./test/kubeconfig"
	defer func() {
		_ = os.Remove(kubeconfigPath)
	}()
	_ = os.WriteFile(kubeconfigPath, []byte(testKubeconfig), 0o640)
	err := RemoveContext("./test/kubeconfig", "test-cluster")
	if err != nil {
		t.Errorf("failed to delete context: %v", err)
	}

	want := `apiVersion: v1
clusters: null
contexts: null
current-context: ""
kind: Config
preferences: {}
users: null
`
	got, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		t.Errorf("failed to read kubeconfig file: %v", err)
	}
	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}
