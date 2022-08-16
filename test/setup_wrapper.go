package test

import (
	"flag"
	"time"

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
	frBuilder := framework.NewFrameworkBuilder(sw.BaseName).WithBuilders(builder)
	if Config.Openshift {
		frBuilder = frBuilder.IsOpenshift(true)
	} else {
		log.Logf("would be using namespaces")
	}
	if Config.GlobalOperator {
		frBuilder = frBuilder.WithGlobalOperator(true)
	} else {
		log.Logf("would be using local operator installation")
	}
	sw.Framework = frBuilder.Build()
	sw.BrokerOperator = sw.Framework.GetFirstContext().OperatorMap[operators.OperatorTypeBroker]
	sw.BrokerClient = sw.BrokerOperator.Interface().(brokerclientset.Interface)
}

func (sw *SetupWrapper) StartOperator() error {
	deployment, err := sw.BrokerOperator.GetDeployment()
	kubeclient := sw.Framework.GetFirstContext().Clients.KubeClient
	if err != nil {
		return err
	}
	replicas := int32(1)
	deployment.Spec.Replicas = &replicas
	sw.BrokerOperator.UpdateDeployment(deployment)
	return framework.WaitForDeployment(kubeclient, sw.BrokerOperator.Namespace(), sw.BrokerOperator.Name(), 1, 5*time.Second, 120*time.Second)
}

func (sw *SetupWrapper) StopOperator() error {
	deployment, err := sw.BrokerOperator.GetDeployment()
	if err != nil {
		return err
	}
	replicas := int32(0)
	deployment.Spec.Replicas = &replicas
	err = sw.BrokerOperator.UpdateDeployment(deployment)
	if err != nil {
		return err
	}
	//flat wait because I can't figure an event to subscribe to for operator deletion
	time.Sleep(time.Second * 60)
	return nil
}

func (sw *SetupWrapper) RestartOperator() error {
	err := sw.StopOperator()
	if err != nil {
		return err
	}
	return sw.StartOperator()
}

func (sw *SetupWrapper) RedeployOperator() error {
	err := sw.BrokerOperator.DeleteDeployment()
	if err != nil {
		return err
	}
	kubeclient := sw.Framework.GetFirstContext().Clients.KubeClient
	// There is no obvious event I could subscribe to to receive feedback on deployment deletion.
	log.Logf("Waiting for 60 seconds for operator shutdown")
	time.Sleep(60 * time.Second)
	if err != nil {
		return err
	}
	err = sw.BrokerOperator.CreateDeployment()
	if err != nil {
		return err
	}
	return framework.WaitForDeployment(kubeclient, sw.BrokerOperator.Namespace(), sw.BrokerOperator.Name(), 1, 5*time.Second, 120*time.Second)
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
