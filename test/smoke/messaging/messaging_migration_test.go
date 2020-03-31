package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp/qeclients"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"time"
)

var _ = ginkgo.Describe("MessagingMigrationTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		dw test.DeploymentWrapper
		sender amqp.Client
		receiver amqp.Client
		url string
		err error
	)

	const (
		MessageBody = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount = 100
	)

	// PrepareNamespace after framework has been created. Doesn'
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).
			WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName).
			WithPersistence(true).WithMigration(true)
		url = "amqp://ex-aao-ss-1:5672/"
		sender, err = qeclients.NewSenderBuilder("sender", qeclients.Python, *ctx1, url).Content(MessageBody).Count(MessageCount).Build()//, MessageBody, MessageCount)
		if err!=nil {
			panic(err)
		}
		url = "amqp://ex-aao-ss-0:5672/"
		receiver, err = qeclients.NewReceiverBuilder("receiver", qeclients.Python, *ctx1, url).Build()
		if err!=nil {
			panic(err)
		}

	})

	ginkgo.It("Deploy double broker instance, migrate to single", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers( 2)
		gomega.Expect(err).To(gomega.BeNil())

		sender.Deploy()
		sender.Wait()
		senderResult := sender.Result()
		gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
		dw.Scale(1)

		//Wait for a drainer pod to do its deed
		err = framework.WaitForDeployment(ctx1.Clients.KubeClient,ctx1.Namespace, "drainer", 1, time.Second * 10, time.Minute * 5)
		gomega.Expect(err).To(gomega.BeNil())
		//wait for Drainer to disappear.
		list,_ := ctx1.ListPodsForDeploymentName("drainer")
		//framework.WaitForDeletion(ctx1.Clients.KubeClient, list, time.Second*10, time.Minute*10)
		// err = framework.WaitForDeletion(ctx1.Clients.DynClient, DrainerObject, time.Second * 10, time.Minute * 10)

		receiver.Deploy()
		receiver.Wait()
		receiverResult := receiver.Result()
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
		for _, msg:= range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy 4 brokers, migrate everything to single", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers(4)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("Deploy 4 brokers, migrate last one ", func () {
		err := dw.DeployBrokers(4)
		gomega.Expect(err).To(gomega.BeNil())
	})

})
