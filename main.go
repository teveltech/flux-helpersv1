package main

import (
	"encoding/json"
	"flag"
	"fmt"
	gitRepositoryClient "github.com/teveltech/flux-helpers/clientset/gitrepository"
	kustomizationClient "github.com/teveltech/flux-helpers/clientset/kustomization"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	kubeconfig = "/Users/orshefi/workspace/kubeconfig_itzik"

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	kustomizationClientSet, err := kustomizationClient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gitClientSet, err := gitRepositoryClient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	defer close(stop)

	kustomizationInformer := kustomizationClient.NewInformer(kustomizationClientSet, "default")

	kustomizationInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println(ObjToJson(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println(ObjToJson(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println(ObjToJson(obj))
		},
	})

	gitRepositoriesInformer := gitRepositoryClient.NewInformer(gitClientSet, "default")

	gitRepositoriesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println(ObjToJson(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println(ObjToJson(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println(ObjToJson(obj))
		},
	})

	go kustomizationInformer.Run(stop)
	go gitRepositoriesInformer.Run(stop)
	<-stop
}

func ObjToJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
