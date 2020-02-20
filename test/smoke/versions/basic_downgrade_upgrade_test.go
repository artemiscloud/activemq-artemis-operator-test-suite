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
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName)
	)

	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName)
	})

	ginkgo.It("Deploy broker and downgrade it to another version", func() {
		gomega.Expect(dw.DeployBrokers(1)).To(gomega.BeNil())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersion))
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageNameOld).ChangeImage())
		podLog, _ = ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersionOld))
	})

	ginkgo.It("Deploy broker and upgrade it to another version", func() {
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageNameOld).DeployBrokers( 1)).To(gomega.BeNil())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersionOld))
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageName).ChangeImage())
		podLog, _ = ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersionOld))
	})

	ginkgo.It("Deploy bogus image and replace it with broker", func() {
		gomega.Expect(dw.WithCustomImage("non-url-at-all").DeployBrokers( 1)).To(gomega.BeNil())
		// Non-image should result in pod not being created and flawlessly replaced with proper image
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageName).ChangeImage())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersionOld))
	})

	ginkgo.It("Deploy wrong image and replace it with broker", func() {
		gomega.Expect(dw.WithCustomImage(test.TestConfig.OperatorImageName).DeployBrokers( 1)).To(gomega.BeNil())
		// Pod in failed state should be replaced by new pod with proper broker
		gomega.Expect(dw.WithCustomImage(test.TestConfig.BrokerImageName).ChangeImage())
		podLog, _ := ctx1.GetLogs("ex-aao-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.TestConfig.BrokerVersionOld))
	})

	ginkgo.It("Deploy broker and upgrade it to another version", func() {

	})
})
