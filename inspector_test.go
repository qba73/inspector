package inspector_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/inspector"
	"k8s.io/apimachinery/pkg/version"

	appsv1 "k8s.io/api/apps/v1"
	coordv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	fakediscovery "k8s.io/client-go/discovery/fake"
	testClient "k8s.io/client-go/kubernetes/fake"
)

func TestInspectorCollectsK8sVersion(t *testing.T) {
	t.Parallel()

	i := newTestInspector()
	got, err := i.ClusterVersion()
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

	c := newTestInspector(kubeSystemNameSpace)
	got, err := c.ClusterID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "421766aa-5d78-4c9e-8736-7faad1f2e927"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsNumberOfNodesInTheCluster(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		clusterNode1,
		clusterNode2,
		clusterNode3,
	)
	got, err := i.Nodes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := 3
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnAWSNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(kubeSystemNameSpace, nodeAWS)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "aws"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnAzureNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(kubeSystemNameSpace, nodeAzure)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "azure"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnGCPNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(kubeSystemNameSpace, nodeGCP)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "gce"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnKindNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(kubeSystemNameSpace, nodeKind)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "kind"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnVSphereNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeVSphere,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "vsphere"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnK3SNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeK3S,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "k3s"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnIBMCloudNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeIBMCloud,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "ibmcloud"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnIBMPowerNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeIBMPowerVS,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "ibmpowervs"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnCloudStackNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeCloudStack,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	want := "cloudstack"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnOpenStackNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeOpenStack,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "openstack"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnDigitalOceanNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeDigitalOcean,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "digitalocean"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnEquinixMetallNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeEquinixMetal,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "equinixmetal"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorCollectsPlatformNameOnAlibabaNode(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeAlibaba,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "alicloud"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorDeterminesUnknownPlatformOnMissingPlatformIDField(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeMalformedBlankPlatformID,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "unknown"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorDeterminesUnknownPlatformOnMalformedPlatformIDField(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeMalformedPlatformID,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "//4232e3c7-d83c-d72b-758c-71d07a3d9310"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorDeterminesUnknownPlatformOnMalformedBlankPlatformIDField(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeMalformedBlankPlatformID,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "unknown"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorDeterminesUnknownPlatformOnMalformedEmptyPlatformIDField(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nodeMalformedEmptyPlatformID,
	)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "unknown"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorDeterminesUnknownPlatformOnMalformedPartialPlatformIDField(t *testing.T) {
	t.Parallel()

	i := newTestInspector(kubeSystemNameSpace, nodeMalformedPartialPlatformID)
	got, err := i.Platform(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := "unknown"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInspectorListsPodsInNotExistingNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		nginxIngressNameSpace,
		pod1,
	)
	got, err := i.Pods(context.Background(), "notExistingNamespace")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.PodList{}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNotExistingPodsInDefaultNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(kubeSystemNameSpace, defaultNameSpace)
	got, err := i.Pods(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.PodList{}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsExistingPodsInDefaultNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		kubeSystemNameSpace,
		defaultNameSpace,
		podDefaultNamespace,
	)
	got, err := i.Pods(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := podListDefaultNamespace
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsEventsOccuredInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(event1, event2)
	got, err := i.Events(context.Background(), "nginx-ingress")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.EventList{
		Items: []corev1.Event{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Event",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx-config.17cca449e4d825fa",
					Namespace: "nginx-ingress",
				},
				Reason:  "Updated",
				Message: "Configuration from nginx-ingress/nginx-config was updated ",
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNoEventsOccuredInNotExistingNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(event1, event2)
	got, err := i.Events(context.Background(), "notExistingNamespace")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.EventList{
		Items: nil,
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNoEventsOccuredInExistingNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(defaultNameSpace)
	got, err := i.Events(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.EventList{
		Items: nil,
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsConfigMapsInGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(nginxIngressNameSpace, configMapNginxIngress)
	got, err := i.ConfigMaps(context.Background(), "nginx-ingress")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.ConfigMapList{
		Items: []corev1.ConfigMap{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx-config",
					Namespace: "nginx-ingress",
				},
				Data: map[string]string{"testdata": "hello inspector!"},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsEmptyConfigMapListOnNamespaceWithoutConfigMaps(t *testing.T) {
	t.Parallel()

	i := newTestInspector(nginxIngressNameSpace)
	got, err := i.ConfigMaps(context.Background(), "nginx-ingress")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.ConfigMapList{}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsExistingServicesInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		defaultNameSpace,
		teaServiceDefaultNS,
		coffeeServiceDefaultNS,
	)
	got, err := i.Services(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.ServiceList{
		Items: []corev1.Service{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-coffee",
					Namespace: "default",
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service-tea",
					Namespace: "default",
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNotExistingServicesInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(defaultNameSpace)
	got, err := i.Services(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.ServiceList{
		Items: nil,
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNotExistingDeploymentsInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(defaultNameSpace)
	got, err := i.Deployments(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &appsv1.DeploymentList{}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsDeploymentsInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		defaultNameSpace,
		fooBarNameSpace,
		deploymentDefaultNS,
		deploymentFooBarNS,
	)
	got, err := i.Deployments(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-deployment",
					Namespace: "default",
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	got, err = i.Deployments(context.Background(), "foobar")
	if err != nil {
		t.Fatal(err)
	}
	want = &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foobar-deployment",
					Namespace: "foobar",
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNotExistingStatefulSetsInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(defaultNameSpace)
	got, err := i.StatefulSets(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &appsv1.StatefulSetList{}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsStatefulSetsInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(defaultNameSpace, statefulSet)
	got, err := i.StatefulSets(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &appsv1.StatefulSetList{
		Items: []appsv1.StatefulSet{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "StatefulSet",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "web",
					Namespace: "default",
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsReplicaSetsInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(defaultNameSpace, replicaSet)
	got, err := i.ReplicaSets(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &appsv1.ReplicaSetList{
		Items: []appsv1.ReplicaSet{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ReplicaSet",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "replica",
					Namespace: "default",
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsLeasesInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(nginxIngressNameSpace, lease)
	got, err := i.Leases(context.Background(), "nginx-ingress")
	if err != nil {
		t.Fatal(err)
	}
	want := &coordv1.LeaseList{
		Items: []coordv1.Lease{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Lease",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx-ingress-leader-election",
					Namespace: "nginx-ingress",
				},
				Spec: coordv1.LeaseSpec{},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsNotExistingLeasesInAGivenNamespace(t *testing.T) {
	t.Parallel()

	i := newTestInspector(nginxIngressNameSpace)
	got, err := i.Leases(context.Background(), "nginx-ingress")
	if err != nil {
		t.Fatal(err)
	}
	want := &coordv1.LeaseList{}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListIngressClasses(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		defaultNameSpace,
		nginxIngressNameSpace,
		ingressClass,
	)
	got, err := i.IngressClasses(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := &netv1.IngressClassList{
		Items: []netv1.IngressClass{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "IngressClass",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "nginx",
				},
				Spec: netv1.IngressClassSpec{
					Controller: "nginx.org/ingress-controller",
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListIngresses(t *testing.T) {
	t.Parallel()

	i := newTestInspector(
		defaultNameSpace,
		nginxIngressNameSpace,
		ingress,
	)
	got, err := i.Ingresses(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	ingressClassName := "nginx"
	want := &netv1.IngressList{
		Items: []netv1.Ingress{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hello-ingress",
					Namespace: "default",
				},
				Spec: netv1.IngressSpec{
					IngressClassName: &ingressClassName,
				},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestInspectorListsCustomResourceDefinitions(t *testing.T) {
	t.Parallel()

}

func TestInspectorCollectsMetricsFromNodes(t *testing.T) {
	t.Parallel()
}

func TestInspectorCollectsHelmInformation(t *testing.T) {
	t.Parallel()
}

func TestInspectorCollectsHelmDeployments(t *testing.T) {
	t.Parallel()
}

// newTestInspector returns Inspector configured to use
// underlying K8s client and fake K8s runtime objects.
func newTestInspector(k8sobjects ...k8sruntime.Object) *inspector.Inspector {
	return &inspector.Inspector{K8sClient: newTestClientset(k8sobjects...)}
}

// newTestClientset takes K8s runtime objects and returns a k8s fake clientset.
func newTestClientset(objects ...k8sruntime.Object) *testClient.Clientset {
	client := testClient.NewSimpleClientset(objects...)
	client.Discovery().(*fakediscovery.FakeDiscovery).FakedServerVersion = &version.Info{
		GitVersion: "v1.29.2",
	}
	return client
}

// K8s cluster namespace.
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

	defaultNameSpace = &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			UID:  "121766aa-5d78-4c9e-8736-7faad1f2e345",
		},
		Spec: corev1.NamespaceSpec{},
	}

	nginxIngressNameSpace = &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-ingress",
			Namespace: "nginx-ingress",
			UID:       "441766aa-5d78-4c9e-8736-7faad1f2e987",
		},
		Spec: corev1.NamespaceSpec{},
	}

	fooBarNameSpace = &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foobar",
			Namespace: "foobar",
			UID:       "541766aa-5d78-4c9e-8736-7faad1f2e864",
		},
		Spec: corev1.NamespaceSpec{},
	}
)

// K8s cluster nodes.
var (
	clusterNode1 = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-node-1",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{},
	}

	clusterNode2 = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-node-2",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{},
	}

	clusterNode3 = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-node-3",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{},
	}
)

// Cloud providers' nodes for testing ProviderID lookups.
var (
	nodeAWS = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "aws:///eu-central-1a/i-088b4f07708408cc0",
		},
	}

	nodeAzure = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "azure:///subscriptions/ba96ef31-4a42-40f5-8740-03f7e3c439eb/resourceGroups/mc_hibrid-weu_be3rr5ovr8ulf_westeurope/providers/Microsoft.Compute/virtualMachines/aks-pool1-27255451-0",
		},
	}

	nodeGCP = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "gce://gcp-banzaidevgcp-nprd-38306/europe-north1-a/gke-vzf3z1vvleco9-pool1-7e48d363-8qz1",
		},
	}

	nodeKind = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "kind://docker/local/local-control-plane",
		},
	}

	nodeVSphere = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "vsphere://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeK3S = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "k3s://ip-1.2.3.4",
		},
	}

	nodeIBMCloud = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "ibmcloud://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeIBMPowerVS = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "ibmpowervs://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeCloudStack = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "cloudstack://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeOpenStack = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "openstack://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeDigitalOcean = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "digitalocean://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeEquinixMetal = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "equinixmetal://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeAlibaba = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "alicloud://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}
)

// Nodes with missing or malformed ProviderID field.
var (
	nodeMalformedPlatformID = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "//4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeMalformedPartialPlatformID = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "://4232e3c7-d83c-d72b-758c-71d07a3d9310",
		},
	}

	nodeMalformedEmptyPlatformID = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "",
		},
	}

	nodeMalformedBlankPlatformID = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: " ",
		},
	}
)

var replicaSetID = "239766ff-5a78-4a1e-8736-7faad1f2e122"

// Pods for testing.
var (
	pod1 = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-ingress",
			Namespace: "nginx-ingress",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "ReplicaSet",
					Name: "nginx-ingress",
					UID:  types.UID(replicaSetID),
				},
			},
			Labels: map[string]string{
				"app":                    "nginx-ingress",
				"app.kubernetes.io/name": "nginx-ingress",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "nginx-ingress",
					Image:           "nginx-ingress",
					ImagePullPolicy: "Always",
					Env: []corev1.EnvVar{
						{
							Name:  "POD_NAMESPACE",
							Value: "nginx-ingress",
						},
						{
							Name:  "POD_NAME",
							Value: "nginx-ingress",
						},
					},
				},
			},
		},
	}

	podDefaultNamespace = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "inspector",
			Namespace:       "default",
			OwnerReferences: []metav1.OwnerReference{},
			Labels: map[string]string{
				"app":                    "inspector",
				"app.kubernetes.io/name": "inspector",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "inspector",
					Image:           "inspector",
					ImagePullPolicy: "Always",
					Env: []corev1.EnvVar{
						{
							Name:  "POD_NAMESPACE",
							Value: "default",
						},
						{
							Name:  "POD_NAME",
							Value: "inspector",
						},
					},
				},
			},
		},
	}
)

