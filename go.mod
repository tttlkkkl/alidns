module github.com/tttlkkkl/alidns

go 1.12

require (
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190813065001-bd59ef2e00ef
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/jetstack/cert-manager v0.9.0
	github.com/json-iterator/go v1.1.7
	k8s.io/apiextensions-apiserver v0.0.0-20190810101755-ebc439d6a67b
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190413052642-108c485f896e

replace github.com/evanphx/json-patch => github.com/evanphx/json-patch v0.0.0-20190203023257-5858425f7550
