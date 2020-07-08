package test

import (
	"fmt"
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp/qeclients"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
)

type SenderReceiverWrapper struct {
	ctx1          *framework.ContextData
	messageBody   string
	messageCount  int
	receiverCount int
	sendUrl       string
	receiveUrl    string
}

const AMQP = "amqp"
const OPENWIRE = "openwire"
const CORE = "core"

var (
	clientProtocolMap = map[string]string{
		AMQP:     "cli-qpid-%s",
		CORE:     "cli-artemis-%s",
		OPENWIRE: "cli-activemq-%s",
	}
)

type SenderReceiverCallback func() (interface{}, error)

func (srw *SenderReceiverWrapper) WithMessageBody(body string) *SenderReceiverWrapper {
	srw.messageBody = body
	return srw
}

func (srw *SenderReceiverWrapper) WithMessageCount(count int) *SenderReceiverWrapper {
	srw.messageCount = count
	return srw
}

func (srw *SenderReceiverWrapper) WithSendUrl(url string) *SenderReceiverWrapper {
	srw.sendUrl = url
	return srw
}

func (srw *SenderReceiverWrapper) WithReceiveUrl(url string) *SenderReceiverWrapper {
	srw.receiveUrl = url
	return srw
}

func (srw *SenderReceiverWrapper) WithContext(ctx1 *framework.ContextData) *SenderReceiverWrapper {
	srw.ctx1 = ctx1
	return srw
}

func (srw *SenderReceiverWrapper) PrepareSenderReceiver() (*qeclients.AmqpQEClientCommon, *qeclients.AmqpQEClientCommon) {
	return srw.PrepareSenderReceiverWithProtocol(AMQP)
}

func (srw *SenderReceiverWrapper) PrepareSenderReceiverWithProtocol(protocol string) (*qeclients.AmqpQEClientCommon, *qeclients.AmqpQEClientCommon) {
	sender := srw.PrepareNamedSenderWithProtocol("sender", protocol)
	receiver := srw.PrepareReceiverWithProtocol(protocol)
	return sender, receiver
}

func (srw *SenderReceiverWrapper) PrepareNamedSender(name string) *qeclients.AmqpQEClientCommon {
	return srw.PrepareNamedSenderWithProtocol(name, AMQP)
}

func (srw *SenderReceiverWrapper) PrepareNamedSenderWithProtocol(name string, protocol string) *qeclients.AmqpQEClientCommon {
	clientVer := qeclients.Java
	if Config.IBMz {
		clientVer = qeclients.JavaIBMZ
	}
	senderBuilder := qeclients.NewSenderBuilder(name,
		clientVer,
		*srw.ctx1,
		srw.sendUrl).
		Content(srw.messageBody).
		Count(srw.messageCount).
		Timeout(20)

	senderBuilder.WithCustomCommand(fmt.Sprintf(clientProtocolMap[protocol], "sender"))

	sender, err := senderBuilder.Build()
	if err != nil {
		panic(err)
	}
	return sender
}

func (srw *SenderReceiverWrapper) PrepareSender() *qeclients.AmqpQEClientCommon {
	return srw.PrepareNamedSender("sender")
}

func (srw *SenderReceiverWrapper) WithReceiverCount(count int) *SenderReceiverWrapper {
	srw.receiverCount = count
	return srw
}

func (srw *SenderReceiverWrapper) PrepareReceiverWithProtocol(protocol string) *qeclients.AmqpQEClientCommon {
	clientVer := qeclients.Java
	if Config.IBMz {
		clientVer = qeclients.JavaIBMZ
	}
	if srw.receiverCount == 0 {
		srw.receiverCount = srw.messageCount
	}
	receiverBuilder := qeclients.
		NewReceiverBuilder("receiver", clientVer, *srw.ctx1, srw.receiveUrl).
		Timeout(20).
		WithCount(srw.receiverCount)
	receiverBuilder.WithCustomCommand(fmt.Sprintf(clientProtocolMap[protocol], "receiver"))

	receiver, err := receiverBuilder.Build()
	if err != nil {
		panic(err)
	}
	return receiver
}
func (srw *SenderReceiverWrapper) PrepareReceiver() *qeclients.AmqpQEClientCommon {
	return srw.PrepareReceiverWithProtocol(AMQP)
}

func SendReceiveMessages(sender *qeclients.AmqpQEClientCommon, receiver *qeclients.AmqpQEClientCommon, callback SenderReceiverCallback) (interface{}, error) {
	log.Logf("Started (sync) deployment of sender")
	err := sender.Deploy()
	sender.Wait()
	if err != nil {
		return nil, err
	}
	var result interface{}
	if callback != nil {
		result, err = callback()
		if err != nil {
			return nil, err
		}
	}
	log.Logf("Started (sync) deployment of receiver")
	err = receiver.Deploy()
	if err != nil {
		return nil, err
	}
	receiver.Wait()
	return result, nil
}
