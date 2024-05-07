package inspector

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	coordv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crd "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Inspector is an inspector client.
type Inspector struct {
	Verbose       bool
	K8sClient     kubernetes.Interface
	CRDClient     *crd.Clientset
	MetricsClient *metrics.Clientset
}

// BuildInspectorFromKubeConfig builds an inspector client ready to interact with the K8s cluster.
func BuildInspectorFromKubeConfig() (*Inspector, error) {
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

	i := Inspector{
		Verbose:       false,
		K8sClient:     kubeClient,
		CRDClient:     crdClient,
		MetricsClient: metricsClient,
	}
	return &i, nil
}

// ClusterVersion returns K8s version.
func (i *Inspector) ClusterVersion() (string, error) {
	sv, err := i.K8sClient.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return sv.String(), nil
}

// ClusterID returns kube-system namespace UID representing K8s clusterID.
func (i *Inspector) ClusterID(ctx context.Context) (string, error) {
	cluster, err := i.K8sClient.CoreV1().Namespaces().Get(ctx, "kube-system", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(cluster.UID), nil
}

// Platform returns K8s platform name.
func (i *Inspector) Platform(ctx context.Context) (string, error) {
	nodes, err := i.K8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
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

// Nodes returns the total number of [nodes] in a cluster.
//
// [nodes]: https://kubernetes.io/docs/concepts/architecture/nodes/
func (i *Inspector) Nodes(ctx context.Context) (int, error) {
	nodes, err := i.K8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(nodes.Items), nil
}

// Pods returns list of [pods] in a given [namespace].
//
// [pods]: https://kubernetes.io/docs/concepts/workloads/pods/
// [namespace]: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
func (i *Inspector) Pods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	pods, err := i.K8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// PodLog represents a pod and collected logs
// from containers in the pod.
type PodLog struct {
	Name string `json:"name"`
	Log  string `json:"log"`
}

// Podlogs returns logs from [pods] in a given [namespace].
//
// [pods]: https://kubernetes.io/docs/concepts/workloads/pods/
// [namespace]: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
func (i *Inspector) Podlogs(ctx context.Context, namespace string) ([]PodLog, error) {
	pods, err := i.K8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	logs := []PodLog{}
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			logReq := i.K8sClient.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
			res, err := logReq.Stream(ctx)
			if err != nil {
				return nil, err
			}
			log, err := io.ReadAll(res)
			if err != nil {
				return []PodLog{}, err
			}

			pl := PodLog{
				Name: fmt.Sprintf("%s_%s", pod.Name, container.Name),
				Log:  string(log),
			}
			logs = append(logs, pl)
		}
	}
	return logs, nil
}

