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
		dw test.DeploymentWrapper
	)

	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName)

	})

	ginkgo.It("Create invalid url for deployment", func() {
		gomega.Expect(dw.WithWait(false).WithCustomImage( "gibberish://whatever").DeployBrokers(1)).To(gomega.BeNil())
		gomega.Expect(dw.WithWait(true).WithCustomImage(test.TestConfig.BrokerImageName).DeployBrokers(1)).To(gomega.BeNil())
	})

	ginkgo.It("Deploy broker and upgrade it to another version", func() {

	})
})
