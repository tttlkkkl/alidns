# alidns
为 cert-manager 提供 alidns 的 DNS01 校验。

## 使用前提：
- 你已经安装好 cert-manager 并且开启了 webhook 支持 [install cert-manager](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html)。
- cert-manager 版本必须在0.8以上，本插件在0.9版本下测试通过,0.11下已确定无法使用。
- 你的域名通过阿里云DNS做解析，并已获得api权限。
## 安装
- `git clone git@github.com:tttlkkkl/alidns.git`
- `cd deploy`
- `helm install --name alidns --namespace cert-manager alidns/`
## 编排使用：
### 以为域名 lihuaio.com 签发通配符域名为例，开始之前你需要替换为自己的域名:
### 在 deploy 目录中的 k8s.yaml 是一个完整的例子，用 traefik 作为入口网关，显示 nginx 的欢迎页，可以作为参考。
- 创建阿里云 api 密钥:
你也可以与 `regionId` 同级别的直接设置 `accessKeyId`和`accessKeySecret`字段，但是不推荐这样做。
```yml
apiVersion: v1
kind: Secret
metadata:
  name: alibaba-api-dns-secret
data:
  accessKeyId: <your aliyun accessKeyId>
  accessKeySecret: <your aliyun accessKeySecret>
type: Opaque
```
- 创建 ClusterIssuer
```yaml
# ClusterIssuer
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: tttlkkkl@aliyun.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - selector: 
        dnsNames:
        - '*.lihuaio.com'
      dns01:
        webhook:
          config:
            accessKeyRef:
              accessKeySecretKey: accessKeySecret
              accessKeyIdKey: accessKeyId
              name: alibaba-api-dns-secret
            regionId: "cn-shenzhen"
            accessKeySecret:
            accessKeyId:
            ttl: 600
          groupName: acme.lihuaio.com
```
- 创建 Certificate 证书对象
```yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: lihuaio.com
  namespace: default
spec:
  secretName: lihuaio.com
  commonName: '*.lihuaio.com'
  dnsNames:
  - "*.lihuaio.com"
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
```
- 创建：ingress
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: lihuaio-ingress
  namespace: default
  annotations:
    certmanager.k8s.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - '*.lihuaio.com'
    secretName: lihuaio.com # 这个对应 Certificate 中的 secretName
  rules:
  - host: xx.lihuaio.com
    http:
      paths:
      - path: /
        backend:
          serviceName: backend-service
          servicePort: 80
```