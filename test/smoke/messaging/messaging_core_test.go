package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"strconv"
)

var _ = ginkgo.Describe("MessagingCoreBasicTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		dw *test.DeploymentWrapper
		//	sender   amqp.Client
		//	receiver amqp.Client
		//url      string
		srw *test.SenderReceiverWrapper
	)

	// URL example: https://ex-aao-amqp-0-svc-rte-broker-operator-nd-ssl.apps.ocp43-released.broker-rvais-stable.fw.rhcloud.com
	var (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = int64(test.AcceptorPorts[test.CoreAcceptor])
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
		Protocol      = "tcp"
		ProtocolName  = test.CORE
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

		sendUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)

	})

	ginkgo.It("Deploy single broker instance and send/receive messages", func() {
		testBaseSendReceiveMessages(dw, srw, MessageCount, MessageBody, test.CoreAcceptor, 1, ProtocolName)

	})

	ginkgo.It("Deploy single amqp broker instance and send/receive core messages", func() {
		testBaseSendReceiveMessages(dw, srw, MessageCount, MessageBody, test.AmqpAcceptor, 1, ProtocolName)
	})

	ginkgo.It("Deploy single openwire broker instance and send/receive core messages", func() {
		testBaseSendReceiveMessages(dw, srw, MessageCount, MessageBody, test.OpenwireAcceptor, 1, ProtocolName)

	})

	ginkgo.It("Deploy double broker instances, send messages", func() {
		testBaseSendReceiveMessages(dw, srw, MessageCount, MessageBody, test.CoreAcceptor, 2, ProtocolName)
	})
})
