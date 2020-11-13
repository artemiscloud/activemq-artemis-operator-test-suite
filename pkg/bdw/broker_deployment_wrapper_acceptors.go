package bdw
/* This file contains acceptor-related helper methods for BrokerDeploymentWrapper
 */

import (
	"fmt"
	brokerv3 "github.com/artemiscloud/activemq-artemis-operator/pkg/apis/broker/v2alpha3"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

type AcceptorType int

const (
	AmqpAcceptor AcceptorType = iota
	CoreAcceptor
	OpenwireAcceptor
	MultiAcceptor
	AllAcceptor
	NoChangeAcceptor
)

var (
	AcceptorPorts = map[AcceptorType]int32{
		AmqpAcceptor:     5672,
		OpenwireAcceptor: 61613,
		CoreAcceptor:     61616,
		MultiAcceptor:    5672,
		AllAcceptor:      5672,
	}
	// Interface to ease transition
	acceptors = map[AcceptorType]*brokerv3.AcceptorType{
		AmqpAcceptor:     defaultAcceptor(test.AMQP, AcceptorPorts[AmqpAcceptor]),
		OpenwireAcceptor: defaultAcceptor(test.OPENWIRE, AcceptorPorts[OpenwireAcceptor]),
		CoreAcceptor:     defaultAcceptor(test.CORE, AcceptorPorts[CoreAcceptor]),
		MultiAcceptor:    defaultAcceptor(fmt.Sprintf("%s,%s,%s", test.AMQP, test.OPENWIRE, test.CORE), AcceptorPorts[MultiAcceptor]),
		AllAcceptor:      defaultAcceptor("all", AcceptorPorts[AllAcceptor]),
	}
)


func getAcceptor(name string, port int32, protocols string, sslEnabled bool, sslSecret string, enabledCipherSuites string,
	enabledProtocols string, needsClientAuth bool, wantClientAuth bool, verifyHost bool, sslProvider string, sniHost string,
	expose bool, anycastPrefix string, multicastPrefix string, connectionsAllowed int) *brokerv3.AcceptorType {
	acceptor := &brokerv3.AcceptorType{
		Name:                name,
		Port:                port,
		Protocols:           protocols,
		SSLEnabled:          sslEnabled,
		SSLSecret:           sslSecret,
		EnabledCipherSuites: enabledCipherSuites,
		EnabledProtocols:    enabledProtocols,
		NeedClientAuth:      needsClientAuth,
		WantClientAuth:      wantClientAuth,
		VerifyHost:          verifyHost,
		SSLProvider:         sslProvider,
		SNIHost:             sniHost,
		Expose:              expose,
		AnycastPrefix:       anycastPrefix,
		MulticastPrefix:     multicastPrefix,
		ConnectionsAllowed:  connectionsAllowed,
	}
	return acceptor
}

func defaultAcceptor(protocol string, port int32) *brokerv3.AcceptorType {
	return getAcceptor(protocol,
		port,
		protocol,
		false,
		"",
		"",
		"",
		false,
		false,
		false,
		"JDK",
		"localhost",
		true,
		"",
		"",
		0)
}
