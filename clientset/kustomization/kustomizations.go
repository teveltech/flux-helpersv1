package kustomization

import (
	"context"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type KustomizationInterface interface {
	List(opts metav1.ListOptions) (*kustomizev1.KustomizationList, error)
	Get(name string, options metav1.GetOptions) (*kustomizev1.Kustomization, error)
	Create(*kustomizev1.Kustomization) (*kustomizev1.Kustomization, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type kustomizationClient struct {
	restClient rest.Interface
	ns         string
}

func (c *kustomizationClient) List(opts metav1.ListOptions) (*kustomizev1.KustomizationList, error) {
	result := kustomizev1.KustomizationList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kustomizations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *kustomizationClient) Get(name string, opts metav1.GetOptions) (*kustomizev1.Kustomization, error) {
	result := kustomizev1.Kustomization{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kustomizations").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *kustomizationClient) Create(kustomization *kustomizev1.Kustomization) (*kustomizev1.Kustomization, error) {
	result := kustomizev1.Kustomization{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("kustomizations").
		Body(kustomization).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *kustomizationClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kustomizations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
