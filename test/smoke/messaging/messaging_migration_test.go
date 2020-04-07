package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/events"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	v1 "k8s.io/api/core/v1"
	"strings"
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
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = "5672"
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
	)

	// PrepareNamespace after framework has been created.
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

		sendUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, Port)
		receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, Port)

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
		sendUrls := []string{"3", "2", "1", "0"}
		receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, Port)
		srw.
			WithReceiveUrl(receiveUrl)
		receiver = srw.PrepareReceiver()

		err := dw.DeployBrokers(4)
		for _, url := range sendUrls {
			sender = srw.WithSendUrl(formUrl(url, SubdomainName, ctx1.Namespace, Domain, Port)).PrepareSender()
			_ = sender.Deploy()
			sender.Wait()
		}
		//Scale to 1
		err = dw.Scale(1)
		gomega.Expect(err).To(gomega.BeNil())

		err = framework.WaitForDeployment(ctx1.Clients.KubeClient, ctx1.Namespace, "drainer", 1, time.Second*10, time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
		podRemoved := false
		Framework.EventHandler.AddEventHandler(events.Pod, events.Delete, func(obj ...interface{}) {
			podObj := obj[0].(v1.Pod)
			if strings.Contains(podObj.Name, "drainer") {
				podRemoved = true
				log.Logf("Pod %s has been removed", podObj.Name)
			}
		})

		for !podRemoved {
			i := 0
			log.Logf("Still not finished...")
			time.Sleep(time.Second * 5)
			i++
			if i > 60 {
				break
			}
		}

		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount * len(sendUrls)))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate last one ", func() {
		sendUrl := formUrl("3", SubdomainName, ctx1.Namespace, Domain, Port)
		receiveUrl := formUrl("0", SubdomainName, ctx1.Namespace, Domain, Port)
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
		// Kludge to wait for removal of the drainer pod
		err = framework.WaitForDeployment(ctx1.Clients.KubeClient, ctx1.Namespace, "drainer", 0, time.Second*10, time.Minute*5)
		_ = receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
		for _, msg := range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

})
