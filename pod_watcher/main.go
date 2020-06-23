// main
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func listPods(cs *kubernetes.Clientset) {
	pods, err := cs.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func createClientsetFromLocal() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	path := filepath.Join(homeDir(), ".kube", "config")
	kubeconfig = &path

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, nil
}

func createClientsetFromPod() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, nil
}

func main() {
	log.Println("Start")
	for {
		cs, err := createClientsetFromPod()
		if err != nil {
			panic(err)
		}
		listPods(cs)
		time.Sleep(time.Second * 10)
	}
	log.Println("End")
}
