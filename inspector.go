package inspector

import (
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