// Events returns [events] for a given namespace.
//
// [events]: https://kubernetes.io/docs/reference/kubectl/generated/kubectl_events/
func (i *Inspector) Events(ctx context.Context, namespace string) (*corev1.EventList, error) {
	events, err := i.K8sClient.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ConfigMaps returns a list of [config maps] for a given namespace.
//
// [config maps]: https://kubernetes.io/docs/concepts/configuration/configmap/
func (i *Inspector) ConfigMaps(ctx context.Context, namespace string) (*corev1.ConfigMapList, error) {
	configMaps, err := i.K8sClient.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// json.MarshalIndent(cm, "", "  ")
	return configMaps, nil
}

// Services returns a list of [services] for a given namespace.
//
// [services]: https://kubernetes.io/docs/concepts/services-networking/service/
func (i *Inspector) Services(ctx context.Context, namespace string) (*corev1.ServiceList, error) {
	services, err := i.K8sClient.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return services, nil
}

// Deployments returns a list of [deployments] in a given namespace.
//
// [deployments]: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
func (i *Inspector) Deployments(ctx context.Context, namespace string) (*appsv1.DeploymentList, error) {
	deployments, err := i.K8sClient.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

// StatefulSets returns a list of [stateful sets] in a given namespace.
//
// [stateful set]: https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/
func (i *Inspector) StatefulSets(ctx context.Context, namespace string) (*appsv1.StatefulSetList, error) {
	statefulSets, err := i.K8sClient.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return statefulSets, nil
}

// ReplicaSets returns a list of [replica sets] in a given namespace.
//
// [replica sets]: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/
func (i *Inspector) ReplicaSets(ctx context.Context, namespace string) (*appsv1.ReplicaSetList, error) {
	replicaSets, err := i.K8sClient.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return replicaSets, nil
}

// Leasess returns a list of [leases] in a given namespace.
//
// [leases]: https://kubernetes.io/docs/concepts/architecture/leases/
func (i *Inspector) Leases(ctx context.Context, namespace string) (*coordv1.LeaseList, error) {
	leases, err := i.K8sClient.CoordinationV1().Leases(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return leases, nil
}

// IngressClasses returns a list of [ingress classes] in a cluster.
//
// [ingress classes]: https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-class
func (i *Inspector) IngressClasses(ctx context.Context) (*netv1.IngressClassList, error) {
	ingressClasses, err := i.K8sClient.NetworkingV1().IngressClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingressClasses, nil
}

// Ingresses returns a list of [ingresses] in a given namespace.
//
// [ingresses]: https://kubernetes.io/docs/concepts/services-networking/ingress/
func (i *Inspector) Ingresses(ctx context.Context, namespace string) (*netv1.IngressList, error) {
	ingresses, err := i.K8sClient.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingresses, nil
}

// CustomResourceDefinitions returns a list of [CRDs] in a cluster.
//
// [CRDs]: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
func (i *Inspector) CustomResourceDefinitions(ctx context.Context) (*apiextv1.CustomResourceDefinitionList, error) {
	crds, err := i.CRDClient.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return crds, nil
}

// ClusterNodes returns a list of [nodes] in a [cluster].
//
// [nodes]: https://kubernetes.io/docs/concepts/architecture/nodes/
// [cluster]: https://kubernetes.io/docs/concepts/cluster-administration/
func (i *Inspector) ClusterNodes(ctx context.Context) (*corev1.NodeList, error) {
	nodes, err := i.K8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// NodeMetrics returns a list of [node metrics] in a cluster.
//
// [node metrics]: https://kubernetes.io/docs/concepts/cluster-administration/system-metrics/
func (i *Inspector) NodeMetrics(ctx context.Context) (*v1beta1.NodeMetricsList, error) {
	metrics, err := i.MetricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// PodMetrics returns a list of [pods metrics] in a given namespace.
//
// [pods metrics]: https://kubernetes.io/docs/concepts/cluster-administration/kube-state-metrics/
func (i *Inspector) PodMetrics(ctx context.Context, namespace string) (*v1beta1.PodMetricsList, error) {
	metrics, err := i.MetricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// RunDiagnostics collects cluster data points for a given namespace.
func (i *Inspector) Report(ctx context.Context, namespace string) (Report, error) {
	version, err := i.ClusterVersion()
	if err != nil {
		return Report{}, err
	}
	id, err := i.ClusterID(ctx)
	if err != nil {
		return Report{}, err
	}
	n, err := i.Nodes(ctx)
	if err != nil {
		return Report{}, err
	}
	p, err := i.Platform(ctx)
	if err != nil {
		return Report{}, err
	}

	pods, err := i.Pods(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	podLogs, err := i.Podlogs(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	events, err := i.Events(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	configMaps, err := i.ConfigMaps(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	services, err := i.Services(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	deployments, err := i.Deployments(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	statefulSets, err := i.StatefulSets(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	replicaSets, err := i.ReplicaSets(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	leases, err := i.Leases(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	ingressClasses, err := i.IngressClasses(ctx)
	if err != nil {
		return Report{}, err
	}

	ingresses, err := i.Ingresses(ctx, namespace)
	if err != nil {
		return Report{}, err
	}

	crds, err := i.CustomResourceDefinitions(ctx)
	if err != nil {
		return Report{}, err
	}

	clusterNodes, err := i.ClusterNodes(ctx)
	if err != nil {
		return Report{}, err
	}

	return Report{
		K8sVersion:     version,
		ClusterID:      id,
		Nodes:          n,
		Platform:       p,
		Pods:           pods,
		Podlogs:        podLogs,
		Events:         events,
		ConfigMaps:     configMaps,
		Services:       services,
		Deployments:    deployments,
		StatefulSets:   statefulSets,
		ReplicaSets:    replicaSets,
		Leases:         leases,
		IngressClasses: ingressClasses,
		Ingresses:      ingresses,
		CRDs:           crds,
		ClusterNodes:   clusterNodes,
	}, nil
}

// ReportJSON returns collected metrics in a JSON format.
func ReportJSON(rep Report) (string, error) {
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Report holds collected data points.
type Report struct {
	K8sVersion     string                                 `json:"k8s_version"`
	ClusterID      string                                 `json:"cluster_id"`
	Nodes          int                                    `json:"nodes"`
	Platform       string                                 `json:"platform"`
	Pods           *corev1.PodList                        `json:"pods"`
	Podlogs        []PodLog                               `json:"pod_logs"`
	Events         *corev1.EventList                      `json:"events"`
	ConfigMaps     *corev1.ConfigMapList                  `json:"config_maps"`
	Services       *corev1.ServiceList                    `json:"services"`
	Deployments    *appsv1.DeploymentList                 `json:"deployments"`
	StatefulSets   *appsv1.StatefulSetList                `json:"stateful_sets"`
	ReplicaSets    *appsv1.ReplicaSetList                 `json:"replica_sets"`
	Leases         *coordv1.LeaseList                     `json:"leases"`
	IngressClasses *netv1.IngressClassList                `json:"ingress_classes"`
	Ingresses      *netv1.IngressList                     `json:"ingresses"`
	CRDs           *apiextv1.CustomResourceDefinitionList `json:"crds"`
	ClusterNodes   *corev1.NodeList                       `json:"cluster_nodes"`
}

var usage = `Usage:

	inspector [-h] [-v] [-n] namespace

Collect K8s and Ingress Controller diagnostics in the given namespace.

In verbose mode (-v), prints out progess, steps and all data points to stdout.`

// Main runs the inspector program.
func Main() int {
	namespace := flag.String("n", "default", "K8s namespace")
	verbose := flag.Bool("v", false, "verbose output")
	help := flag.Bool("h", false, "show help")
	flag.Parse()

	if *help {
		fmt.Println(usage)
		return 0
	}

	i, err := BuildInspectorFromKubeConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	i.Verbose = *verbose

	report, err := i.Report(context.Background(), *namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	rep, err := ReportJSON(report)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		return 1
	}
	fmt.Println(rep)
	return 0
}
