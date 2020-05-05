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
		dw       test.DeploymentWrapper
		srw      test.SenderReceiverWrapper
		sender   amqp.Client
		receiver amqp.Client
	)

	const (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = "5672"
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
	)

	// PrepareNamespace after framework has been created.
	ginkgo.JustBeforeEach(func() {
		if !test.Config.IBMz {
			ctx1 = Framework.GetFirstContext()
			dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).
				WithContext(ctx1).WithCustomImage(test.Config.BrokerImageName).
				WithPersistence(true).WithMigration(true).
				WithName(DeployName)
			srw = test.SenderReceiverWrapper{}.WithContext(ctx1).
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

			sendUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)

			sender, receiver = srw.
				WithReceiveUrl(receiveUrl).
				WithSendUrl(sendUrl).
				PrepareSenderReceiver()
			_ = sender.Deploy()
			sender.Wait()

			senderResult := sender.Result()
			gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
			_ = dw.Scale(1)

			WaitForDrainerRemoval(1)

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
			sendUrls := []string{"3", "2", "1", "0"}
			receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			srw.
				WithReceiveUrl(receiveUrl)
			receiver = srw.PrepareReceiver()

			err := dw.DeployBrokers(4)
			for _, url := range sendUrls {
				sender = srw.WithSendUrl(formUrl(url, SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)).PrepareSender()
				_ = sender.Deploy()
				sender.Wait()
			}
			err = dw.Scale(1)
			gomega.Expect(err).To(gomega.BeNil())
			WaitForDrainerRemoval(3)
			err = framework.WaitForDeployment(ctx1.Clients.KubeClient, ctx1.Namespace, "drainer", 1, time.Second*10, time.Minute*5)
			gomega.Expect(err).To(gomega.BeNil())
			WaitForDrainerRemoval(3)
			_ = receiver.Deploy()
			receiver.Wait()
			receiverResult := receiver.Result()
			gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount * len(sendUrls)))
			for _, msg := range receiverResult.Messages {
				gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
			}
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate last one ", func() {
		if !test.Config.IBMz {
			sendUrl := formUrl("3", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port)
			sender, receiver = srw.
				WithReceiveUrl(receiveUrl).
				WithSendUrl(sendUrl).
				PrepareSenderReceiver()
			err := dw.DeployBrokers(4)
			gomega.Expect(err).To(gomega.BeNil())
			_ = sender.Deploy()
			sender.Wait()
			_ = dw.Scale(3)
			WaitForDrainerRemoval(1)
			receiver.Wait()
			receiverResult := receiver.Result()
			gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
			for _, msg := range receiverResult.Messages {
				gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
			}
		}
	})

})
