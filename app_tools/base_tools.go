package app_tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Port struct {
	Protocol   string
	Port       int32
	TargetPort int32
}

type SubSet struct {
	Name   string
	Labels map[string]string
	Weight int
}

type ReplicaSet struct {
	SetLabels map[string]string
	SubSets   []SubSet
}

type ServiceDef struct {
	NameSpace string
	Name      string
	ReplicaSet
	Port Port
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func GetLocalRestConfig() (*restclient.Config, error) {
	var kubeconfig *string
	path := filepath.Join(homeDir(), ".kube", "config")
	kubeconfig = &path

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, err
}

func createK8SClientset(restCfg *restclient.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(restCfg)
	return clientset, nil
}

func CreateK8SService(restCfg *restclient.Config, serviceDef *ServiceDef) error {
	clientset, err := createK8SClientset(restCfg)
	if err != nil {
		return err
	}

	deploy := appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
	}

	service := v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceDef.Name,
			Namespace: serviceDef.NameSpace,
			Labels:    serviceDef.ReplicaSet.SetLabels,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol: v1.Protocol(strings.ToUpper(serviceDef.Port.Protocol)),
					Port:     serviceDef.Port.Port,
					TargetPort: intstr.IntOrString{
						IntVal: serviceDef.Port.TargetPort,
					},
				},
			},
		},
	}
	clientset.CoreV1().Services(serviceDef.NameSpace).
		Create(context.TODO(), &service, metav1.CreateOptions{})
}

func CreateAppService(restCfg *restclient.Config, serviceDef *ServiceDef) error {

}
