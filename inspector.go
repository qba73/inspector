package inspector

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Report struct {
	K8sVersion string
	ClusterID  string
	Nodes      int
}

func (r Report) String() string {
	return fmt.Sprintf(
		"Version: %s\nClusterID: %s\nNodes: %d\n",
		r.K8sVersion,
		r.ClusterID,
		r.Nodes,
	)
}

// Client represents Inspector client.
type Client struct {
	Verbose     bool
	Output      io.Writer
	K8sClient   kubernetes.Interface
	diagnostics Report
}

func NewClient() *Client {
	// todo create default k8s clientset from config
	return &Client{
		Verbose: false,
		Output:  os.Stdout,
	}
}

func (c *Client) Report() Report {
	return c.diagnostics
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

// RunDiagnostics collects cluster data points for given namespace.
func (c *Client) RunDiagnostic(ctx context.Context, namespace string) {
	version, err := c.ClusterVersion()
	if err != nil {
		fmt.Fprint(c.Output, err.Error())
	}
	id, err := c.ClusterID(ctx)
	if err != nil {
		fmt.Fprint(c.Output, err.Error())
	}
	n, err := c.Nodes(ctx)
	if err != nil {
		fmt.Fprint(c.Output, err.Error())
	}

	report := Report{
		K8sVersion: version,
		ClusterID:  id,
		Nodes:      n,
	}
	c.diagnostics = report
}

var usage = `Usage: inspector [-v] namespace

Gather NIC diagnostics for the given namespace

In verbose mode (-v), prints out all data points to stdout.`

func Main() int {
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	namespace := flag.Args()[0]
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	c := NewClient() // todo setup defualt k8s client
	c.Verbose = *verbose

	go func() {
		c.RunDiagnostic(ctx, namespace)
		cancel()
	}()
	<-ctx.Done()

	report := c.Report()
	fmt.Fprintf(os.Stdout, "%s\n", report)
	return 0
}
