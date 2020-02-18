package versions

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1 *framework.ContextData
	)

	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
	})

	ginkgo.It("Deploy broker and downgrade it to another version", func() {
		gomega.Expect(test.DeployBrokers(ctx1, 1,brokerClient, test.BrokerImageName)).To(gomega.BeNil())
		// Check for jolokia call for version, curl from pod
		gomega.Expect(test.ChangeImage(ctx1,brokerClient, test.BrokerImageNameOld))
	})

	ginkgo.It("Deploy broker and upgrade it to another version", func() {
		gomega.Expect(test.DeployBrokers(ctx1, 1,brokerClient, test.BrokerImageNameOld)).To(gomega.BeNil())
		// Check for jolokia call for version, curl from pod
		gomega.Expect(test.ChangeImage(ctx1,brokerClient, test.BrokerImageName))
	})

	ginkgo.It("Deploy broker and upgrade it to another version", func() {

	})
})
