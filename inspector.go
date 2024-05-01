package inspector

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Client represents Inspector client.
type Client struct {
	K8sClient kubernetes.Interface
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
