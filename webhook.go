package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	certmanager_v1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/issuer/acme/dns/util"
	jsoniter "github.com/json-iterator/go"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	AliCloudAccessKeyIDRef     certmanager_v1alpha1.SecretKeySelector `json:"accessKeyIdRef"`
	DNSTtl                     int                                    `json:"ttl"`
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
	fmt.Println("set dns ...", StructToString(ch))
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = ch.ResolvedZone
	request.Type = "TXT"
	request.RR = getRR(ch.ResolvedFQDN)
	request.Value = ch.Key
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}
	request.TTL = requests.NewInteger(cfg.DNSTtl)
	client, err := a.getAliDNSClient(ch, cfg)
	if err != nil {
		return err
	}
	response, err := client.AddDomainRecord(request)

	if err != nil {
		log.Printf("Response:%v\n error: %v\n", response, err)
		return err
	}
	log.Printf("Response: %v\n", response)
	return nil
}

// StructToString formt struct data
func StructToString(s interface{}) string {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("%+v", s)
	}
	return string(b)
}

// CleanUp clean the dns setting
func (a *AlibabaDNSSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	fmt.Println("========================UNset dns ...", StructToString(ch))
	log.Printf("set the dns record,FQDN:%s,zone:%s\n", ch.ResolvedFQDN, ch.ResolvedZone)
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}
	client, err := a.getAliDNSClient(ch, cfg)
	if err != nil {
		return err
	}
	request := alidns.CreateDeleteSubDomainRecordsRequest()
	request.DomainName = ch.ResolvedZone
	request.RR = getRR(ch.ResolvedFQDN)
	request.Type = "TXT"
	response, err := client.DeleteSubDomainRecords(request)
	log.Printf("domain list :%v", response)
	if err != nil {
		log.Printf("delete fail :%v", err)
		return err
	}
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
func loadConfig(cfgJSON *extapi.JSON) (*AlibabaDNSSolverConfig, error) {
	cfg := AlibabaDNSSolverConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return &cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return &cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return &cfg, nil
}

func (a *AlibabaDNSSolver) getAliDNSClient(ch *v1alpha1.ChallengeRequest, cfg *AlibabaDNSSolverConfig) (*alidns.Client, error) {
	var err error
	var accessKeyID, accessKeySecret string
	accessKeySecret = cfg.AliCloudAccessKeySecret
	accessKeyID = cfg.AliCloudAccessKeyID
	if accessKeySecret == "" {
		if cfg.AliCloudAccessKeySecretRef.Key == "" {
			return nil, errors.New("the accessKeySecretRef key not found")
		}
		if cfg.AliCloudAccessKeySecretRef.Name == "" {
			return nil, errors.New("the accessKeySecretRef name not found")
		}
		secret, err := a.K8sClient.CoreV1().Secrets(ch.ResourceNamespace).Get(cfg.AliCloudAccessKeySecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		accessKeySecretRef, ok := secret.Data[cfg.AliCloudAccessKeySecretRef.Key]
		if !ok {
			return nil, errors.New("the accessKeySecret not found ")
		}
		accessKeySecret = fmt.Sprintf("%s", accessKeySecretRef)
	}
	if accessKeyID == "" {
		if cfg.AliCloudAccessKeyIDRef.Key == "" {
			return nil, errors.New("the accessKeyIdRef key not found")
		}
		if cfg.AliCloudAccessKeyIDRef.Name == "" {
			return nil, errors.New("the accessKeyIdRef name not found")
		}
		secret, err := a.K8sClient.CoreV1().Secrets(ch.ResourceNamespace).Get(cfg.AliCloudAccessKeyIDRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		accessKeySecretRef, ok := secret.Data[cfg.AliCloudAccessKeyIDRef.Key]
		if !ok {
			return nil, errors.New("the accessKeySecret not found ")
		}
		accessKeySecret = fmt.Sprintf("%s", accessKeySecretRef)
	}
	client, err := alidns.NewClientWithAccessKey(
		cfg.AliCloudRegionID,
		cfg.AliCloudAccessKeyID,
		accessKeySecret,
	)

	if err != nil {
		return nil, err
	}
	client.OpenLogger()
	return client, nil
}

func getRR(fqdn string) string {
	idx := strings.LastIndex(fqdn, ".")
	if idx == -1 {
		return util.UnFqdn(fqdn)
	}
	return fqdn[:idx]
}
