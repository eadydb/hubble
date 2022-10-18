package gv

import (
	"context"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Gvk group version kind
type Gvk interface {
	ListNode() *corev1.NodeList
	ListService() *corev1.ServiceList
	ListDeployment() *appv1.DeploymentList
	ListPodBySelector(selector string) *corev1.PodList
	ListPod() *corev1.PodList

	GetNode(name string) *corev1.Node
	GetService(name string) *corev1.Service
	GetDeployment(name string) *appv1.Deployment
	GetPod(name string) *corev1.Pod
}

type KubeGvK struct {
	Config    *rest.Config          // kube config
	Client    *kubernetes.Clientset // kube client
	Namespace string                // kube namespace
	Ctx       context.Context       // context
}

// NewDefaultGvk Initialization default gvk
func NewDefaultGvk(config *rest.Config, client *kubernetes.Clientset) *KubeGvK {
	return &KubeGvK{
		Config:    config,
		Client:    client,
		Namespace: "default",
		Ctx:       context.TODO(),
	}
}

// SetNamespace reset namespace
func (kg *KubeGvK) SetNamespace(namespace string) *KubeGvK {
	kg.Namespace = namespace
	return kg
}

// SetContext reset context
func (kg *KubeGvK) SetContext(ctx context.Context) *KubeGvK {
	kg.Ctx = ctx
	return kg
}

func (kube *KubeGvK) ListNode() *corev1.NodeList {
	nodes, err := kube.Client.CoreV1().Nodes().List(kube.Ctx, metav1.ListOptions{})
	handlerError(err)
	return nodes
}

func (kube *KubeGvK) ListService() *corev1.ServiceList {
	svc, err := kube.Client.CoreV1().Services(kube.Namespace).List(kube.Ctx, metav1.ListOptions{})
	handlerError(err)
	return svc
}

func (kube *KubeGvK) ListDeployment() *appv1.DeploymentList {
	deploys, err := kube.Client.AppsV1().Deployments(kube.Namespace).List(kube.Ctx, metav1.ListOptions{})
	handlerError(err)
	return deploys
}

func (kube *KubeGvK) ListPodBySelector(selector string) *corev1.PodList {
	pods, err := kube.Client.CoreV1().Pods(kube.Namespace).List(kube.Ctx, metav1.ListOptions{LabelSelector: selector})
	handlerError(err)
	return pods
}

func (kube *KubeGvK) ListPod() *corev1.PodList {
	pods, err := kube.Client.CoreV1().Pods(kube.Namespace).List(kube.Ctx, metav1.ListOptions{})
	handlerError(err)
	return pods
}

func (kube *KubeGvK) GetNode(name string) *corev1.Node {
	node, err := kube.Client.CoreV1().Nodes().Get(kube.Ctx, name, metav1.GetOptions{})
	handlerError(err)
	return node
}

func (kube *KubeGvK) GetService(name string) *corev1.Service {
	svc, err := kube.Client.CoreV1().Services(kube.Namespace).Get(kube.Ctx, name, metav1.GetOptions{})
	handlerError(err)
	return svc
}

func (kube *KubeGvK) GetDeployment(name string) *appv1.Deployment {
	deploy, err := kube.Client.AppsV1().Deployments(kube.Namespace).Get(kube.Ctx, name, metav1.GetOptions{})
	handlerError(err)
	return deploy
}

func (kube *KubeGvK) GetPod(name string) *corev1.Pod {
	pod, err := kube.Client.CoreV1().Pods(kube.Namespace).Get(kube.Ctx, name, metav1.GetOptions{})
	handlerError(err)
	return pod
}

func handlerError(err error) {
	if err == nil {
		panic("ERR_NOT_FOUNT expected")
	}
	if !errors.IsNotFound(err) {
		panic(err.Error())
	}
}
