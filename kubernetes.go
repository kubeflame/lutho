package main

import (
	"context"
	"flag"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeconfigType int

const (
	Config KubeconfigType = iota
	AccessToken
)

type KubeconfigInit struct {
	MasterURL   string
	AccessToken string
	Type        KubeconfigType
	TLSInsecure bool
}

func initKubeconfig(ki *KubeconfigInit) (config *rest.Config, err error) {
	switch ki.Type {
	case Config:
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
		// Show current-context only when the config is created outside the cluster

		// config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		// 	&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig},
		// 	&clientcmd.ConfigOverrides{
		// 		CurrentContext: "",
		// 	}).RawConfig()

		// if err != nil {
		// 	fmt.Println(err)
		// }
		// currentContext := config.CurrentContext

		// fmt.Println("currentContext:", currentContext)
	case AccessToken:
		config, err = clientcmd.BuildConfigFromFlags(ki.MasterURL, "")
		if err != nil {
			return nil, err
		}
		config.BearerToken = ki.AccessToken
		config.TLSClientConfig = rest.TLSClientConfig{Insecure: ki.TLSInsecure}
	}

	return config, nil
}

func initClientset(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

func initDynamicClient(config *rest.Config) (*dynamic.DynamicClient, error) {
	return dynamic.NewForConfig(config)
}

type ListResource struct {
	Namespace string
	Options   metav1.ListOptions
	GVR       schema.GroupVersionResource
}

func (ld *ListResource) List(client *dynamic.DynamicClient) (*unstructured.UnstructuredList, error) {
	return client.Resource(ld.GVR).Namespace(ld.Namespace).List(context.TODO(), ld.Options)
}

func (ld *ListResource) Watch(client *dynamic.DynamicClient) (watch.Interface, error) {
	return client.Resource(ld.GVR).Namespace(ld.Namespace).Watch(context.TODO(), ld.Options)
}

type GetResource struct {
	Name      string
	Namespace string
	Options   metav1.GetOptions
	GVR       schema.GroupVersionResource
}

func (gr *GetResource) Get(client *dynamic.DynamicClient) (*unstructured.Unstructured, error) {
	resource, err := client.Resource(gr.GVR).Namespace(gr.Namespace).Get(context.TODO(), gr.Name, gr.Options)
	if err != nil {
		return nil, err
	}

	unstructured.RemoveNestedField(resource.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(resource.Object, "metadata", "managedFields")

	return resource, nil
}

type UpdateResource struct {
	Name          string
	Namespace     string
	Data          []byte
	UpdateOptions metav1.UpdateOptions
	GetOptions    metav1.GetOptions
	GVR           schema.GroupVersionResource
	GVK           *schema.GroupVersionKind
}

func (ur *UpdateResource) Update(client *dynamic.DynamicClient) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	if _, _, err := unstructured.UnstructuredJSONScheme.Decode(ur.Data, ur.GVK, obj); err != nil {
		return nil, err
	}

	resource, err := client.Resource(ur.GVR).Namespace(ur.Namespace).Update(context.TODO(), obj, ur.UpdateOptions)
	if err != nil {
		return nil, err
	}

	unstructured.RemoveNestedField(resource.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(resource.Object, "metadata", "managedFields")

	return resource, nil
}

type DeleteResource struct {
	Name      string
	Namespace string
	Options   metav1.DeleteOptions
	GVR       schema.GroupVersionResource
}

func (dr *DeleteResource) Delete(client *dynamic.DynamicClient) error {
	return client.Resource(dr.GVR).Namespace(dr.Namespace).Delete(context.TODO(), dr.Name, dr.Options)
}
