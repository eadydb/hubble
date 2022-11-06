package kubernetes

import (
	"context"
	"fmt"

	kubernetesclient "github.com/eadydb/hubble/pkg/kubernetes/client"
	"github.com/eadydb/hubble/pkg/output/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TopLevelOwnerKey returns a key associated with the top level
// owner of a Kubernetes resource in the form Kind-Name
func TopLevelOwnerKey(ctx context.Context, obj metav1.Object, kubeContext string, kind string) string {
	for {
		or := obj.GetOwnerReferences()
		if or == nil {
			return fmt.Sprintf("%s-%s", kind, obj.GetName())
		}
		var err error
		kind = or[0].Kind
		obj, err = ownerMetaObject(ctx, obj.GetNamespace(), kubeContext, or[0])
		if err != nil {
			log.Entry(ctx).Warnf("unable to get owner from reference: %v", err)
			return ""
		}
	}
}

func ownerMetaObject(ctx context.Context, ns string, kubeContext string, owner metav1.OwnerReference) (metav1.Object, error) {
	client, err := kubernetesclient.Client(kubeContext)
	if err != nil {
		return nil, err
	}
	switch owner.Kind {
	case "Deployment":
		return client.AppsV1().Deployments(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	case "StatefulSet":
		return client.AppsV1().StatefulSets(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	case "ReplicaSet":
		return client.AppsV1().ReplicaSets(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	case "Job":
		return client.BatchV1().Jobs(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	case "CronJob":
		return client.BatchV1beta1().CronJobs(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	case "ReplicationController":
		return client.CoreV1().ReplicationControllers(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	case "Pod":
		return client.CoreV1().Pods(ns).Get(ctx, owner.Name, metav1.GetOptions{})
	default:
		return nil, fmt.Errorf("kind %s is not supported", owner.Kind)

	}
}
