package inspector

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Report struct {
	K8sVersion string
	ClusterID  string
	Nodes      int
	Platform   string
}

func (r Report) String() string {
	return fmt.Sprintf(
		"Version: %s\nClusterID: %s\nNodes: %d\nPlatform: %s\n",
		r.K8sVersion,
		r.ClusterID,
		r.Nodes,
		r.Platform,
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
	p, err := c.Platform(ctx)
	if err != nil {
		fmt.Fprint(c.Output, err.Error())
	}

	report := Report{
		K8sVersion: version,
		ClusterID:  id,
		Nodes:      n,
		Platform:   p,
	}
	c.diagnostics = report
}

// configDir returns path to the K8s configuration.
//
// If user exported the env var XDG_DATA_HOME inspector
// will use this location to look for k8s config file.
// If XDG_DATA_HOME is not set inspector looks for k8s config
// in K8s default directory: $HOME/.kube/.
func configDir() string {
	path := os.Getenv("XDG_DATA_HOME")
	if path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

func k8sClient(path string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// NewClientFromConfig takes a path to the k8s config path
// and returns inspector client ready to interact with the cluster.
func NewClientFromConfig(path string) (*Client, error) {
	k8sClient, err := k8sClient(path)
	if err != nil {
		return nil, err
	}
	c := Client{
		Verbose:   false,
		Output:    os.Stdout,
		K8sClient: k8sClient,
	}
	return &c, nil
}

var usage = `Usage: inspector [-v] namespace

Gather NIC diagnostics for the given namespace

In verbose mode (-v), prints out all data points to stdout.`

func Main() int {
	kubeconfig := flag.String("kubeconfig", filepath.Join(configDir(), ".kube", "config"), "path to the kubeconfig file")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	namespace := flag.Args()[0]

	c, err := NewClientFromConfig(*kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	c.Verbose = *verbose

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		c.RunDiagnostic(ctx, namespace)
		cancel()
	}()
	<-ctx.Done()

	report := c.Report()
	fmt.Fprintf(os.Stdout, "%s\n", report)
	return 0
}
