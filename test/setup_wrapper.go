package test

import (
	"flag"
	brokerclientset "github.com/artemiscloud/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"github.com/rh-messaging/shipshape/pkg/framework/operators"
)

type SetupWrapper struct {
	DeployName string
	BaseName   string
	Framework  *framework.Framework
	// Basic Operator instance
	BrokerOperator operators.OperatorSetup
	BrokerClient   brokerclientset.Interface
	odw            *OperatorDeploymentWrapper
}

func (sw *SetupWrapper) WithBaseName(baseName string) *SetupWrapper {
	sw.BaseName = baseName
	return sw
}

func (sw *SetupWrapper) WithDeployName(deployName string) *SetupWrapper {
	sw.DeployName = deployName
	return sw
}

func (sw *SetupWrapper) WithOperatorWrapper(odw *OperatorDeploymentWrapper) *SetupWrapper {
	sw.odw = odw
	return sw
}

func (sw *SetupWrapper) mayBeInitWrapper() {
	if sw.odw == nil {
		sw.odw = &OperatorDeploymentWrapper{}
	}
}

func (sw *SetupWrapper) SetOperatorEnv(vars map[string]string) {
	sw.mayBeInitWrapper()
	sw.odw.EnvVariables = vars
}

func (sw *SetupWrapper) AddOperatorEnv(name string, value string) {
	sw.mayBeInitWrapper()
	sw.odw.EnvVariables[name] = value
}

func (sw *SetupWrapper) BeforeEach() {
	sw.mayBeInitWrapper()
	builder := sw.odw.PrepareOperator()
	frBuilder := framework.NewFrameworkBuilder(sw.BaseName).
		WithBuilders(builder)
	if Config.Openshift {
		frBuilder = frBuilder.IsOpenshift(true)
	} else {
		log.Logf("Would be using namespaces")
	}
	sw.Framework = frBuilder.Build()
	sw.BrokerOperator = sw.Framework.GetFirstContext().OperatorMap[operators.OperatorTypeBroker]
	sw.BrokerClient = sw.BrokerOperator.Interface().(brokerclientset.Interface)
	log.Logf("We got: %v", Config.NeedsLatestCR)
}

func (sw *SetupWrapper) AfterEach() {
	if Config.DebugRun {
		log.Logf("Not removing namespace due to debug option")
	} else {
		sw.Framework.AfterEach()
	}
}

func (sw *SetupWrapper) JustBeforeEach() {
	//Nothing for now
}

func (sw *SetupWrapper) InitFlags() {
	RegisterFlags()
	flag.Parse()
}
