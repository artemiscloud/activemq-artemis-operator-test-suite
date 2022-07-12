module github.com/artemiscloud/activemq-artemis-operator-test-suite

go 1.12

require (
	github.com/artemiscloud/activemq-artemis-operator v0.17.0
	github.com/fgiorgetti/qpid-dispatch-go-tests v0.0.0-20190923194420-c3f992ce0eee
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/onsi/ginkgo v1.12.3
	github.com/onsi/gomega v1.10.1
	github.com/rh-messaging/shipshape v0.2.8
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/tools v0.0.0-20200331025713-a30bf2db82d4 // indirect
	k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v1.0.0

)

// For local override of dependencies, use following:
// replace github.com/rh-messaging/activemq-artemis-operator v0.0.0+incompatible => ../../../github.com/rh-messaging/activemq-artemis-operator

// replace github.com/rh-messaging/shipshape v0.2.7 => ../../../github.com/rh-messaging/shipshape

replace bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d
