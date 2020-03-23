package messaging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp/qeclients"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

var _ = ginkgo.Describe("MessagingBasicTests", func() {

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

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
		dw = test.DeploymentWrapper{}.WithWait(true).WithBrokerClient(brokerClient).WithContext(ctx1).WithCustomImage(test.TestConfig.BrokerImageName)
		url = "amqp://ex-aao-ss-0:5672/"
		sender, err = qeclients.NewSenderBuilder("sender", qeclients.Python, *ctx1, url).Content(MessageBody).Count(MessageCount).Build()//, MessageBody, MessageCount)
		if err!=nil {
			panic(err)
		}
		receiver, err = qeclients.NewReceiverBuilder("receiver", qeclients.Python, *ctx1, url).Build()
		if err!=nil {
			panic(err)
		}

	})

	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers( 1)
		gomega.Expect(err).To(gomega.BeNil())

		sender.Deploy()
		receiver.Deploy()
		sender.Wait()
		receiver.Wait()
		senderResult := sender.Result()
		receiverResult := receiver.Result()
		gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
		gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
		for _, msg:= range receiverResult.Messages {
			gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
		}
	})

	ginkgo.It("Deploy double broker instances", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		err := dw.DeployBrokers(2)
		gomega.Expect(err).To(gomega.BeNil())
	})

})
