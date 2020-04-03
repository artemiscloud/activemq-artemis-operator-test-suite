package messaging

import (
	"github.com/onsi/ginkgo"
	brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/operators"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "messaging"
)

var (
	// Framework instance that holds the generated resources
	Framework *framework.Framework
	// Basic Operator instance
	brokerOperator operators.OperatorSetup
	brokerClient   brokerclientset.Interface
)

// Create the Framework instance to be used oneinterior test
var _ = ginkgo.BeforeEach(func() {
	// Setup the topology
	builder := operators.SupportedOperators[operators.OperatorTypeBroker]
	//Set image to parameter if one is supplied, otherwise use default from shipshape.
	if len(test.Config.OperatorImageName) != 0 {
		builder.WithImage(test.Config.OperatorImageName)
	}
	if test.Config.DownstreamBuild {
		builder.WithCommand("/home/amq-broker-operator/bin/entrypoint")
	}
	Framework = framework.NewFrameworkBuilder("broker-framework").
		WithBuilders(builder).
		Build()
	brokerOperator = Framework.GetFirstContext().OperatorMap[operators.OperatorTypeBroker]
	brokerClient = brokerOperator.Interface().(brokerclientset.Interface)
}, 60)

var _ = ginkgo.JustBeforeEach(func() {

})

// After each test completes, run cleanup actions to save resources (otherwise resources will remain till
// all specs from this suite are done.
var _ = ginkgo.AfterEach(func() {
	/*	if (test.TestConfig.DebugRun) {
			log.Logf("Not removing namespace due to debug option")
		} else {
			Framework.AfterEach()
		}*/
})
