package test

import (
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp/qeclients"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

type SenderReceiverWrapper struct {
	ctx1          *framework.ContextData
	messageBody   string
	messageCount  int
	receiverCount int
	sendUrl       string
	receiveUrl    string
}

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
	return srw.PrepareSenderReceiverWithProtocol("amqp")
}

func (srw *SenderReceiverWrapper) PrepareSenderReceiverWithProtocol(protocol string) (*qeclients.AmqpQEClientCommon, *qeclients.AmqpQEClientCommon) {
	sender := srw.PrepareNamedSenderWithProtocol("sender", protocol)
	receiver := srw.PrepareReceiverWithProtocol(protocol)
	return sender, receiver
}

func (srw *SenderReceiverWrapper) PrepareNamedSender(name string) *qeclients.AmqpQEClientCommon {
	return srw.PrepareNamedSenderWithProtocol(name, "amqp")
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

	if protocol == "amqp" {
		senderBuilder.WithCustomCommand("cli-qpid-sender")
	} else if protocol == "core" {
		senderBuilder.WithCustomCommand("cli-artemis-sender")
	} else if protocol == "openwire" {
		senderBuilder.WithCustomCommand("cli-activemq-sender")
	}

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
	if protocol == "amqp" {
		receiverBuilder.WithCustomCommand("cli-qpid-receiver")
	} else if protocol == "core" {
		receiverBuilder.WithCustomCommand("cli-artemis-receiver")
	} else if protocol == "openwire" {
		receiverBuilder.WithCustomCommand("cli-activemq-receiver")
	}
	receiver, err := receiverBuilder.Build()
	if err != nil {
		panic(err)
	}
	return receiver
}
func (srw *SenderReceiverWrapper) PrepareReceiver() *qeclients.AmqpQEClientCommon {
	return srw.PrepareReceiverWithProtocol("amqp")
}
