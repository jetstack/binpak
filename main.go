
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type info struct {
	nodes []node
}

type node struct {
	name string
	capacity resources 
	workloads []workload
}

type resources struct {
	memoryMi int
	cpuM int
}

type workload struct {
	id string
	name string
	requests resources 
	limits resources
}

func main() {
	contextName := ""
	if len(os.Args) >= 2 {
		contextName = os.Args[1]
	}

	client, err := newClient(contextName)
	if err != nil {
		log.Fatal(err)
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	i := info{}

	for _, node := range nodes.Items {
		n := node{
			name: node.Name,
			capacity: resources{
				memoryMi: node.Status.Capacity.Memory,
				cpuM: node.Status.Capacity.Cpu,
			},
		}
		i.nodes = append(i.nodes, n)
		// fmt.Println(node.Name)
	}

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

func memoryToMemoryMi(memory string) int {
	
}

func cpuToCpuM(memory string) int {
	
}

// [
// 		{
// 		  id: "nodepool-id",
// 		  availableResources: {
// 			memory: 3072,
// 			cpu: 4000,
// 		  },
// 		  nodes: [
// 			{
// 			  id: "node-id",
// 			  workloads: [
// 				{
// 				  id: "workload-id",
// 				  name: "workload-name",
// 				  allocatedResources: {
// 					memory: 64,
// 					cpu: 250,
// 				  },
// 				  resourcesLimits: {
// 					memory: 128,
// 					cpu: 500,
// 				  },
// 				},
// 			  ],
// 			},
// 		  ],
// 		},
// 	  ]
// }