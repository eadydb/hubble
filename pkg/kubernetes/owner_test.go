package kubernetes

import (
	"context"
	"testing"

	kubernetesclient "github.com/eadydb/hubble/pkg/kubernetes/client"
	"github.com/eadydb/hubble/pkg/testutil"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	fakekubeclientset "k8s.io/client-go/kubernetes/fake"
)

func mockClient(m kubernetes.Interface) func(string) (kubernetes.Interface, error) {
	return func(string) (kubernetes.Interface, error) {
		return m, nil
	}
}

func TestTopLevelOwnerKey(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{
				{
					Name: "rs",
					Kind: "ReplicaSet",
				},
			},
		},
	}

	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rs",
			Namespace: "ns",
			OwnerReferences: []metav1.OwnerReference{
				{
					Name: "dep",
					Kind: "Deployment",
				},
			},
		},
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dep",
			Namespace: "ns",
		},
	}

	tests := []struct {
		description   string
		initialObject metav1.Object
		kind          string
		objects       []runtime.Object
		expected      string
	}{
		{
			description:   "owner is two levels up",
			initialObject: pod,
			kind:          "Pod",
			objects:       []runtime.Object{pod, rs, deployment},
			expected:      "Deployment-dep",
		}, {
			description:   "object is owner",
			initialObject: deployment,
			kind:          "Deployment",
			objects:       []runtime.Object{pod, rs, deployment},
			expected:      "Deployment-dep",
		}, {
			description:   "error, owner doesn't exist",
			initialObject: pod,
			kind:          "Pod",
			objects:       []runtime.Object{pod, rs},
		},
	}

	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			client := fakekubeclientset.NewSimpleClientset(test.objects...)
			t.Override(&kubernetesclient.Client, mockClient(client))

			actual := TopLevelOwnerKey(context.Background(), test.initialObject, "", test.kind)

			t.CheckDeepEqual(test.expected, actual)
		})
	}
}

func TestOwnerMetaObject(t *testing.T) {
	tests := []struct {
		description string
		or          metav1.OwnerReference
		objects     []runtime.Object
		expected    metav1.Object
	}{
		{
			description: "getting a deployment",
			or: metav1.OwnerReference{
				Kind: "Deployment",
				Name: "dep",
			},
			objects: []runtime.Object{
				&v1.Service{},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dep",
						Namespace: "ns",
					},
				},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dep",
						Namespace: "ns2",
					},
				},
			},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dep",
					Namespace: "ns",
				},
			},
		}, {
			description: "getting a replica set",
			or: metav1.OwnerReference{
				Kind: "ReplicaSet",
				Name: "rs",
			},
			objects: []runtime.Object{
				&v1.Service{},
				&appsv1.ReplicaSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rs",
						Namespace: "ns",
					},
				},
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dep",
						Namespace: "ns2",
					},
				},
			},
			expected: &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rs",
					Namespace: "ns",
				},
			},
		}, {
			description: "getting a job",
			or: metav1.OwnerReference{
				Kind: "Job",
				Name: "job",
			},
			objects: []runtime.Object{
				&batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "job",
						Namespace: "ns",
					},
				},
			},
			expected: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job",
					Namespace: "ns",
				},
			},
		}, {
			description: "getting a cronjob",
			or: metav1.OwnerReference{
				Kind: "CronJob",
				Name: "cj",
			},
			objects: []runtime.Object{
				&batchv1beta1.CronJob{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cj",
						Namespace: "ns",
					},
				},
			},
			expected: &batchv1beta1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cj",
					Namespace: "ns",
				},
			},
		}, {
			description: "getting a statefulset",
			or: metav1.OwnerReference{
				Kind: "StatefulSet",
				Name: "ss",
			},
			objects: []runtime.Object{
				&appsv1.StatefulSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ss",
						Namespace: "ns",
					},
				},
			},
			expected: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ss",
					Namespace: "ns",
				},
			},
		}, {
			description: "getting a replicationcontroller",
			or: metav1.OwnerReference{
				Kind: "ReplicationController",
				Name: "rc",
			},
			objects: []runtime.Object{
				&v1.ReplicationController{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rc",
						Namespace: "ns",
					},
				},
			},
			expected: &v1.ReplicationController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rc",
					Namespace: "ns",
				},
			},
		}, {
			description: "getting a pod",
			or: metav1.OwnerReference{
				Kind: "Pod",
				Name: "po",
			},
			objects: []runtime.Object{
				&v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "po",
						Namespace: "ns",
					},
				},
			},
			expected: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "po",
					Namespace: "ns",
				},
			},
		},
	}

	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			client := fakekubeclientset.NewSimpleClientset(test.objects...)
			t.Override(&kubernetesclient.Client, mockClient(client))

			actual, err := ownerMetaObject(context.Background(), "ns", "", test.or)

			t.CheckNoError(err)
			t.CheckDeepEqual(test.expected, actual)
		})
	}
}
