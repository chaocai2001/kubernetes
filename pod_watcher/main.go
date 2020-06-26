// main
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var cnt int = 0
var cntUpdate int = 0

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
		cnt++
		log.Println("Add Event", cnt)
		listPods(cs)
	}

	updateHandler := func(oldObj, newObj interface{}) {
		cntUpdate++
		log.Println("Update Event", cntUpdate)
		log.Println("old:", oldObj)
		log.Println("new:", newObj)
		listPods(cs)
	}

	deleteHandler := func(obj interface{}) {
		log.Println("Delete Event")
		listPods(cs)
	}
	sel, err := fields.ParseSelector("status.phase=Running")
	if err != nil {
		panic(err)
	}
	watchList := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), string(corev1.ResourcePods), metav1.NamespaceAll,
		sel)
	_, controller := cache.NewInformer(
		watchList,
		&corev1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addHandler,
			UpdateFunc: updateHandler,
			DeleteFunc: deleteHandler,
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
}

func watchPodsShared(cs *kubernetes.Clientset) {
	addHandler := func(obj interface{}) {
		cnt++
		log.Println("Add Event", cnt)
		listPods(cs)
	}

	updateHandler := func(oldObj, newObj interface{}) {
		cntUpdate++
		log.Println("Update Event", cntUpdate)
		log.Println("old:", oldObj)
		log.Println("new:", newObj)
		listPods(cs)
	}

	deleteHandler := func(obj interface{}) {
		log.Println("Delete Event")
		listPods(cs)
	}
	factory := informers.NewSharedInformerFactory(cs, 0)

	// Get the informer for the right resource, in this case a Pod
	informer := factory.Core().V1().Pods().Informer()

	// Create a channel to stops the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// When a new pod gets created
		AddFunc: addHandler,
		// When a pod gets updated
		UpdateFunc: updateHandler,
		// When a pod gets deleted
		DeleteFunc: deleteHandler,
	})

	// You need to start the informer, in my case, it runs in the background
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func watchPodsShared2(cs *kubernetes.Clientset) {
	addHandler := func(obj interface{}) {
		cnt++
		log.Println("Add Event", cnt)
		listPods(cs)
	}

	updateHandler := func(oldObj, newObj interface{}) {
		cntUpdate++
		log.Println("Update Event", cntUpdate)
		log.Println("old:", oldObj)
		log.Println("new:", newObj)
		listPods(cs)
	}

	deleteHandler := func(obj interface{}) {
		log.Println("Delete Event")
		listPods(cs)
	}
	factory := informers.NewSharedInformerFactory(cs, 0)

	// Get the informer for the right resource, in this case a Pod
	informer := factory.Core().V1().Pods().Informer()

	// Create a channel to stops the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// When a new pod gets created
		AddFunc: addHandler,
		// When a pod gets updated
		UpdateFunc: updateHandler,
		// When a pod gets deleted
		DeleteFunc: deleteHandler,
	})

	// You need to start the informer, in my case, it runs in the background
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
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
		//watchPodsShared(cs)
		time.Sleep(time.Second * 10000000)

	}
	log.Println("End")
}
