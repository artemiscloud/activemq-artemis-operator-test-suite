package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"strconv"
)

var _ = ginkgo.Describe("MessagingAmqpBasicTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		bdw *test.BrokerDeploymentWrapper
		//url      string
		srw *test.SenderReceiverWrapper
	)

	var (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = int64(test.AcceptorPorts[test.AmqpAcceptor])
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
		Protocol      = test.AMQP
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		bdw = &test.BrokerDeploymentWrapper{}
		bdw.WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName)

		sendUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)
	})

	ginkgo.It("Deploy double broker instances, send messages", func() {
		testBaseSendReceiveMessages(bdw, srw, MessageCount, MessageBody, test.AmqpAcceptor, 2, Protocol)
	})

	ginkgo.It("Deploy single broker instances, send messages", func() {
		testBaseSendReceiveMessages(bdw, srw, MessageCount, MessageBody, test.AmqpAcceptor, 1, Protocol)
	})

	ginkgo.It("Deploy double instances with migration disabled, send messages, receive", func() {
		bdw.WithPersistence(true).WithMigration(false)
		testBaseSendReceiveMessages(bdw, srw, MessageCount, MessageBody, test.AmqpAcceptor, 2, Protocol)
	})

	ginkgo.It("Deploy double instances with migration disabled, send messages, scaledown, scaleup, receive", func() {
		bdw.WithPersistence(true).WithMigration(false)
		callback := func() (interface{}, error) {
			err := bdw.Scale(1)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			err = bdw.Scale(2)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			return nil, nil
		}
		testBaseSendReceiveMessagesWithCallback(bdw, srw, MessageCount, MessageBody, test.AmqpAcceptor, 2, Protocol, callback)
	})
})
