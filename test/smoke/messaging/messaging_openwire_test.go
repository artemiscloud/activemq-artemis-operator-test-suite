package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("MessagingBasicTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		dw *test.DeploymentWrapper
		//	sender   amqp.Client
		receiver amqp.Client
		//url      string
		srw *test.SenderReceiverWrapper
	)

	// URL example: https://ex-aao-amqp-0-svc-rte-broker-operator-nd-ssl.apps.ocp43-released.broker-rvais-stable.fw.rhcloud.com
	const (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = "5672"
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
		Protocol      = "openwire"
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

		sendUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		receiveUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)

	})

	ginkgo.It("Deploy single broker instance and send/receive messages", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()

		testBaseSendReceiveSingleBroker(dw, srw, MessageCount, MessageBody, test.OpenwireAcceptor, 1)

	})

	ginkgo.It("Deploy double broker instances, send messages", func() {
		testBaseSendReceiveSingleBroker(dw, srw, MessageCount, MessageBody, test.OpenwireAcceptor, 2)
	})

	ginkgo.It("Deploy single broker, Send 1k of 1MB messages", func() {
		testSizedMessage(srw, dw, receiver, 1024*1024, 1024)
	})

	ginkgo.It("Deploy single broker, send 100 of 10MB messages", func() {
		testSizedMessage(srw, dw, receiver, 10*1024*1024, 100)
	})

	ginkgo.It("Deploy single broker, send 1000000 of 1kb messages", func() {
		testSizedMessage(srw, dw, receiver, 1024, 1024*1024)
	})
})
