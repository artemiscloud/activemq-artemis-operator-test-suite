package test

import (
	"github.com/rh-messaging/shipshape/pkg/api/client/amqp/qeclients"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

type SenderReceiverWrapper struct {
	ctx1         *framework.ContextData
	messageBody  string
	messageCount int
	sendUrl      string
	receiveUrl   string
}

func (srw SenderReceiverWrapper) WithMessageBody(body string) SenderReceiverWrapper {
	srw.messageBody = body
	return srw
}

func (srw SenderReceiverWrapper) WithMessageCount(count int) SenderReceiverWrapper {
	srw.messageCount = count
	return srw
}

func (srw SenderReceiverWrapper) WithSendUrl(url string) SenderReceiverWrapper {
	srw.sendUrl = url
	return srw
}

func (srw SenderReceiverWrapper) WithReceiveUrl(url string) SenderReceiverWrapper {
	srw.receiveUrl = url
	return srw
}

func (srw SenderReceiverWrapper) WithContext(ctx1 *framework.ContextData) SenderReceiverWrapper {
	srw.ctx1 = ctx1
	return srw
}

func (srw SenderReceiverWrapper) PrepareSenderReceiver() (*qeclients.AmqpQEClientCommon, *qeclients.AmqpQEClientCommon) {
	sender := srw.PrepareSender()
	receiver := srw.PrepareReceiver()
	return sender, receiver
}

func (srw SenderReceiverWrapper) PrepareSender() *qeclients.AmqpQEClientCommon {
	sender, err := qeclients.NewSenderBuilder("sender",
		qeclients.Java,
		*srw.ctx1,
		srw.sendUrl).
		Content(srw.messageBody).
		Count(srw.messageCount).
		Timeout(20).
		Build()
	if err != nil {
		panic(err)
	}
	return sender
}

func (srw SenderReceiverWrapper) PrepareReceiver() *qeclients.AmqpQEClientCommon {
	sender, err := qeclients.
		NewReceiverBuilder("sender", qeclients.Java, *srw.ctx1, srw.sendUrl).
		Timeout(20).
		Build()
	if err != nil {
		panic(err)
	}
	return sender
}
