package hpa_test

import (
	"context"
	"testing"

	//	"fmt"
	//	"log"
	"os"
	"path/filepath"

	//	"time"

	//	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"k8s.io/apimachinery/pkg/fields"
	//	"k8s.io/apimachinery/pkg/util/runtime"
	//	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	//	"k8s.io/client-go/rest"
	//	"k8s.io/client-go/tools/cache"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
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

func TestHPA(t *testing.T) {
	clientSet, err := createClientsetFromLocal()
	if err != nil {
		t.Error(err)
	}
	clientSet.AutoscalingV1().HorizontalPodAutoscalers().Create(
		context.TODO(),
		autoscalingv1.HorizontalPodAutoscaler{
			// metav1.TypeMeta{
			// 	APIVersion: "autoscaling/v1",
			// 	Kind:       "HorizontalPodAutoscaler",
			// },
			Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
					Kind:"Service",
					Name:"hello-kn-client-2"
				}
			},
		},
		metav1.CreateOptions{},
	)
}
