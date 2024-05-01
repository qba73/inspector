package inspector_test

import (
	"context"
	"testing"

	"github.com/qba73/inspector"
	"k8s.io/apimachinery/pkg/version"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	fakediscovery "k8s.io/client-go/discovery/fake"
	testClient "k8s.io/client-go/kubernetes/fake"
)

func TestInspectorCollectsK8sVersion(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(),
	}
	got, err := c.ClusterVersion()
	if err != nil {
		t.Fatal(err)
	}
	want := "v1.29.2"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsK8sClusterID(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(kubeSystemNameSpace),
	}
	got, err := c.ClusterID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "421766aa-5d78-4c9e-8736-7faad1f2e927"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

// newTestClientset takes K8s runtime objects and returns a k8s fake clientset.
func newTestClientset(objects ...k8sruntime.Object) *testClient.Clientset {
	client := testClient.NewSimpleClientset(objects...)
	client.Discovery().(*fakediscovery.FakeDiscovery).FakedServerVersion = &version.Info{
		GitVersion: "v1.29.2",
	}
	return client
}

var (
	kubeSystemNameSpace = &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-system",
			UID:  "421766aa-5d78-4c9e-8736-7faad1f2e927",
		},
		Spec: corev1.NamespaceSpec{},
	}
)
