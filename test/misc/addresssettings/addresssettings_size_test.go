package addresssettings

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"strconv"
)

var _ = ginkgo.Describe("AddressSettingsSizeTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		brokerDeployer *bdw.BrokerDeploymentWrapper
		//url      string
		srw *test.SenderReceiverWrapper
	)

	var (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = int64(bdw.AcceptorPorts[bdw.AmqpAcceptor])
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
		Protocol      = test.AMQP
	)

	ginkgo.BeforeEach(func() {
		if brokerDeployer != nil {
			brokerDeployer.PurgeAddressSettings()
		}
	})

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		setEnv(ctx1, brokerDeployer)
		sendURL := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveURL := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendURL).
			WithReceiveUrl(receiveURL)

		brokerDeployer.SetUpDefaultAddressSettings(AddressBit)

	})

	ginkgo.It("maxSizeBytes limit 1KB positive scenarios - DROP", func() {
		brokerDeployer.WithAddressSize(AddressBit, "1K").WithAddressPolicy(AddressBit, bdw.DropPolicy)
		test_helpers.TestBaseSendMessages(brokerDeployer, srw, 50, MessageBody, bdw.AmqpAcceptor, 1, Protocol, "sender", nil)
	})

	ginkgo.It("maxSizeBytes limit 1KB positive scenarios - PAGE", func() {
		brokerDeployer.WithAddressSize(AddressBit, "1K").WithAddressPolicy(AddressBit, bdw.PagePolicy)
		err := test_helpers.TestBaseSendMessages(brokerDeployer, srw, 50, MessageBody, bdw.AmqpAcceptor, 1, Protocol, "sender", nil)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("maxSizeBytes limit 10KB - DROP", func() {
		brokerDeployer.WithAddressSize(AddressBit, "10K").WithAddressPolicy(AddressBit, bdw.DropPolicy)
		err := test_helpers.TestBaseSendMessages(brokerDeployer, srw, 200, MessageBody, bdw.AmqpAcceptor, 1, Protocol, "sender", nil)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("maxSizeBytes not affecting other addresses", func() {
		brokerDeployer.WithAddressSize(AddressBit+"someOtherAddress", "1K").WithAddressPolicy(AddressBit, bdw.DropPolicy)
		err := test_helpers.TestBaseSendMessages(brokerDeployer, srw, 200, MessageBody, bdw.AmqpAcceptor, 1, Protocol, "sender", nil)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("maxSizeBytes limit 1KB negative - FAIL", func() {
		brokerDeployer.WithAddressSize(AddressBit, "1K").WithAddressPolicy(AddressBit, bdw.FailPolicy)
		err := test_helpers.TestBaseSendReceiveMessages(brokerDeployer, srw, 200, MessageBody, bdw.AmqpAcceptor, 1, Protocol)
		gomega.Expect(err).NotTo(gomega.BeNil()) //ToDo: error validation through logs!
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("is not equal to sent/delivered"))
		log.Logf("Expected error received: %s", err.Error())
	})

	ginkgo.It("maxSizeBytes limit 1KB negative - BLOCK", func() {
		brokerDeployer.WithAddressSize(AddressBit, "1K").WithAddressPolicy(AddressBit, bdw.BlockPolicy)
		// Should be already filled.
		err := test_helpers.TestBaseSendReceiveMessages(brokerDeployer, srw, 200, MessageBody, bdw.AmqpAcceptor, 1, Protocol)
		gomega.Expect(err).NotTo(gomega.BeNil()) //ToDo: error validation through logs!
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("is not equal to sent/delivered"))
		log.Logf("Expected error received: %s", err.Error())
	})

	ginkgo.It("maxSizeBytes limit exact message size (27b) - DROP", func() {
		brokerDeployer.WithAddressSize(AddressBit, "27").WithAddressPolicy(AddressBit, bdw.DropPolicy)
		err := test_helpers.TestBaseSendMessages(brokerDeployer, srw, 1, MessageBody, bdw.AmqpAcceptor, 1, Protocol, "sender", nil)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("maxSizeBytes limit by address regexp - DROP", func() {
		brokerDeployer.WithAddressSize("some*", "1K").WithAddressPolicy(AddressBit, bdw.DropPolicy)
		err := test_helpers.TestBaseSendMessages(brokerDeployer, srw, 1, MessageBody, bdw.AmqpAcceptor, 1, Protocol, "sender", nil)
		gomega.Expect(err).To(gomega.BeNil())
	})
})
