module github.com/artemiscloud/activemq-artemis-operator-test-suite

go 1.12

require (
	github.com/artemiscloud/activemq-artemis-operator v1.0.4
	github.com/fgiorgetti/qpid-dispatch-go-tests v0.0.0-20190923194420-c3f992ce0eee
	github.com/ghodss/yaml v1.0.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
	github.com/rh-messaging/shipshape v0.2.7
	k8s.io/api v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/client-go v0.24.2
	k8s.io/klog v1.0.0

)

// For local override of dependencies, use following:
// replace github.com/rh-messaging/activemq-artemis-operator v0.0.0+incompatible => ../../../github.com/rh-messaging/activemq-artemis-operator

replace github.com/rh-messaging/shipshape v0.2.7 => ../../../github.com/rh-messaging/shipshape

replace bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d
