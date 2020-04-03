package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"time"
)

var _ = ginkgo.Describe("MessagingMigrationTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		dw       test.DeploymentWrapper
		srw      test.SenderReceiverWrapper
		sender   amqp.Client
		receiver amqp.Client
	)

	const (
		MessageBody  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount = 100
	)

	// PrepareNamespace after framework has been created. Doesn'
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).
			WithContext(ctx1).WithCustomImage(test.Config.BrokerImageName).
			WithPersistence(true).WithMigration(true).
			WithName(DeployName)
		srw = test.SenderReceiverWrapper{}.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount)

	})

	ginkgo.It("Deploy double broker instance, migrate to single", func() {
		err := dw.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())

		sendUrl := "amqp://ex-aao-ss-1:5672/"
		receiveUrl := "amqp://ex-aao-ss-0:5672/"

		sender, receiver = srw.
			WithReceiveUrl(receiveUrl).
			WithSendUrl(sendUrl).
			PrepareSenderReceiver()
		_ = sender.Deploy()
		sender.Wait()

		senderResult := sender.Result()
		gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
		_ = dw.Scale(1)

		//Wait for a drainer pod to do its deed
		err = framework.WaitForDeployment(ctx1.Clients.KubeClient, ctx1.Namespace, "drainer", 1, time.Second*10, time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate everything to single", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		sendUrls := []string{"amqp://ex-aao-ss-3:5672/", "amqp://ex-aao-ss-2:5672/", "amqp://ex-aao-ss-1:5672/", "amqp://ex-aao-ss-0:5672/"}
		receiveUrl := "amqp://ex-aao-ss-0:5672/"

		receiver = srw.
			WithReceiveUrl(receiveUrl).
			PrepareReceiver()

		err := dw.DeployBrokers(4)
		for _, url := range sendUrls {
			sender = srw.WithSendUrl(url).PrepareSender()
			_ = sender.Deploy()
			sender.Wait()
		}
		//Scale to 1
		err = dw.Scale(1)
		gomega.Expect(err).To(gomega.BeNil())

		err = framework.WaitForDeployment(ctx1.Clients.KubeClient, ctx1.Namespace, "drainer", 1, time.Second*10, time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount * len(sendUrls)))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate last one ", func() {
		sendUrl := "amqp://ex-aao-ss-3:5672/"
		receiveUrl := "amqp://ex-aao-ss-0:5672/"
		sender, receiver = srw.
			WithReceiveUrl(receiveUrl).
			WithSendUrl(sendUrl).
			PrepareSenderReceiver()
		err := dw.DeployBrokers(4)
		gomega.Expect(err).To(gomega.BeNil())

		_ = sender.Deploy()
		sender.Wait()
		_ = dw.Scale(3)
		//Wait for a drainer pod to do its deed
		err = framework.WaitForDeployment(ctx1.Clients.KubeClient, ctx1.Namespace, "drainer", 1, time.Second*10, time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}

	})

})
