package versions

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("VersionsTests", func() {

	var (
		ctx1 *framework.ContextData
		dw   = &test.DeploymentWrapper{}
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = &test.DeploymentWrapper{}
		dw.WithWait(true).
			WithBrokerClient(brokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName)
	})

	ginkgo.It("Deploy broker and downgrade it to another version", func() {
		gomega.Expect(dw.DeployBrokers(1)).To(gomega.BeNil())
		podLog, _ := ctx1.GetLogs(DeployName + "-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.Config.BrokerVersion))
		gomega.Expect(dw.WithCustomImage(test.Config.BrokerImageNameOld).ChangeImage())
		podLog, _ = ctx1.GetLogs(DeployName + "-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.Config.BrokerVersionOld))
	})

	ginkgo.It("Deploy broker and upgrade it to another version", func() {
		gomega.Expect(dw.WithCustomImage(test.Config.BrokerImageNameOld).DeployBrokers(1)).To(gomega.BeNil())
		podLog, _ := ctx1.GetLogs(DeployName + "-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.Config.BrokerVersionOld))
		gomega.Expect(dw.WithCustomImage(test.Config.BrokerImageName).ChangeImage())
		podLog, _ = ctx1.GetLogs(DeployName + "-ss-0")
		gomega.Expect(podLog).To(gomega.ContainSubstring(test.Config.BrokerVersion))
	})
})
