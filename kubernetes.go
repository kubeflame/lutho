package main

import (
	"context"

	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type CreateResource struct {
	Namespace string
	Options   metav1.CreateOptions
	GVR       schema.GroupVersionResource
	Object    *unstructured.Unstructured
}

func (cd *CreateResource) Create(client *dynamic.DynamicClient) (*unstructured.Unstructured, error) {
	return client.Resource(cd.GVR).Namespace(cd.Namespace).Create(context.TODO(), cd.Object, cd.Options)
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

type SelfSubjectAuth struct {
	Options metav1.CreateOptions
	Access  *authv1.SelfSubjectAccessReview
	Rules   *authv1.SelfSubjectRulesReview
}

func (ssa *SelfSubjectAuth) AccessReview(client kubernetes.Interface) (*authv1.SelfSubjectAccessReview, error) {
	return client.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), ssa.Access, ssa.Options)
}

func (ssa *SelfSubjectAuth) RulesReview(client kubernetes.Interface) (*authv1.SelfSubjectRulesReview, error) {
	return client.AuthorizationV1().SelfSubjectRulesReviews().Create(context.TODO(), ssa.Rules, ssa.Options)
}
