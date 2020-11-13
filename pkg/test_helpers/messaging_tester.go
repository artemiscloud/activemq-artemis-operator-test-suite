package test_helpers

import (
	"errors"
	"fmt"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/pkg/bdw"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

func TestBaseSendReceiveMessages(bdw *bdw.BrokerDeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	acceptorType bdw.AcceptorType, BrokerCount int, protocol string) error {
	return TestBaseSendReceiveMessagesWithCallback(bdw, srw, MessageCount, MessageBody, acceptorType, BrokerCount, protocol, nil)
}

func TestBaseSendReceiveMessagesWithCallback(bdw *bdw.BrokerDeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	acceptorType bdw.AcceptorType, BrokerCount int, protocol string, callback test.SenderReceiverCallback) error {
	err := bdw.DeployBrokersWithAcceptor(BrokerCount, acceptorType)
	if err != nil {
		return err
	}
	sender, receiver := srw.PrepareSenderReceiverWithProtocol(protocol)
	_, err = test.SendReceiveMessages(sender, receiver, callback)
	if err != nil {
		return nil
	}
	senderResult := sender.Result()
	receiverResult := receiver.Result()

	if senderResult.Delivered != receiverResult.Delivered {
		return errors.New(fmt.Sprintf("Sent count is not equal to delivered count.\nSent: %d\n Delivered:%d\n Expected:%d\n",
			senderResult.Delivered, receiverResult.Delivered, MessageCount))
	}
	if senderResult.Delivered != MessageCount {
		return errors.New(fmt.Sprintf("MessageCount(%d) is not equal to sent/delivered (%d) count", MessageCount, senderResult.Delivered))
	}
	log.Logf("MessageCount is fine")
	for _, msg := range receiverResult.Messages {
		if msg.Content != MessageBody {
			return errors.New(fmt.Sprintf("Message content corrupted, expected: %s, real: %s", MessageBody, msg.Content))
		}
	}
	return nil
}

func TestBaseSendMessages(bdw *bdw.BrokerDeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	acceptorType bdw.AcceptorType, BrokerCount int, protocol, senderName string, callback test.SenderReceiverCallback) error {
	err := bdw.DeployBrokersWithAcceptor(BrokerCount, acceptorType)
	if err != nil {
		return err
	}
	srw.WithMessageCount(MessageCount).WithMessageBody(MessageBody)
	sender := srw.PrepareNamedSenderWithProtocol(senderName, protocol)
	_, err = test.SendMessages(sender, callback)
	if err != nil {
		return err
	}
	senderResult := sender.Result()
	if senderResult.Delivered != MessageCount {
		return errors.New(fmt.Sprintf("Sender Delivered count (%d) not equal to MessageCount (%d)", senderResult.Delivered, MessageCount))
	}
	return nil
}

func TestBaseReceiveMessages(bdw *bdw.BrokerDeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	protocol string) {
	receiver := srw.PrepareReceiverWithProtocol(protocol)
	err := test.ReceiveMessages(receiver)
	gomega.Expect(err).To(gomega.BeNil())
	receiverResult := receiver.Result()
	gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))
	log.Logf("MessageCount is fine")
	for _, msg := range receiverResult.Messages {
		gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
	}
}
