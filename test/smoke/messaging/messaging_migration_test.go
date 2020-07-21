package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("MessagingMigrationTests", func() {

	var (
		ctx1     *framework.ContextData
		bdw      *test.BrokerDeploymentWrapper
		srw      *test.SenderReceiverWrapper
		sender   amqp.Client
		receiver amqp.Client
	)

	const (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 10
		Port          = "5672"
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
		Protocol      = test.AMQP
	)

	// PrepareNamespace after framework has been created.
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		bdw = &test.BrokerDeploymentWrapper{}
		bdw.WithWait(true).WithBrokerClient(sw.BrokerClient).
			WithContext(ctx1).WithCustomImage(test.Config.BrokerImageName).
			WithPersistence(true).WithMigration(true).
			WithName(DeployName)
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount)
	})

	ginkgo.JustAfterEach(func() {
		log.Logf("Test failures in this suite might be related to ENTMQBR-3597")
		sw.Framework.GetFirstContext().EventHandler.ClearCallbacks()
	})

	// This test might fail due to ENTMQBR-3597
	ginkgo.It("Deploy double broker instance, migrate to single", func() {
		err := bdw.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())

		sendUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)

		sender, receiver := srw.
			WithReceiveUrl(receiveUrl).
			WithSendUrl(sendUrl).
			PrepareSenderReceiverWithProtocol(test.AMQP)

		callback := func() (interface{}, error) {
			senderResult := sender.Result()
			gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
			_ = bdw.Scale(1)
			drainerCompleted := test.WaitForDrainerRemoval(sw, 1)
			gomega.Expect(drainerCompleted).To(gomega.BeTrue())
			return drainerCompleted, nil
		}
		_, err = test.SendReceiveMessages(sender, receiver, callback)
		gomega.Expect(err).To(gomega.BeNil())

		receiverResult := receiver.Result()

		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}

	})

	// This test might fail due to ENTMQBR-3597
	ginkgo.It("Deploy 4 brokers, migrate everything to single", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		podNumbers := []string{"3", "2", "1"}

		err := bdw.DeployBrokers(4)
		gomega.Expect(err).To(gomega.BeNil())
		for _, number := range podNumbers {
			url := test.FormUrl(Protocol, DeployName, number, SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			sender = srw.WithSendUrl(url).PrepareNamedSender("sender-" + string(number))
			_ = sender.Deploy()
			sender.Wait()
		}
		err = bdw.Scale(1)
		gomega.Expect(err).To(gomega.BeNil())
		drainerCompleted := test.WaitForDrainerRemoval(sw, 3)
		gomega.Expect(drainerCompleted).To(gomega.BeTrue())
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		receiver = srw.
			WithReceiveUrl(receiveUrl).
			WithReceiverCount(len(podNumbers) * MessageCount).
			PrepareReceiver()

		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount * len(podNumbers)))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	// This test might fail due to ENTMQBR-3597
	ginkgo.It("Deploy 4 brokers, migrate last one", func() {
		sendUrl := test.FormUrl(Protocol, DeployName, "3", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
		sender, receiver = srw.
			WithReceiveUrl(receiveUrl).
			WithSendUrl(sendUrl).
			PrepareSenderReceiver()
		err := bdw.DeployBrokers(4)
		gomega.Expect(err).To(gomega.BeNil())
		_ = sender.Deploy()
		sender.Wait()
		_ = bdw.Scale(3)
		drainerCompleted := test.WaitForDrainerRemoval(sw, 1)
		gomega.Expect(drainerCompleted).To(gomega.BeTrue())
		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	// TODO: redesign mass migration test to be actually able to run it with giant message sizes and message quantities
})
