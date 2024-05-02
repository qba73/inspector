package inspector_test

import (
	"context"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			clusterNode1,
			clusterNode2,
			clusterNode3,
		),
	}
	c.Output = io.Discard
	c.RunDiagnostic(context.Background(), "default")

	got := c.Report()
	want := inspector.Report{
		K8sVersion: "v1.29.2",
		ClusterID:  "421766aa-5d78-4c9e-8736-7faad1f2e927",
		Nodes:      3,
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

// Nodes with missing or malformed PorviderID.
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
