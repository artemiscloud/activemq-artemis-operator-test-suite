# activemq-artemis-operator-test-suite
The ActiveMQ Artemis Operator Test Suite.

## Quick Start
Install Ginkgo, see https://onsi.github.io/ginkgo/
```
go install github.com/onsi/ginkgo/ginkgo@v1.16.5
```

Set up a kubernetes cluster, see https://kind.sigs.k8s.io/docs/user/quick-start/
```
kind create cluster
```

Export KUBECONFIG, see https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/
```
export KUBECONFIG=$HOME/.kube/config
```

Clone activemq-artemis-operator, see https://github.com/artemiscloud/activemq-artemis-operator
```
git clone https://github.com/artemiscloud/activemq-artemis-operator.git
```

### Use artemiscloud container images
Exexute the test suite using artemiscloud container images from quay.io and the local clone of the activemq-artemis-operator repository:
```
ginkgo -r -keepGoing test/smoke/basic/... -- \
        -broker-name activemq-artemis \
        -broker-image quay.io/artemiscloud/activemq-artemis-broker-kubernetes:dev.latest \
        -broker-image-second quay.io/artemiscloud/activemq-artemis-broker-kubernetes:dev.latest \
        -operator-image quay.io/artemiscloud/activemq-artemis-operator:dev.latest \
        -repository ../activemq-artemis-operator/deploy \
        -v2 -debug-run
```

### Use local container images
The default operator image pull policy is `Always`, change it to use a local operator container image, see https://kubernetes.io/docs/concepts/containers/images/#updating-images
```
sed -i 's/imagePullPolicy:.*/imagePullPolicy: IfNotPresent/' ../activemq-artemis-operator/deploy/operator.yaml
```

Load local container images into the kubernetes cluster, see https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster
```
kind load docker-image activemq-artemis-broker-kubernetes:dev.latest
```

Exexute the test suite using local container images and the local clone of activemq-artemis-operator repository
```
ginkgo -r -keepGoing test/smoke/basic/... -- \
        -broker-name activemq-artemis \
        -broker-image activemq-artemis-broker-kubernetes:dev.latest \
        -broker-image-second activemq-artemis-broker-kubernetes:dev.latest \
        -operator-image activemq-artemis-operator:dev.latest \
        -repository ../activemq-artemis-operator/deploy \
        -v2 -debug-run
```

## Delete test suite namespaces
```
kubectl get namespaces | grep -o '^e2e[^ ]*' | while read NAMESPACE; do kubectl delete namespaces $NAMESPACE; done
```

## Troubleshooting
Use the following command to investigate on kubernetes cluster failures:
- kubectl get pods --all-namespaces
- kubectl describe pod
- kubectl logs
