package kustomization

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type ExampleV1Alpha1Interface interface {
	Kustomizations(namespace string) KustomizationInterface
}

type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: "kustomize.toolkit.fluxcd.io", Version: "v1beta2"}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func (c *ExampleV1Alpha1Client) Kustomizations(namespace string) KustomizationInterface {
	return &kustomizationClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func NewInformer(clientSet *ExampleV1Alpha1Client, namespace string) cache.SharedIndexInformer {
	kustomizev1.SchemeBuilder.AddToScheme(scheme.Scheme)

	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.Kustomizations(namespace).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Kustomizations(namespace).Watch(lo)
			},
		},
		&kustomizev1.Kustomization{},
		0,
		cache.Indexers{},
	)
}
