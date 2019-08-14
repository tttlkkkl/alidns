package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	certmanager_v1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	jsoniter "github.com/json-iterator/go"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

// AlibabaDNSSolver interface realization
type AlibabaDNSSolver struct {
	K8sClient *kubernetes.Clientset
}

// AlibabaDNSSolverConfig alibaba DNS config
// see https://github.com/aliyun/alibaba-cloud-sdk-go
type AlibabaDNSSolverConfig struct {
	AliCloudRegionID           string                                 `json:"regionId"`
	AliCloudAccessKeyID        string                                 `json:"accessKeyId"`
	AliCloudAccessKeySecret    string                                 `json:"accessKeySecret"`
	AliCloudAccessKeySecretRef certmanager_v1alpha1.SecretKeySelector `json:"accessKeySecretRef"`
	DNSTtl                     *int                                   `json:"ttl"`
}

// NewAlibabaDNSSolver new the Solver
func NewAlibabaDNSSolver() *AlibabaDNSSolver {
	return &AlibabaDNSSolver{}
}

//NewAlibabaDNSSolverConfig new the config
func NewAlibabaDNSSolverConfig() *AlibabaDNSSolverConfig {
	return &AlibabaDNSSolverConfig{
		AliCloudAccessKeyID:     os.Getenv("ALICLOUD_ACCESS_KEY"),
		AliCloudAccessKeySecret: os.Getenv("ALICLOUD_SECRET_KEY"),
		AliCloudRegionID:        os.Getenv("REGIONID"),
	}
}

// Name return the name of dns solver
func (a *AlibabaDNSSolver) Name() string {
	return "alidns"
}

// Present handel the dns request
func (a *AlibabaDNSSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	fmt.Println("set dns ...", ch)
	_, err := getAliDNSClient(ch)
	if err != nil {
		return err
	}
	return nil
}

// CleanUp clean the dns setting
func (a *AlibabaDNSSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	fmt.Println("set dns ...", ch)
	// TODO: add code that deletes a record from the DNS provider's console
	return nil
}

// Initialize the init function
func (a *AlibabaDNSSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}
	a.K8sClient = cl
	return nil
}
func loadConfig(cfgJSON *extapi.JSON) (AlibabaDNSSolverConfig, error) {
	cfg := AlibabaDNSSolverConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}

func getAliDNSClient(ch *v1alpha1.ChallengeRequest) (*alidns.Client, error) {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return nil, err
	}
	var accessKeyID, accessKeySecret string
	accessKeyID = cfg.AliCloudAccessKeyID
	accessKeySecret = cfg.AliCloudAccessKeySecret
	log.Println(accessKeyID, accessKeySecret)
	client, err := alidns.NewClientWithAccessKey(
		cfg.AliCloudRegionID,
		accessKeyID,
		accessKeySecret,
	)

	if err != nil {
		return nil, err
	}
	client.OpenLogger()
	return client, nil
}
