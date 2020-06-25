// main
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

func watchPods(cs *kubernetes.Clientset) {

	addHandler := func(obj interface{}) {
		log.Println("Add Event")
		listPods(cs)
	}

	updateHandler := func(oldObj, newObj interface{}) {
		log.Println("Update Event")
		listPods(cs)
	}

	deleteHandler := func(obj interface{}) {
		log.Println("Delete Event")
		listPods(cs)
	}
	watchList := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", metav1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(
		watchList,
		&corev1.Pod{},
		time.Second*300,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addHandler,
			UpdateFunc: updateHandler,
			DeleteFunc: deleteHandler,
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
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
		cs, err := createClientsetFromLocal()
		if err != nil {
			panic(err)
		}
		//listPods(cs)
		watchPods(cs)
		time.Sleep(time.Second * 1000)
	}
	log.Println("End")
}
