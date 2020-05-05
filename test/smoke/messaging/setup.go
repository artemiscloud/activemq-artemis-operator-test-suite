package messaging

import (
	"errors"
	"github.com/onsi/ginkgo"
	brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/events"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"github.com/rh-messaging/shipshape/pkg/framework/operators"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	v1 "k8s.io/api/core/v1"
	"strings"
	"time"
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
	builder := test.PrepareOperator()
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
	if test.Config.DebugRun {
		log.Logf("Not removing namespace due to debug option")
	} else {
		Framework.AfterEach()
	}
})

func formUrl(number, subdomain, namespace, domain, address, port string) string {
	//messaging-ss-0.messaging-hdls-svc.e2e-tests-broker-framework-msp9m.svc.cluster.local
	return "amqp://" + DeployName + "-ss-" + number + "." + DeployName + subdomain + "." + namespace + "." + domain + ":" + port +
		"/" + address
}

func WaitForDrainerRemoval(count int) {
	podsRemoved := 0
	Framework.GetFirstContext().EventHandler.AddEventHandler(events.Pod, events.Delete, func(obj ...interface{}) {
		podObj := obj[0].(v1.Pod)
		log.Logf("Pod %s has been removed", podObj.Name)
		if strings.Contains(podObj.Name, "ss") {
			podsRemoved += 1
			log.Logf("Pod %s has completed draining", podObj.Name)
		}
	})
	log.Logf("Waiting for drainer pods to execute")
	for podsRemoved < count {
		i := 0
		time.Sleep(time.Second * 5)
		i++
		if i > 60 {
			panic(errors.New("drainer pod failed to be created/removed over 5 minutes, exiting test"))
		}
	}

}
