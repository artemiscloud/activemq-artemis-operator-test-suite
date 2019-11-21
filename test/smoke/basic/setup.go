package basic

import (
	"github.com/onsi/ginkgo"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/operators"
	brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"


)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "basic"
	DeploySize = 1
	ImageName = "brew-pulp-docker01.web.prod.ext.phx2.redhat.com:8888/amq7/amq-broker-operator:0.9"
)

var (
	// Framework instance that holds the generated resources
	Framework *framework.Framework
	// Basic Operator instance
	brokerOperator operators.OperatorSetup
	brokerClient brokerclientset.Interface
)

// Create the Framework instance to be used oneinterior test
var _ = ginkgo.BeforeEach(func() {
	// Setup the topology
	Framework = framework.NewFrameworkBuilder("broker-framework").
		WithBuilders(operators.SupportedOperators[operators.OperatorTypeBroker]).
		Build()
	brokerOperator = Framework.GetFirstContext().OperatorMap[operators.OperatorTypeBroker]
	brokerClient = brokerOperator.Interface().(brokerclientset.Interface)
}, 60)

// Deploy Interconnect
var _ = ginkgo.JustBeforeEach(func() {

})

// After each test completes, run cleanup actions to save resources (otherwise resources will remain till
// all specs from this suite are done.
var _ = ginkgo.AfterEach(func() {
	Framework.AfterEach()
})
