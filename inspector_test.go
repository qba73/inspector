package inspector_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/inspector"
	"k8s.io/apimachinery/pkg/version"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

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

func TestInspectorCollectsNumberOfNodesInTheCluster(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			clusterNode1,
			clusterNode2,
			clusterNode3,
		),
	}
	got, err := c.Nodes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := 3
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestInspectorCollectsDiagnosticData(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeAWS,
			nodeAWS2,
			nodeAWS3,
		),
	}
	got, err := c.Report(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}

	want := inspector.Report{
		K8sVersion: "v1.29.2",
		ClusterID:  "421766aa-5d78-4c9e-8736-7faad1f2e927",
		Nodes:      3,
		Platform:   "aws",
		Pods:       podListEmpty.String(),
	}

	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestInspectorCollectsPlatformNameOnAWSNode(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeAWS,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeAzure,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeGCP,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeKind,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeVSphere,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeK3S,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeIBMCloud,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeIBMPowerVS,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeCloudStack,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeOpenStack,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeDigitalOcean,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeEquinixMetal,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeAlibaba,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeMalformedBlankPlatformID,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeMalformedPlatformID,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeMalformedBlankPlatformID,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeMalformedEmptyPlatformID,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nodeMalformedPartialPlatformID,
		),
	}
	got, err := c.Platform(context.Background())
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			nginxIngressNameSpace,
			pod1,
		),
	}
	got, err := c.Pods(context.Background(), "notExistingNamespace")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.PodList{}
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestInspectorListsNotExistingPodsInDefaultNamespace(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			defaultNameSpace,
		),
	}
	got, err := c.Pods(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := &corev1.PodList{}
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestInspectorListsExistingPodsInDefaultNamespace(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(
			kubeSystemNameSpace,
			defaultNameSpace,
			podDefaultNamespace,
		),
	}
	got, err := c.Pods(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	want := podListDefaultNamespace
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestInspectorListsEventsOccuredInAGivenNamespace(t *testing.T) {
	t.Parallel()

	c := inspector.Client{
		K8sClient: newTestClientset(
			event1,
			event2,
		),
	}
	got, err := c.Events(context.Background(), "nginx-ingress")
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			event1,
			event2,
		),
	}
	got, err := c.Events(context.Background(), "notExistingNamespace")
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			defaultNameSpace,
		),
	}
	got, err := c.Events(context.Background(), "default")
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			nginxIngressNameSpace,
			configMapNginxIngress,
		),
	}
	got, err := c.ConfigMaps(context.Background(), "nginx-ingress")
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			nginxIngressNameSpace,
		),
	}
	got, err := c.ConfigMaps(context.Background(), "nginx-ingress")
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			defaultNameSpace,
			teaServiceDefaultNS,
			coffeeServiceDefaultNS,
		),
	}
	got, err := c.Services(context.Background(), "default")
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

	c := inspector.Client{
		K8sClient: newTestClientset(
			defaultNameSpace,
		),
	}
	got, err := c.Services(context.Background(), "default")
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

func TestInspectorListsDeploymentsInAGivenNamespace(t *testing.T) {
	t.Parallel()

}

func TestInspectorListsStatefulSetsInAGivenNamespace(t *testing.T) {
	t.Parallel()

}

func TestInspectorListsReplicaSetsInAGivenNamespace(t *testing.T) {
	t.Parallel()

}

func TestInspectorListsLeasesInAGivenNamespace(t *testing.T) {
	t.Parallel()

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

	nodeAWS2 = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node-aws-2",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "aws:///eu-central-1a/i-088b4f07708408ca0",
		},
	}

	nodeAWS3 = &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node-aws-3",
			Namespace: "default",
		},
		Spec: corev1.NodeSpec{
			ProviderID: "aws:///eu-central-1a/i-088b4f07708408va0",
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

var (
	replicaSetID = "239766ff-5a78-4a1e-8736-7faad1f2e122"
	daemonSetID  = "319766ff-5c78-4a9a-8736-7faad1f2e234"
)

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

	pod2 = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-ingress-2",
			Namespace: "nginx-ingress",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "DaemonSet",
					Name: "nginx-ingress",
					UID:  types.UID(daemonSetID),
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
var (
	podListDefaultNamespace = &corev1.PodList{
		Items: []corev1.Pod{*podDefaultNamespace},
	}

	podListEmpty = &corev1.PodList{
		Items: []corev1.Pod{},
	}
)

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

	waterServiceWaterNS = &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-water",
			Namespace: "water",
		},
	}
)
