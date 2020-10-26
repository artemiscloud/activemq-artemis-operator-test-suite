package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/rh-messaging/shipshape/pkg/framework"
	bdw "gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/pkg/bdw"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"strconv"
)

var _ = ginkgo.Describe("MessagingAllAcceptorTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		brokerDeployer *bdw.BrokerDeploymentWrapper
		//	sender   amqp.Client
		//	receiver amqp.Client
		//url      string
		srw *test.SenderReceiverWrapper
	)

	// URL example: https://ex-aao-amqp-0-svc-rte-broker-operator-nd-ssl.apps.ocp43-released.broker-rvais-stable.fw.rhcloud.com
	var (
		MessageBody          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount         = 100
		Port                 = int64(bdw.AcceptorPorts[bdw.AllAcceptor])
		Domain               = "svc.cluster.local"
		SubdomainName        = "-hdls-svc"
		AddressBit           = "someQueue"
		Protocol             = "tcp"
		ProtocolAmqp         = "amqp"
		ProtocolNameOpenwire = test.OPENWIRE
		ProtocolNameAmqp     = test.AMQP
		ProtocolNameCore     = test.CORE
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		brokerDeployer.WithWait(true).
			WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName).
			WithLts(!test.Config.NeedsV2)

		sendUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)

	})

	ginkgo.It("Deploy single broker instance and send/receive messages through openwire", func() {
		testBaseSendReceiveMessages(brokerDeployer, srw, MessageCount, MessageBody, bdw.AllAcceptor, 1, ProtocolNameOpenwire)
	})

	ginkgo.It("Deploy single broker instance and send/receive messages through amqp", func() {
		sendUrl := test.FormUrl(ProtocolAmqp, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := test.FormUrl(ProtocolAmqp, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)
		testBaseSendReceiveMessages(brokerDeployer, srw, MessageCount, MessageBody, bdw.AllAcceptor, 1, ProtocolNameAmqp)
	})

	ginkgo.It("Deploy single broker instance and send/receive messages through core", func() {
		sendUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)
		testBaseSendReceiveMessages(brokerDeployer, srw, MessageCount, MessageBody, bdw.AllAcceptor, 1, ProtocolNameCore)
	})

	ginkgo.It("Deploy single broker instance and send messages through openwire, receive through amqp", func() {
		_ = true
		sendUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := test.FormUrl(ProtocolAmqp, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)
		testBaseSendMessages(brokerDeployer, srw, MessageCount, MessageBody, bdw.AllAcceptor, 1, ProtocolNameOpenwire, "sender-openwire", nil)
		testBaseReceiveMessages(brokerDeployer, srw, MessageCount, MessageBody, ProtocolNameAmqp)
	})
})
