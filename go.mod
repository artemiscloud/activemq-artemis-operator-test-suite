module gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang

go 1.12

require (
	github.com/PuerkitoBio/purell v1.1.1
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/jsonpointer v0.19.3
	github.com/go-openapi/jsonreference v0.19.3
	github.com/go-openapi/spec v0.19.5
	github.com/go-openapi/swag v0.19.6
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/google/btree v1.0.0
	github.com/google/gofuzz v1.0.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/hpcloud/tail v1.0.1-0.20180514194441-a1dbeea552b7
	github.com/imdario/mergo v0.3.7
	github.com/interconnectedcloud/qdr-operator v0.0.0-20200122133240-3984fddc8ad8 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/mailru/easyjson v0.7.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/openshift/api v0.0.0-20180801171038-322a19404e37
	github.com/openshift/client-go v0.0.0-20190412095722-0255926f5393
	github.com/pborman/uuid v1.2.0
	github.com/petar/GoLLRB v0.0.0-20190514000832-33fb24c13b99
	github.com/rh-messaging/activemq-artemis-operator v0.0.0
	github.com/rh-messaging/shipshape v0.0.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20200115085410-6d4e4cb37c7d
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20200113162924-86b910548bc1
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20200331025713-a30bf2db82d4
	google.golang.org/appengine v1.6.5
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/gengo v0.0.0-20200114144118-36b2048a9120
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools v0.1.8
	sigs.k8s.io/yaml v1.1.0

)

// For local override of dependencies, use following:
replace github.com/rh-messaging/activemq-artemis-operator v0.0.0 => ../../../github.com/rh-messaging/activemq-artemis-operator
