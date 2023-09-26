package main

import (
	"fmt"

	v1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func initClusterConfig(authReq *AuthRequest) (config *rest.Config, err error) {
	switch authReq.Type {
	case KubernetesConfigType.inClusterConfig:
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	case KubernetesConfigType.kubeconfigPath:
		config, err = clientcmd.BuildConfigFromFlags("", authReq.KubeconfigPath)
		if err != nil {
			return nil, err
		}
	case KubernetesConfigType.kubeconfigRaw:
		cc, errcc := clientcmd.NewClientConfigFromBytes([]byte(authReq.KubeconfigRaw))
		if errcc != nil {
			return nil, errcc
		}

		restConfig, errRestConfig := cc.ClientConfig()
		if errRestConfig != nil {
			return nil, errRestConfig
		}

		config = restConfig
	case KubernetesConfigType.accessToken:
		config, err = clientcmd.BuildConfigFromFlags(authReq.MasterURL, "")
		if err != nil {
			return nil, err
		}
		config.BearerToken = authReq.AccessToken
		config.TLSClientConfig = rest.TLSClientConfig{Insecure: true}
	}

	return config, nil
}

func initClientset(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

func initDynamicClient(config *rest.Config) (*dynamic.DynamicClient, error) {
	return dynamic.NewForConfig(config)
}

func authInit(ar *APIResource) (err error) {
	cfg, err := initClusterConfig(ar.AuthRequest)
	if err != nil {
		return fmt.Errorf("could not initialize cluster config: %s", err)
	}
	ar.Config = cfg

	cs, err := initClientset(ar.Config)
	if err != nil {
		return fmt.Errorf("could not initialize kubernetes clientset: %s", err)
	}
	ar.Clientset = cs

	dc, err := initDynamicClient(ar.Config)
	if err != nil {
		return fmt.Errorf("could not initialize kubernetes dynamic client: %s", err)
	}
	ar.DynamicClient = dc

	hac, hes, err := initHelm(ar)
	if err != nil {
		return fmt.Errorf("could not initialize helm client: %s", err)
	}

	ar.Helm = &Helm{ActionConfig: hac, EnvSettings: hes}

	ssarr := &SelfSubjectAccessReviewResource{
		SSAR: &v1.SelfSubjectAccessReview{
			Spec: v1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &v1.ResourceAttributes{
					Verb:     "*",
					Resource: "*",
				},
			},
		},
	}

	_, errReview := ssarr.Review(ar.Clientset)
	if errReview != nil {
		return errReview
	}

	return nil
}
