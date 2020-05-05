package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("MessagingBasicTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		dw       test.DeploymentWrapper
		sender   amqp.Client
		receiver amqp.Client
		//url      string
	)

	// URL example: https://ex-aao-amqp-0-svc-rte-broker-operator-nd-ssl.apps.ocp43-released.broker-rvais-stable.fw.rhcloud.com
	const (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = "5672"
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.
			WithWait(true).
			WithBrokerClient(brokerClient).
			WithContext(ctx1).
			WithCustomImage(test.Config.BrokerImageName).
			WithName(DeployName)

		sendUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		srw := test.SenderReceiverWrapper{}.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount)

		sender, receiver = srw.
			WithReceiveUrl(receiveUrl).
			WithSendUrl(sendUrl).
			PrepareSenderReceiver()

	})

	ginkgo.It("Deploy single broker instance and send/receive messages", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers(1)
		gomega.Expect(err).To(gomega.BeNil())

		_ = sender.Deploy()
		_ = receiver.Deploy()
		log.Logf("Started (sync) deployment of clients")
		sender.Wait()
		receiver.Wait()
		log.Logf("Wait finished")
		senderResult := sender.Result()
		receiverResult := receiver.Result()
		log.Logf("Finished (sync) deployment")
		log.Logf("Count sent: %d", senderResult.Delivered)
		log.Logf("Count received: %d", receiverResult.Delivered)
		log.Logf("Len of received: %d", len(receiverResult.Messages))
		gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))

		log.Logf("MessageCount is fine")
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy double broker instances, send messages", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())
		_ = sender.Deploy()
		_ = receiver.Deploy()
		sender.Wait()
		receiver.Wait()

		senderResult := sender.Result()
		receiverResult := receiver.Result()

		gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))

		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy broker with persistence but without migration", func() {
		err := dw.WithPersistence(false).WithMigration(false).DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())
	})

})
