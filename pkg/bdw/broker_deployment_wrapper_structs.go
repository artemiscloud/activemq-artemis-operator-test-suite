package bdw

/* This file contains structs for BrokerDeploymentWrapper
 */

import (
	brokerbeta "github.com/artemiscloud/activemq-artemis-operator/api/v1beta1"
	brokerclientset "github.com/artemiscloud/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
)

// BrokerDeploymentWrapper takes care of deployment of Broker
type BrokerDeploymentWrapper struct {
	wait            bool
	brokerClient    brokerclientset.Interface
	ctx1            *framework.ContextData
	customImage     string
	customInit      string
	migration       bool
	persistence     bool
	name            string
	sslEnabled      bool
	exposeConsole   bool
	deploymentSize  int
	isLtsDeployment bool
	storageSize     string
	timeoutMult     int
	securities      map[string]brokerbeta.ActiveMQArtemisSecurity
	AddressSettings
	ResourcesLimits
	ResourcesRequests
}
type ResourcesLimits struct {
	cpu string
	mem string
}

type ResourcesRequests struct {
	cpu string
	mem string
}

type AddressSettings struct {
	knownAddresses                     []string
	maxSizeBytes                       map[string]string
	addressFullPolicy                  map[string]AddressFullPolicy
	deadLetterAddress                  map[string]string
	autoCreateDeadLetterResources      map[string]bool
	dlqPrefix                          map[string]string
	dlqSuffix                          map[string]string
	expiryAddress                      map[string]string
	autoCreateExpiryResources          map[string]bool
	expirySuffix                       map[string]string
	expiryPrefix                       map[string]string
	expiryDelay                        map[string]int32
	minExpiryDelay                     map[string]int32
	maxExpiryDelay                     map[string]int32
	redeliveryDelay                    map[string]int32
	maxRedeliveryDelay                 map[string]int32
	redeliveryDelayMult                map[string]string
	redeliveryCollisionsAvoidance      map[string]string
	maxRedeliveryAttempts              map[string]int32
	maxSizeBytesRejectThreshold        map[string]int32
	pageSizeBytes                      map[string]string
	pageMaxCacheSize                   map[string]int32
	messageCounterHistoryDayLimit      map[string]int32
	lastValueQueue                     map[string]bool
	defaultLastValueQueue              map[string]bool
	defaultLastValueKey                map[string]string
	defaultNonDestructive              map[string]bool
	defaultExclusiveQueue              map[string]bool
	defaultGroupRebalance              map[string]bool
	defaultGroupRebalancePauseDispatch map[string]bool
	defaultGroupBuckets                map[string]int32
	defaultGroupFirstKey               map[string]string
	defaultConsumerBeforeDispatch      map[string]int32
	defaultDelayBeforeDispatch         map[string]int32
	redistributionDelay                map[string]int32
	sendToDLAOnNoRoute                 map[string]bool
	slowConsumerThreshold              map[string]int32
	slowConsumerPolicy                 map[string]SlowConsumerPolicy
	slowConsumerCheckPeriod            map[string]int32
	autoCreateJmsQueues                map[string]bool
	autoDeleteJmsQueues                map[string]bool
	autoCreateJmsTopics                map[string]bool
	autoDeleteJmsTopics                map[string]bool
	autoCreateQueues                   map[string]bool
	autoDeleteQueues                   map[string]bool
	autoDeleteCreatedQueues            map[string]bool
	autoDeleteQueuesDelay              map[string]int32
	audoDeleteQueuesMessageCount       map[string]int32
	configDeleteQueues                 map[string]ConfigDelete
	autoCreateAddresses                map[string]bool
	autoDeleteAddressesDelay           map[string]int32
	configDeleteAddresses              map[string]ConfigDelete
	managementBrowsePageSize           map[string]int32
	defaultPurgeOnNoConsumers          map[string]bool
	defaultMaxConsumers                map[string]int32
	defaultQueueRoutingType            map[string]RoutingType
	defaultAddressRoutingType          map[string]RoutingType
	defaultConsumerWindowSize          map[string]int32
	defaultRingSize                    map[string]int32
	defaultRetroMessageCount           map[string]int32
	enableMetrics                      map[string]bool
	autoDeleteAddresses                map[string]bool
}

type RoutingType int
type ConfigDelete int
type SlowConsumerPolicy int //TODO
type AddressFullPolicy int

const (
	DROP      = "DROP"
	FAIL      = "FAIL"
	PAGE      = "PAGE"
	BLOCK     = "BLOCK"
	OFF       = "OFF"
	FORCE     = "FORCE"
	NOTIFY    = "NOTIFY"
	KILL      = "KILL"
	ANYCAST   = "ANYCAST"
	MULTICAST = "MULTICAST"

	DropPolicy AddressFullPolicy = iota
	FailPolicy
	PagePolicy
	BlockPolicy

	Off ConfigDelete = iota
	Force

	Kill SlowConsumerPolicy = iota
	Notify

	Anycast RoutingType = iota
	Multicast
)

func (a ConfigDelete) String() string {
	if a == Force {
		return FORCE
	}
	return OFF
}
func (a SlowConsumerPolicy) String() string {
	if a == Kill {
		return KILL
	}
	return NOTIFY
}

func (a RoutingType) String() string {
	if a == Anycast {
		return ANYCAST
	}
	return MULTICAST
}

func (a AddressFullPolicy) String() string {
	switch a {
	case DropPolicy:
		return DROP
	case FailPolicy:
		return FAIL
	case PagePolicy:
		return PAGE
	case BlockPolicy:
		return BLOCK
	}
	return PAGE
}
