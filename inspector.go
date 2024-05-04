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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Report struct {
	K8sVersion string
	ClusterID  string
	Nodes      int
	Platform   string
	Pods       string
	Podlogs    string
}

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
	Verbose   bool
	K8sClient kubernetes.Interface
}

// BuildClientFromKubeConfig inspector client ready to interact with the cluster.
func BuildClientFromKubeConfig() (*Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", home+"/.kube/config")
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	c := Client{
		Verbose:   false,
		K8sClient: clientset,
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

var usage = `Usage: inspector [-v] namespace

Gather K8s and NIC diagnostics in the given namespace

In verbose mode (-v), prints out all data points to stdout.`

func Main() int {
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	namespace := flag.Args()[0]

	c, err := BuildClientFromKubeConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	c.Verbose = *verbose

	report, err := c.Report(context.Background(), namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	fmt.Println(report)
	return 0
}
