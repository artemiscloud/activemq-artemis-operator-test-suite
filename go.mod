module gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang

go 1.12

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/interconnectedcloud/qdr-operator v0.0.0-20200122133240-3984fddc8ad8 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.10.0
	github.com/rh-messaging/activemq-artemis-operator v0.0.0+incompatible
	github.com/rh-messaging/shipshape v0.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/tools v0.0.0-20200331025713-a30bf2db82d4 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/klog v1.0.0

)

// For local override of dependencies, use following:
replace github.com/rh-messaging/activemq-artemis-operator v0.0.0+incompatible => ../../../github.com/rh-messaging/activemq-artemis-operator

replace github.com/rh-messaging/shipshape v0.0.0 => ../../../github.com/rh-messaging/shipshape
