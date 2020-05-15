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
		ctx1     *framework.ContextData
		dw       *test.DeploymentWrapper
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
		Protocol      = "amqp"
	)

	// PrepareNamespace after framework has been created.
	ginkgo.JustBeforeEach(func() {
		if !test.Config.IBMz {
			ctx1 = Framework.GetFirstContext()
			dw = &test.DeploymentWrapper{}
			dw.WithWait(true).WithBrokerClient(brokerClient).
				WithContext(ctx1).WithCustomImage(test.Config.BrokerImageName).
				WithPersistence(true).WithMigration(true).
				WithName(DeployName)
			srw = &test.SenderReceiverWrapper{}
			srw.WithContext(ctx1).
				WithMessageBody(MessageBody).
				WithMessageCount(MessageCount)

		}
	})

	ginkgo.JustAfterEach(func() {
		Framework.GetFirstContext().EventHandler.ClearCallbacks()
	})

	ginkgo.It("Deploy double broker instance, migrate to single", func() {
		if !test.Config.IBMz {
			err := dw.DeployBrokers(2)
			gomega.Expect(err).To(gomega.BeNil())

			sendUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			receiveUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)

			sender, receiver = srw.
				WithReceiveUrl(receiveUrl).
				WithSendUrl(sendUrl).
				PrepareSenderReceiver()
			_ = sender.Deploy()
			sender.Wait()

			senderResult := sender.Result()
			gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
			_ = dw.Scale(1)

			drainerCompleted := WaitForDrainerRemoval(1)
			gomega.Expect(drainerCompleted).To(gomega.BeTrue())
			gomega.Expect(err).To(gomega.BeNil())
			_ = receiver.Deploy()
			receiver.Wait()
			receiverResult := receiver.Result()
			gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
			for _, msg := range receiverResult.Messages {
				gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
			}
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate everything to single", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		if !test.Config.IBMz {
			podNumbers := []string{"3", "2", "1"}

			err := dw.DeployBrokers(4)
			gomega.Expect(err).To(gomega.BeNil())
			for _, number := range podNumbers {
				url := formUrl(Protocol, number, SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
				sender = srw.WithSendUrl(url).PrepareNamedSender("sender-" + string(number))
				_ = sender.Deploy()
				sender.Wait()
			}
			err = dw.Scale(1)
			gomega.Expect(err).To(gomega.BeNil())
			drainerCompleted := WaitForDrainerRemoval(3)
			gomega.Expect(drainerCompleted).To(gomega.BeTrue())
			receiveUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
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
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate last one ", func() {
		if !test.Config.IBMz {
			sendUrl := formUrl(Protocol, "3", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			receiveUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			sender, receiver = srw.
				WithReceiveUrl(receiveUrl).
				WithSendUrl(sendUrl).
				PrepareSenderReceiver()
			err := dw.DeployBrokers(4)
			gomega.Expect(err).To(gomega.BeNil())
			_ = sender.Deploy()
			sender.Wait()
			_ = dw.Scale(3)
			drainerCompleted := WaitForDrainerRemoval(1)
			gomega.Expect(drainerCompleted).To(gomega.BeTrue())
			_ = receiver.Deploy()
			receiver.Wait()
			receiverResult := receiver.Result()
			gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
			for _, msg := range receiverResult.Messages {
				gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
			}
		}
	})

	ginkgo.It("Mass migration of messages", func() {
		if !test.Config.IBMz {
			sendUrl := formUrl(Protocol, "1", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			receiveUrl := formUrl(Protocol, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			BigMultiplier := 10000
			srw.WithMessageCount(BigMultiplier * MessageCount)
			sender, receiver = srw.
				WithReceiveUrl(receiveUrl).
				WithSendUrl(sendUrl).
				PrepareSenderReceiver()
			err := dw.DeployBrokers(2)
			gomega.Expect(err).To(gomega.BeNil())
			_ = sender.Deploy()
			sender.Wait()
			_ = dw.Scale(1)
			drainerCompleted := WaitForDrainerRemovalSlow(1, time.Second*time.Duration(10), 1000)
			gomega.Expect(drainerCompleted).To(gomega.BeTrue())
			_ = receiver.Deploy()
			receiver.Wait()
			receiverResult := receiver.Result()
			gomega.Expect(receiverResult.Delivered).To(gomega.Equal(BigMultiplier * MessageCount))
			for _, msg := range receiverResult.Messages {
				gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
			}
		}
	})
})