// List of pods in default namespace.
var podListDefaultNamespace = &corev1.PodList{
	Items: []corev1.Pod{*podDefaultNamespace},
}

// K8s Events used for testing.
var (
	event1 = &corev1.Event{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Event",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test event",
			Namespace: "default",
		},
		Reason:  "Updated",
		Message: "human readable message",
	}

	event2 = &corev1.Event{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Event",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-config.17cca449e4d825fa",
			Namespace: "nginx-ingress",
		},
		Reason:  "Updated",
		Message: "Configuration from nginx-ingress/nginx-config was updated ",
	}
)

// Config Maps used for testing.
var (
	configMapNginxIngress = &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-config",
			Namespace: "nginx-ingress",
		},
		Data: map[string]string{"testdata": "hello inspector!"},
	}
)

// Services used for testing.
var (
	coffeeServiceDefaultNS = &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-coffee",
			Namespace: "default",
		},
	}

	teaServiceDefaultNS = &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-tea",
			Namespace: "default",
		},
	}
)

// Deployments used for testing.
var (
	deploymentDefaultNS = &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-deployment",
			Namespace: "default",
		},
		Spec:   appsv1.DeploymentSpec{},
		Status: appsv1.DeploymentStatus{},
	}

	deploymentFooBarNS = &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foobar-deployment",
			Namespace: "foobar",
		},
		Spec:   appsv1.DeploymentSpec{},
		Status: appsv1.DeploymentStatus{},
	}
)

// StatefulSet used for testing.
var statefulSet = &appsv1.StatefulSet{
	TypeMeta: metav1.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "web",
		Namespace: "default",
	},
	Spec: appsv1.StatefulSetSpec{},
}

// ReplicaSet used for testing.
var replicaSet = &appsv1.ReplicaSet{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ReplicaSet",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "replica",
		Namespace: "default",
	},
	Spec: appsv1.ReplicaSetSpec{},
}

// Lease used for testing.
var lease = &coordv1.Lease{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Lease",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "nginx-ingress-leader-election",
		Namespace: "nginx-ingress",
	},
	Spec: coordv1.LeaseSpec{},
}

// IngressClass and Ingress for testing.
var (
	ingressClass = &netv1.IngressClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IngressClass",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: netv1.IngressClassSpec{
			Controller: "nginx.org/ingress-controller",
		},
	}

	ingressClassName = "nginx"
	ingress          = &netv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello-ingress",
			Namespace: "default",
		},
		Spec: netv1.IngressSpec{
			IngressClassName: &ingressClassName,
		},
	}
)
