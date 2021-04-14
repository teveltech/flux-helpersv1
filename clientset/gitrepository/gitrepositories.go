package gitrepository

import (
	"context"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type GitRepositoryInterface interface {
	List(opts metav1.ListOptions) (*sourcev1.GitRepositoryList, error)
	Get(name string, options metav1.GetOptions) (*sourcev1.GitRepository, error)
	Create(*sourcev1.GitRepository) (*sourcev1.GitRepository, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type gitRepositoryClient struct {
	restClient rest.Interface
	ns         string
}

func (c *gitRepositoryClient) List(opts metav1.ListOptions) (*sourcev1.GitRepositoryList, error) {
	result := sourcev1.GitRepositoryList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("gitrepositories").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *gitRepositoryClient) Get(name string, opts metav1.GetOptions) (*sourcev1.GitRepository, error) {
	result := sourcev1.GitRepository{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("gitrepositories").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *gitRepositoryClient) Create(gitRepository *sourcev1.GitRepository) (*sourcev1.GitRepository, error) {
	result := sourcev1.GitRepository{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("gitrepositories").
		Body(gitRepository).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *gitRepositoryClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("gitrepositories").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
