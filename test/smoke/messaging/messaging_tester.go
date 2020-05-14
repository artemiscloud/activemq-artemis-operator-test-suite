package messaging

import (
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"math/rand"
)

func testBaseSendReceiveSingleBroker(dw *test.DeploymentWrapper,
	srw *test.SenderReceiverWrapper,
	MessageCount int, MessageBody string,
	acceptorType test.AcceptorType, BrokerCount int, protocol string) {
	err := dw.DeployBrokersWithAcceptor(BrokerCount, acceptorType)
	gomega.Expect(err).To(gomega.BeNil())
	sender, receiver := srw.PrepareSenderReceiverWithProtocol(protocol)
	_ = sender.Deploy()
	_ = receiver.Deploy()
	log.Logf("Started (sync) deployment of clients")
	sender.Wait()
	receiver.Wait()
	log.Logf("Wait finished")
	senderResult := sender.Result()
	receiverResult := receiver.Result()
	log.Logf("Finished (sync) deployment")
	log.Logf("Count sent: %d", senderResult.Delivered)
	log.Logf("Count received: %d", receiverResult.Delivered)
	log.Logf("Len of received: %d", len(receiverResult.Messages))
	gomega.Expect(senderResult.Delivered).To(gomega.Equal(MessageCount))
	gomega.Expect(receiverResult.Delivered).To(gomega.Equal(MessageCount))

	log.Logf("MessageCount is fine")
	for _, msg := range receiverResult.Messages {
		gomega.Expect(msg.Content).To(gomega.Equal(MessageBody))
	}
}

func testSizedMessage(srw *test.SenderReceiverWrapper, dw *test.DeploymentWrapper, receiver amqp.Client, size, count int) {
	err := dw.DeployBrokers(1)
	gomega.Expect(err).To(gomega.BeNil())
	body := RandStringBytes(size)
	srw.WithMessageCount(count).
		WithMessageBody(body)
	sender := srw.PrepareSender()
	_ = sender.Deploy()
	_ = receiver.Deploy()
	sender.Wait()
	receiver.Wait()

	senderResult := sender.Result()
	receiverResult := receiver.Result()

	gomega.Expect(senderResult.Delivered).To(gomega.Equal(count))
	gomega.Expect(receiverResult.Delivered).To(gomega.Equal(count))

	for _, msg := range receiverResult.Messages {
		gomega.Expect(msg.Content).To(gomega.Equal(body))
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
