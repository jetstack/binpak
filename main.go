
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

	for _, node := range nodes.Items {
		fmt.Println(node.Name)
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