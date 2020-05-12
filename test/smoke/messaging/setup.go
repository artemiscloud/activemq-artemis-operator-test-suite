package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"github.com/rh-messaging/shipshape/pkg/framework/operators"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"index/suffixarray"
	"regexp"
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
	return "amqp://" + DeployName + "-ss-" + number + "." + DeployName + subdomain + "." + namespace + "." + domain + ":" + port +
		"/" + address
}

func WaitForDrainerRemovalSlow(count int, timeout time.Duration, retries int) bool {
	expectedLog := "Deleting drain pod"
	loop := 0
	r := regexp.MustCompile(expectedLog)
	label := "amq-broker-operator"
	if !test.Config.DownstreamBuild {
		label = "activemq-artemis-operator"
	}
	operatorPodName, err := Framework.GetFirstContext().GetPodName(label)
	log.Logf("loading logs from pod %s", operatorPodName)
	gomega.Expect(err).To(gomega.BeNil())
	for loop < retries {
		if loop%10 == 0 {
			log.Logf("(still) waiting for drainer completion")
		}
		operatorLog, _ := Framework.GetFirstContext().GetLogs(operatorPodName)
		if strings.Contains(operatorLog, expectedLog) {
			index := suffixarray.New([]byte(operatorLog))
			results := index.FindAllIndex(r, -1)
			if len(results) == count {
				return true
			}
		}
		time.Sleep(timeout)
		loop++
	}
	return false
}

// WaitForDrainerRemoval would check logs for amount of drainer finished messages.
// Wait for up to 60 seconds * count
// Returns true when found all drainers expected, false otherwise
func WaitForDrainerRemoval(count int) bool {
	return WaitForDrainerRemovalSlow(count, time.Second*time.Duration(10), count*6)
}
