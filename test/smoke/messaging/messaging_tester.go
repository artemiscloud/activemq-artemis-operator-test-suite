package messaging

import (
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

func testBaseSendReceiveMessages(dw *test.DeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	acceptorType test.AcceptorType, BrokerCount int, protocol string) {
	testBaseSendReceiveMessagesWithCallback(dw, srw, MessageCount, MessageBody, acceptorType, BrokerCount, protocol, nil)
}

func testBaseSendReceiveMessagesWithCallback(dw *test.DeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	acceptorType test.AcceptorType, BrokerCount int, protocol string, callback test.SenderReceiverCallback) {
	err := dw.DeployBrokersWithAcceptor(BrokerCount, acceptorType)
	gomega.Expect(err).To(gomega.BeNil())
	sender, receiver := srw.PrepareSenderReceiverWithProtocol(protocol)
	_, err = test.SendReceiveMessages(sender, receiver, callback)
	gomega.Expect(err).To(gomega.BeNil())
	senderResult := sender.Result()
	receiverResult := receiver.Result()
	gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
	gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))

	log.Logf("MessageCount is fine")
	for _, msg := range receiverResult.Messages {
		gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
	}
}
