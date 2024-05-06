package inspector

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	coordv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crd "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Report holds collected data points.
type Report struct {
	K8sVersion string
	ClusterID  string
	Nodes      int
	Platform   string
	Pods       string
	Podlogs    string
}

// String implements stringer interface.
func (r Report) String() string {
	const report = "=== Cluster Info ===\nVersion: %s\nClusterID: %s\nNodes: %d\nPlatform: %s\n=== Pods ===\n%s\n=== Pod logs ===\n%s\n"
	return fmt.Sprintf(report,
		r.K8sVersion,
		r.ClusterID,
		r.Nodes,
		r.Platform,
		r.Pods,
		r.Podlogs,
	)
}

// Client represents Inspector client.
type Client struct {
	Verbose       bool
	K8sClient     kubernetes.Interface
	CRDClient     *crd.Clientset
	MetricsClient *metrics.Clientset
}

// BuildClientFromKubeConfig builds inspector client ready to interact with the cluster.
func BuildClientFromKubeConfig() (*Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", home+"/.kube/config")
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	crdClient, err := crd.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	c := Client{
		Verbose:       false,
		K8sClient:     kubeClient,
		CRDClient:     crdClient,
		MetricsClient: metricsClient,
	}
	return &c, nil
}

// ClusterVersion returns K8s version.
func (c *Client) ClusterVersion() (string, error) {
	sv, err := c.K8sClient.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return sv.String(), nil
}

// ClusterID returns kube-system namespace UID representing K8s clusterID.
func (c *Client) ClusterID(ctx context.Context) (string, error) {
	cluster, err := c.K8sClient.CoreV1().Namespaces().Get(ctx, "kube-system", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(cluster.UID), nil
}

// Platform returns K8s platform name.
func (c *Client) Platform(ctx context.Context) (string, error) {
	nodes, err := c.K8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(nodes.Items) == 0 {
		return "", errors.New("cannot verify platform name")
	}
	return platformName(nodes.Items[0].Spec.ProviderID), nil
}

// platformName takes a string representing a K8s PlatformID
// retrieved from a cluster node and returns a string
// representing the platform name.
func platformName(providerID string) string {
	provider := strings.TrimSpace(providerID)
	if provider == "" {
		return "unknown"
	}
	provider = strings.ToLower(providerID)
	p := strings.Split(provider, ":")
	if len(p) == 0 {
		return "unknown"
	}
	if p[0] == "" {
		return "unknown"
	}
	return p[0]
}

// Nodes returns the total number of nodes in the cluster.
func (c *Client) Nodes(ctx context.Context) (int, error) {
	nodes, err := c.K8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(nodes.Items), nil
}

// Pods returns info about all pods in a given namespace.
func (c *Client) Pods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	pods, err := c.K8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// Podlogs returns logs from pods in the given namespace.
func (c *Client) Podlogs(ctx context.Context, namespace string) (map[string]string, error) {
	pods, err := c.K8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	podLogs := make(map[string]string)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			logReq := c.K8sClient.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
			res, err := logReq.Stream(ctx)
			if err != nil {
				return nil, err
			}
			buf := &bytes.Buffer{}
			_, err = io.Copy(buf, res)
			if err != nil {
				return nil, err
			}
			podLogs[pod.Name+"_"+container.Name] = buf.String()
		}
	}
	return podLogs, nil
}

// Events returns events from given namespace.
func (c *Client) Events(ctx context.Context, namespace string) (*corev1.EventList, error) {
	events, err := c.K8sClient.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *Client) ConfigMaps(ctx context.Context, namespace string) (*corev1.ConfigMapList, error) {
	cm, err := c.K8sClient.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// json.MarshalIndent(cm, "", "  ")
	return cm, nil
}

func (c *Client) Services(ctx context.Context, namespace string) (*corev1.ServiceList, error) {
	sl, err := c.K8sClient.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return sl, nil
}

func (c *Client) Deployments(ctx context.Context, namespace string) (*appsv1.DeploymentList, error) {
	dl, err := c.K8sClient.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return dl, nil
}

func (c *Client) StatefulSets(ctx context.Context, namespace string) (*appsv1.StatefulSetList, error) {
	ss, err := c.K8sClient.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (c *Client) ReplicaSets(ctx context.Context, namespace string) (*appsv1.ReplicaSetList, error) {
	rs, err := c.K8sClient.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (c *Client) Leases(ctx context.Context, namespace string) (*coordv1.LeaseList, error) {
	leases, err := c.K8sClient.CoordinationV1().Leases(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return leases, nil
}

func (c *Client) CustomResourceDefinitions(ctx context.Context) (*apiextv1.CustomResourceDefinitionList, error) {
	crds, err := c.CRDClient.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return crds, nil
}

func (c *Client) ClusterNodes(ctx context.Context) (*corev1.NodeList, error) {
	nodes, err := c.K8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c *Client) NodeMetrics(ctx context.Context) (*v1beta1.NodeMetricsList, error) {
	metrics, err := c.MetricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (c *Client) PodMetrics(ctx context.Context, namespace string) (*v1beta1.PodMetricsList, error) {
	metrics, err := c.MetricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// RunDiagnostics collects cluster data points for given namespace.
func (c *Client) Report(ctx context.Context, namespace string) (Report, error) {
	version, err := c.ClusterVersion()
	if err != nil {
		return Report{}, err
	}
	id, err := c.ClusterID(ctx)
	if err != nil {
		return Report{}, err
	}
	n, err := c.Nodes(ctx)
	if err != nil {
		return Report{}, err
	}
	p, err := c.Platform(ctx)
	if err != nil {
		return Report{}, err
	}
	pods, err := c.Pods(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	pl, err := c.Podlogs(ctx, namespace)
	if err != nil {
		return Report{}, err
	}
	var podl strings.Builder
	for k, v := range pl {
		podl.WriteString(fmt.Sprintf("%s\n%s\n", k, v))
	}
	return Report{
		K8sVersion: version,
		ClusterID:  id,
		Nodes:      n,
		Platform:   p,
		Pods:       pods.String(),
		Podlogs:    podl.String(),
	}, nil
}

var usage = `Usage:

	inspector [-v] [-n] namespace

Collect K8s and Ingress Controller diagnostics in the given namespace.

In verbose mode (-v), prints out progess, steps and all data points to stdout.`

func Main() int {
	namespace := flag.String("n", "default", "K8s namespace")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()

	c, err := BuildClientFromKubeConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	c.Verbose = *verbose

	report, err := c.Report(context.Background(), *namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	fmt.Println(report)
	return 0
}
