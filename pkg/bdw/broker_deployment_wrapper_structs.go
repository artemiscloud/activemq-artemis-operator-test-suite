package bdw
/* This file contains structs for BrokerDeploymentWrapper
 */


import (
	brokerclientset "github.com/artemiscloud/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

// BrokerDeploymentWrapper takes care of deployment of Broker
type BrokerDeploymentWrapper struct {
	wait           bool
	brokerClient   brokerclientset.Interface
	ctx1           *framework.ContextData
	customImage    string
	migration      bool
	persistence    bool
	name           string
	sslEnabled     bool
	exposeConsole  bool
	deploymentSize int
	isLtsDeployment bool
}

