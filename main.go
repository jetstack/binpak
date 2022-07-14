package main

import (
	"context"
	"github.com/gin-gonic/gin"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
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
}

func main() {
	// Default the Kube context and read as an argument if set.
	contextName := ""
	if len(os.Args) >= 2 {
		contextName = os.Args[1]
	}

	// Create Kube client to access the cluster.
	client, err := newClient(contextName)
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
		i, err := updateInfo(client)
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

func updateInfo(client kubernetes.Interface) (*Info, error) {
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

		inst := Instance{
			Name: nodeName,
			Capacity: &Resources{
				MemoryMi: nodeMemoryMi,
				CpuM:     nodeCpuM,
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
			podContainerRequestCpuMi := podContainerRequestCpu.MilliValue()
			podRequestCpuM = podRequestCpuM + podContainerRequestCpuMi

			podContainerLimitsMemory := container.Resources.Requests["memory"]
			podContainerLimitsMemoryMi := (podContainerLimitsMemory.Value() / 1024) / 1024
			podLimitMemoryMi = podLimitMemoryMi + podContainerLimitsMemoryMi

			podContainerLimitsCpu := container.Resources.Requests["cpu"]
			podContainerLimitsCpuMi := podContainerLimitsCpu.MilliValue()
			podLimitCpuM = podLimitCpuM + podContainerLimitsCpuMi
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

func newClient(contextName string) (kubernetes.Interface, error) {
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
