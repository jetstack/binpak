package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"net/http"
	"os"
)

// Info is all the information displayed by the UI.
type Info struct {
	ClusterName string      `json:"clusterName"`
	Instances   []*Instance `json:"instances"`
}

// Instance represents a cluster node.
type Instance struct {
	Name      string      `json:"name"`
	Capacity  *Resources  `json:"capacity"`
	Usage     *Resources  `json:"usage"`
	Workloads []*Workload `json:"workloads"`
}

// Resources represents CPU and memory in standardised units.
type Resources struct {
	MemoryMi int64 `json:"memoryMi"`
	CpuM     int64 `json:"cpuM"`
}

// Workload represents an instance of a running workload.
type Workload struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Requests *Resources `json:"requests"`
	Limits   *Resources `json:"limits"`
	Usage    *Resources `json:"usage"`
}

func main() {
	// Default the Kube context and read as an argument if set.
	contextName := ""
	if len(os.Args) >= 2 {
		contextName = os.Args[1]
	}

	_, usageMetrics := os.LookupEnv("BINPAK_ENABLE_USAGE_METRICS")

	// Create Kube client to access the cluster.
	client, metricsClient, err := newClient(contextName)
	if err != nil {
		log.Fatal(err)
	}

	// Can use this for debugging to just print the info
	//i, err := updateInfo(client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//b, err := json.MarshalIndent(i, "", "    ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(b))
	//os.Exit(0)

	router := gin.Default()
	router.GET("/info", func(c *gin.Context) {
		i, err := updateInfo(client, metricsClient, usageMetrics)
		if err != nil {
			log.Fatal(err)
		}
		c.IndentedJSON(http.StatusOK, i)
	})
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.Run("0.0.0.0:8080")
}

func updateInfo(client kubernetes.Interface, metricsClient *metricsv.Clientset, enableUsageMetrics bool) (*Info, error) {
	// Get a list of nodes from the cluster.
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Create a list of instances to represent cluster nodes and the workloads running on them.
	instanceList := make([]*Instance, 0)
	for _, node := range nodes.Items {
		nodeName := node.Name
		nodeMemory := node.Status.Capacity["memory"]
		nodeMemoryMi := (nodeMemory.Value() / 1024) / 1024
		nodeCpu := node.Status.Capacity["cpu"]
		nodeCpuM := nodeCpu.MilliValue()

		nodeMetrices, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(context.TODO(), nodeName, meta_v1.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}

		inst := Instance{
			Name: nodeName,
			Capacity: &Resources{
				MemoryMi: nodeMemoryMi,
				CpuM:     nodeCpuM,
			},
			Usage: &Resources{
				MemoryMi: nodeMetrices.Usage.Memory().Value(),
				CpuM:     (nodeMetrices.Usage.Cpu().Value() / 1024) / 1024,
			},
			Workloads: make([]*Workload, 0),
		}

		instanceList = append(instanceList, &inst)
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		podName := pod.Name
		podRequestMemoryMi := int64(0)
		podRequestCpuM := int64(0)
		podLimitMemoryMi := int64(0)
		podLimitCpuM := int64(0)
		for _, container := range pod.Spec.Containers {
			podContainerRequestMemory := container.Resources.Requests["memory"]
			podContainerRequestMemoryMi := (podContainerRequestMemory.Value() / 1024) / 1024
			podRequestMemoryMi = podRequestMemoryMi + podContainerRequestMemoryMi

			podContainerRequestCpu := container.Resources.Requests["cpu"]
			podContainerRequestCpuM := podContainerRequestCpu.MilliValue()
			podRequestCpuM = podRequestCpuM + podContainerRequestCpuM

			podContainerLimitsMemory := container.Resources.Limits["memory"]
			podContainerLimitsMemoryMi := (podContainerLimitsMemory.Value() / 1024) / 1024
			podLimitMemoryMi = podLimitMemoryMi + podContainerLimitsMemoryMi

			podContainerLimitsCpu := container.Resources.Limits["cpu"]
			podContainerLimitsCpuM := podContainerLimitsCpu.MilliValue()
			podLimitCpuM = podLimitCpuM + podContainerLimitsCpuM
		}

		podUsageMemoryMi := int64(0)
		podUsageCpuM := int64(0)
		if enableUsageMetrics == true {
			podUsageMemoryMi, podUsageCpuM, err = getPodUsageMetrics(metricsClient, pod, podName)
			if err != nil {
				return nil, err
			}
		}
		w := &Workload{
			Name: podName,
			Requests: &Resources{
				MemoryMi: podRequestMemoryMi,
				CpuM:     podRequestCpuM,
			},
			Limits: &Resources{
				MemoryMi: podLimitMemoryMi,
				CpuM:     podLimitCpuM,
			},
			Usage: &Resources{
				MemoryMi: podUsageMemoryMi,
				CpuM:     podUsageCpuM,
			},
		}

		for _, inst := range instanceList {
			if pod.Spec.NodeName == inst.Name {
				inst.Workloads = append(inst.Workloads, w)
				break
			}
		}
	}

	i := &Info{
		ClusterName: "test",
		Instances:   instanceList,
	}
	return i, nil
}

func getPodUsageMetrics(metricsClient *metricsv.Clientset, pod v1.Pod, podName string) (int64, int64, error) {
	podUsageMemoryMi := int64(0)
	podUsageCpuM := int64(0)
	podMetrices, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(context.TODO(), podName, meta_v1.GetOptions{})
	if err != nil {
		return 0, 0, err
	}
	for _, containerMetrices := range podMetrices.Containers {
		podUsageMemoryMi = podUsageMemoryMi + (containerMetrices.Usage.Memory().Value()/1024)/1024
		podUsageCpuM = podUsageCpuM + containerMetrices.Usage.Cpu().Value()
	}
	return podUsageMemoryMi, podUsageCpuM, nil
}

func newClient(contextName string) (kubernetes.Interface, *metricsv.Clientset, error) {
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return client, metricsClient, nil
}
