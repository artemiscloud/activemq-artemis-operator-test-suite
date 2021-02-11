package persistence

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"strconv"
)

var _ = ginkgo.Describe("MessagingPersistenceTests", func() {

	var (
		ctx1 *framework.ContextData
		//brokerClient brokerclientset.Interface
		brokerDeployer *bdw.BrokerDeploymentWrapper
		//url      string
		srw *test.SenderReceiverWrapper
	)

	var (
		MessageBody   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		MessageCount  = 100
		Port          = int64(bdw.AcceptorPorts[bdw.AmqpAcceptor])
		Domain        = "svc.cluster.local"
		SubdomainName = "-hdls-svc"
		AddressBit    = "someQueue"
		Protocol      = test.AMQP
	)

	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
		log.Logf("Value is: %v", test.Config.NeedsLatestCR)
		setEnv(ctx1, brokerDeployer)

		sendUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		receiveUrl := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, strconv.FormatInt(Port, 10))
		srw = &test.SenderReceiverWrapper{}
		srw.WithContext(ctx1).
			WithMessageBody(MessageBody).
			WithMessageCount(MessageCount).
			WithSendUrl(sendUrl).
			WithReceiveUrl(receiveUrl)
	})

	ginkgo.It("Deploy double instances with migration disabled, send messages, receive", func() {
		brokerDeployer.WithPersistence(true).WithMigration(false)
		test_helpers.TestBaseSendReceiveMessages(brokerDeployer, srw, MessageCount, MessageBody, bdw.AmqpAcceptor, 2, Protocol)
	})

	ginkgo.It("Deploy double instances with migration disabled, send messages, scaledown, scaleup, receive", func() {
		brokerDeployer.WithPersistence(true).WithMigration(false)
		callback := func() (interface{}, error) {
			err := brokerDeployer.Scale(1)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			err = brokerDeployer.Scale(2)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			return nil, nil
		}
		test_helpers.TestBaseSendReceiveMessagesWithCallback(brokerDeployer, srw, MessageCount, MessageBody, bdw.AmqpAcceptor, 2, Protocol, callback)
	})
})
