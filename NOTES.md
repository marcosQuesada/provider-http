helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane

10542  docker build -t eu.gcr.io/king-container-registry-shared/cloud-platform/crossplane/function-sequencer:0.0.1 .
10543  docker push eu.gcr.io/king-container-registry-shared/cloud-platform/crossplane/function-sequencer:0.0.1



    apiVersion: pkg.crossplane.io/v1
    kind: Provider
    metadata:
      name: provider-http
    spec:
      package: "xpkg.upbound.io/crossplane-contrib/provider-http:v0.2.0"


crossplane --verbose xpkg push -f /home/marcos/code/king/crossplane/providers/provider-king/_output/xpkg/linux_amd64/provider-king-v0.0.0-175.g3c51016.xpkg eu.gcr.io/king-container-registry-shared/cloud-platform/crossplane/provider-poc:0.0.11


172.17.0.1

❯ curl http://127.0.0.1:8081/api/v1/pets

[{"color":"redxx","id":6661,"name":"foxy","price":1666631,"state":"Created","tag":"MyTag"}]


❯ curl -X 'POST' http://127.0.0.1:8081/api/v1/pets -H 'Content-Type: application/json' -d '{"color":"foo", "id":999, "name":"foo-pet", "price":102909, "state":"fake-state", "tag":"dog"}'
{"color":"foo","id":999,"name":"foo-pet","price":102909,"state":"Created","tag":"dog"}

