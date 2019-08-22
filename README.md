# alidns
[简体中文](README-zh.md)

the alidns provider solver for k8s cert-manager.
cert-manager version >= 0.8


## install
- `git clone git@github.com:tttlkkkl/alidns.git`
- `cd deploy`
- `helm install --name alidns --namespace cert-manager alidns/`

[example](deploy/k8s.yml),the example use the domain lihuaio.com,you need to replace it with your own.

run scripts/fetch-test-binaries.sh install the kubebuilder tools.
run ` make test `
