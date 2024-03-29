package persistence

import (
	"strconv"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/test_helpers"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

var _ = ginkgo.Describe("PersistenceVolumeSizeTest", func() {

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
	ginkgo.BeforeEach(func() {
	})
	// PrepareNamespace after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = sw.Framework.GetFirstContext()
		brokerDeployer = &bdw.BrokerDeploymentWrapper{}
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

	ginkgo.It("Deploy with smaller PVC", func() {
		brokerDeployer.WithPersistence(true).WithMigration(false).WithStorageSize("1Gi")
		test_helpers.CheckVolumeSize(ctx1, "1Gi")
	})

	ginkgo.It("Deploy with smallest PVC", func() {
		brokerDeployer.WithPersistence(true).WithMigration(false).WithStorageSize("1")
		test_helpers.CheckVolumeSize(ctx1, "1Gi")
	})

	ginkgo.It("Deploy with default PVC", func() {
		brokerDeployer.WithPersistence(true).WithMigration(false)
		test_helpers.CheckVolumeSize(ctx1, "2Gi")
	})

	ginkgo.It("Deploy with bigger PVC", func() {
		brokerDeployer.WithPersistence(true).WithMigration(false).WithStorageSize("3Gi")
		test_helpers.CheckVolumeSize(ctx1, "3Gi")
	})

})
