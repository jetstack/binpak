
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"encoding/json"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type info struct {
	instances []instance
}

func (i *info) AddItem(item instance) []instance {
	i.instances = append(i.instances, item)
	return i.instances
  }

type instance struct {
	name string
	capacity resources 
	workloads []workload
}

type resources struct {
	memoryMi int64
	cpuM int64
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

	// pods, err := client.CoreV1().Pods("").List(context.TODO(), meta_v1.ListOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, pod := range pods.Items {
	// 	fmt.Println(pod.Name)
	// }

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	instanceList := make([]instance, 0)

	for _, node := range nodes.Items {
		nodeName := node.Name
		nodeMemory := node.Status.Capacity["memory"]
		nodeMemoryMi := nodeMemory.MilliValue()
		nodeCpu := node.Status.Capacity["cpu"]
		nodeCpuM := nodeCpu.MilliValue()


		n := instance{
			name: nodeName,
			capacity: resources{
				memoryMi: nodeMemoryMi,
				cpuM: nodeCpuM,
			},
		}
		fmt.Println(node.Name)

		instanceList = append(instanceList, n)
	}

	fmt.Println(instanceList)

	i := info{
		instances: instanceList,
	}


	b, err := json.Marshal(i)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(b))	

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