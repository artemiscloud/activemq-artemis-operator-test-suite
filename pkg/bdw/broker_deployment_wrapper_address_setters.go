package bdw

/* This file contains configuration setter methods for BrokerDeploymentWrapper
 */

func (bdw *BrokerDeploymentWrapper) WithMaxSizeBytes(addressName string, maxSizeBytes string) *BrokerDeploymentWrapper {
	if bdw.maxSizeBytes == nil {
		bdw.maxSizeBytes = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.maxSizeBytes[addressName] = maxSizeBytes
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAddressFullPolicy(addressName string, addressFullPolicy AddressFullPolicy) *BrokerDeploymentWrapper {
	if bdw.addressFullPolicy == nil {
		bdw.addressFullPolicy = map[string]AddressFullPolicy{}
	}
	bdw.addKnownAddress(addressName)
	bdw.addressFullPolicy[addressName] = addressFullPolicy
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDeadLetterAddress(addressName string, deadLetterAddress string) *BrokerDeploymentWrapper {
	if bdw.deadLetterAddress == nil {
		bdw.deadLetterAddress = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.deadLetterAddress[addressName] = deadLetterAddress
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoCreateDeadLetterResources(addressName string, autoCreateDeadLetterResources bool) *BrokerDeploymentWrapper {
	if bdw.autoCreateDeadLetterResources == nil {
		bdw.autoCreateDeadLetterResources = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoCreateDeadLetterResources[addressName] = autoCreateDeadLetterResources
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDlqPrefix(addressName string, dlqPrefix string) *BrokerDeploymentWrapper {
	if bdw.dlqPrefix == nil {
		bdw.dlqPrefix = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.dlqPrefix[addressName] = dlqPrefix
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDlqSuffix(addressName string, dlqSuffix string) *BrokerDeploymentWrapper {
	if bdw.dlqSuffix == nil {
		bdw.dlqSuffix = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.dlqSuffix[addressName] = dlqSuffix
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithExpiryAddress(addressName string, expiryAddress string) *BrokerDeploymentWrapper {
	if bdw.expiryAddress == nil {
		bdw.expiryAddress = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.expiryAddress[addressName] = expiryAddress
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoCreateExpiryResources(addressName string, autoCreateExpiryResources bool) *BrokerDeploymentWrapper {
	if bdw.autoCreateExpiryResources == nil {
		bdw.autoCreateExpiryResources = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoCreateExpiryResources[addressName] = autoCreateExpiryResources
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithExpirySuffix(addressName string, expirySuffix string) *BrokerDeploymentWrapper {
	if bdw.expirySuffix == nil {
		bdw.expirySuffix = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.expirySuffix[addressName] = expirySuffix
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithExpiryPrefix(addressName string, expiryPrefix string) *BrokerDeploymentWrapper {
	if bdw.expiryPrefix == nil {
		bdw.expiryPrefix = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.expiryPrefix[addressName] = expiryPrefix
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithExpiryDelay(addressName string, expiryDelay int32) *BrokerDeploymentWrapper {
	if bdw.expiryDelay == nil {
		bdw.expiryDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.expiryDelay[addressName] = expiryDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMinExpiryDelay(addressName string, minExpiryDelay int32) *BrokerDeploymentWrapper {
	if bdw.minExpiryDelay == nil {
		bdw.minExpiryDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.minExpiryDelay[addressName] = minExpiryDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMaxExpiryDelay(addressName string, maxExpiryDelay int32) *BrokerDeploymentWrapper {
	if bdw.maxExpiryDelay == nil {
		bdw.maxExpiryDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.maxExpiryDelay[addressName] = maxExpiryDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithRedeliveryDelay(addressName string, redeliveryDelay int32) *BrokerDeploymentWrapper {
	if bdw.redeliveryDelay == nil {
		bdw.redeliveryDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.redeliveryDelay[addressName] = redeliveryDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMaxRedeliveryDelay(addressName string, maxRedeliveryDelay int32) *BrokerDeploymentWrapper {
	if bdw.maxRedeliveryDelay == nil {
		bdw.maxRedeliveryDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.maxRedeliveryDelay[addressName] = maxRedeliveryDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithRedeliveryDelayMult(addressName string, redeliveryDelayMult int32) *BrokerDeploymentWrapper {
	if bdw.redeliveryDelayMult == nil {
		bdw.redeliveryDelayMult = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.redeliveryDelayMult[addressName] = redeliveryDelayMult
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithRedeliveryCollisionsAvoidance(addressName string, redeliveryCollisionsAvoidance int32) *BrokerDeploymentWrapper {
	if bdw.redeliveryCollisionsAvoidance == nil {
		bdw.redeliveryCollisionsAvoidance = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.redeliveryCollisionsAvoidance[addressName] = redeliveryCollisionsAvoidance
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMaxRedeliveryAttempts(addressName string, maxRedeliveryAttempts int32) *BrokerDeploymentWrapper {
	if bdw.maxRedeliveryAttempts == nil {
		bdw.maxRedeliveryAttempts = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.maxRedeliveryAttempts[addressName] = maxRedeliveryAttempts
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMaxSizeBytesRejectThreshold(addressName string, maxSizeBytesRejectThreshold int32) *BrokerDeploymentWrapper {
	if bdw.maxSizeBytesRejectThreshold == nil {
		bdw.maxSizeBytesRejectThreshold = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.maxSizeBytesRejectThreshold[addressName] = maxSizeBytesRejectThreshold
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithPageSizeBytes(addressName string, pageSizeBytes string) *BrokerDeploymentWrapper {
	if bdw.pageSizeBytes == nil {
		bdw.pageSizeBytes = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.pageSizeBytes[addressName] = pageSizeBytes
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithPageMaxCacheSize(addressName string, pageMaxCacheSize int32) *BrokerDeploymentWrapper {
	if bdw.pageMaxCacheSize == nil {
		bdw.pageMaxCacheSize = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.pageMaxCacheSize[addressName] = pageMaxCacheSize
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMessageCounterHistoryDayLimit(addressName string, messageCounterHistoryDayLimit int32) *BrokerDeploymentWrapper {
	if bdw.messageCounterHistoryDayLimit == nil {
		bdw.messageCounterHistoryDayLimit = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.messageCounterHistoryDayLimit[addressName] = messageCounterHistoryDayLimit
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithLastValueQueue(addressName string, lastValueQueue bool) *BrokerDeploymentWrapper {
	if bdw.lastValueQueue == nil {
		bdw.lastValueQueue = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.lastValueQueue[addressName] = lastValueQueue
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultLastValueQueue(addressName string, defaultLastValueQueue bool) *BrokerDeploymentWrapper {
	if bdw.defaultLastValueQueue == nil {
		bdw.defaultLastValueQueue = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultLastValueQueue[addressName] = defaultLastValueQueue
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultLastValueKey(addressName string, defaultLastValueKey string) *BrokerDeploymentWrapper {
	if bdw.defaultLastValueKey == nil {
		bdw.defaultLastValueKey = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultLastValueKey[addressName] = defaultLastValueKey
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultNonDestructive(addressName string, defaultNonDestructive bool) *BrokerDeploymentWrapper {
	if bdw.defaultNonDestructive == nil {
		bdw.defaultNonDestructive = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultNonDestructive[addressName] = defaultNonDestructive
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultExclusiveQueue(addressName string, defaultExclusiveQueue bool) *BrokerDeploymentWrapper {
	if bdw.defaultExclusiveQueue == nil {
		bdw.defaultExclusiveQueue = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultExclusiveQueue[addressName] = defaultExclusiveQueue
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultGroupRebalance(addressName string, defaultGroupRebalance bool) *BrokerDeploymentWrapper {
	if bdw.defaultGroupRebalance == nil {
		bdw.defaultGroupRebalance = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultGroupRebalance[addressName] = defaultGroupRebalance
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultGroupRebalancePauseDispatch(addressName string, defaultGroupRebalancePauseDispatch bool) *BrokerDeploymentWrapper {
	if bdw.defaultGroupRebalancePauseDispatch == nil {
		bdw.defaultGroupRebalancePauseDispatch = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultGroupRebalancePauseDispatch[addressName] = defaultGroupRebalancePauseDispatch
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultGroupBuckets(addressName string, defaultGroupBuckets int32) *BrokerDeploymentWrapper {
	if bdw.defaultGroupBuckets == nil {
		bdw.defaultGroupBuckets = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultGroupBuckets[addressName] = defaultGroupBuckets
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultGroupFirstKey(addressName string, defaultGroupFirstKey string) *BrokerDeploymentWrapper {
	if bdw.defaultGroupFirstKey == nil {
		bdw.defaultGroupFirstKey = map[string]string{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultGroupFirstKey[addressName] = defaultGroupFirstKey
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultConsumerBeforeDispatch(addressName string, defaultConsumerBeforeDispatch int32) *BrokerDeploymentWrapper {
	if bdw.defaultConsumerBeforeDispatch == nil {
		bdw.defaultConsumerBeforeDispatch = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultConsumerBeforeDispatch[addressName] = defaultConsumerBeforeDispatch
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultDelayBeforeDispatch(addressName string, defaultDelayBeforeDispatch int32) *BrokerDeploymentWrapper {
	if bdw.defaultDelayBeforeDispatch == nil {
		bdw.defaultDelayBeforeDispatch = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultDelayBeforeDispatch[addressName] = defaultDelayBeforeDispatch
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithRedistributionDelay(addressName string, redistributionDelay int32) *BrokerDeploymentWrapper {
	if bdw.redistributionDelay == nil {
		bdw.redistributionDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.redistributionDelay[addressName] = redistributionDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithSendToDLAOnNoRoute(addressName string, sendToDLAOnNoRoute bool) *BrokerDeploymentWrapper {
	if bdw.sendToDLAOnNoRoute == nil {
		bdw.sendToDLAOnNoRoute = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.sendToDLAOnNoRoute[addressName] = sendToDLAOnNoRoute
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithSlowConsumerThreshold(addressName string, slowConsumerThreshold int32) *BrokerDeploymentWrapper {
	if bdw.slowConsumerThreshold == nil {
		bdw.slowConsumerThreshold = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.slowConsumerThreshold[addressName] = slowConsumerThreshold
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithSlowConsumerPolicy(addressName string, slowConsumerPolicy SlowConsumerPolicy) *BrokerDeploymentWrapper {
	if bdw.slowConsumerPolicy == nil {
		bdw.slowConsumerPolicy = map[string]SlowConsumerPolicy{}
	}
	bdw.addKnownAddress(addressName)
	bdw.slowConsumerPolicy[addressName] = slowConsumerPolicy
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithSlowConsumerCheckPeriod(addressName string, slowConsumerCheckPeriod int32) *BrokerDeploymentWrapper {
	if bdw.slowConsumerCheckPeriod == nil {
		bdw.slowConsumerCheckPeriod = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.slowConsumerCheckPeriod[addressName] = slowConsumerCheckPeriod
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoCreateJmsQueues(addressName string, autoCreateJmsQueues bool) *BrokerDeploymentWrapper {
	if bdw.autoCreateJmsQueues == nil {
		bdw.autoCreateJmsQueues = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoCreateJmsQueues[addressName] = autoCreateJmsQueues
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoDeleteJmsQueues(addressName string, autoDeleteJmsQueues bool) *BrokerDeploymentWrapper {
	if bdw.autoDeleteJmsQueues == nil {
		bdw.autoDeleteJmsQueues = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoDeleteJmsQueues[addressName] = autoDeleteJmsQueues
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoCreateJmsTopics(addressName string, autoCreateJmsTopics bool) *BrokerDeploymentWrapper {
	if bdw.autoCreateJmsTopics == nil {
		bdw.autoCreateJmsTopics = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoCreateJmsTopics[addressName] = autoCreateJmsTopics
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoDeleteJmsTopics(addressName string, autoDeleteJmsTopics bool) *BrokerDeploymentWrapper {
	if bdw.autoDeleteJmsTopics == nil {
		bdw.autoDeleteJmsTopics = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoDeleteJmsTopics[addressName] = autoDeleteJmsTopics
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoCreateQueues(addressName string, autoCreateQueues bool) *BrokerDeploymentWrapper {
	if bdw.autoCreateQueues == nil {
		bdw.autoCreateQueues = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoCreateQueues[addressName] = autoCreateQueues
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoDeleteQueues(addressName string, autoDeleteQueues bool) *BrokerDeploymentWrapper {
	if bdw.autoDeleteQueues == nil {
		bdw.autoDeleteQueues = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoDeleteQueues[addressName] = autoDeleteQueues
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoDeleteCreatedQueues(addressName string, autoDeleteCreatedQueues bool) *BrokerDeploymentWrapper {
	if bdw.autoDeleteCreatedQueues == nil {
		bdw.autoDeleteCreatedQueues = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoDeleteCreatedQueues[addressName] = autoDeleteCreatedQueues
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoDeleteQueuesDelay(addressName string, autoDeleteQueuesDelay int32) *BrokerDeploymentWrapper {
	if bdw.autoDeleteQueuesDelay == nil {
		bdw.autoDeleteQueuesDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoDeleteQueuesDelay[addressName] = autoDeleteQueuesDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAudoDeleteQueuesMessageCount(addressName string, audoDeleteQueuesMessageCount int32) *BrokerDeploymentWrapper {
	if bdw.audoDeleteQueuesMessageCount == nil {
		bdw.audoDeleteQueuesMessageCount = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.audoDeleteQueuesMessageCount[addressName] = audoDeleteQueuesMessageCount
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithConfigDeleteQueues(addressName string, configDeleteQueues ConfigDeleteQueues) *BrokerDeploymentWrapper {
	if bdw.configDeleteQueues == nil {
		bdw.configDeleteQueues = map[string]ConfigDeleteQueues{}
	}
	bdw.addKnownAddress(addressName)
	bdw.configDeleteQueues[addressName] = configDeleteQueues
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoCreateAddresses(addressName string, autoCreateAddresses bool) *BrokerDeploymentWrapper {
	if bdw.autoCreateAddresses == nil {
		bdw.autoCreateAddresses = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoCreateAddresses[addressName] = autoCreateAddresses
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAudoDeleteAddresses(addressName string, audoDeleteAddresses bool) *BrokerDeploymentWrapper {
	if bdw.audoDeleteAddresses == nil {
		bdw.audoDeleteAddresses = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.audoDeleteAddresses[addressName] = audoDeleteAddresses
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithAutoDeleteAddressesDelay(addressName string, autoDeleteAddressesDelay int32) *BrokerDeploymentWrapper {
	if bdw.autoDeleteAddressesDelay == nil {
		bdw.autoDeleteAddressesDelay = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.autoDeleteAddressesDelay[addressName] = autoDeleteAddressesDelay
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithConfigDeleteAddresses(addressName string, configDeleteAddresses ConfigDeleteAddresses) *BrokerDeploymentWrapper {
	if bdw.configDeleteAddresses == nil {
		bdw.configDeleteAddresses = map[string]ConfigDeleteAddresses{}
	}
	bdw.addKnownAddress(addressName)
	bdw.configDeleteAddresses[addressName] = configDeleteAddresses
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithManagementBrowsePageSize(addressName string, managementBrowsePageSize int32) *BrokerDeploymentWrapper {
	if bdw.managementBrowsePageSize == nil {
		bdw.managementBrowsePageSize = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.managementBrowsePageSize[addressName] = managementBrowsePageSize
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultPurgeOnNoConsumers(addressName string, defaultPurgeOnNoConsumers bool) *BrokerDeploymentWrapper {
	if bdw.defaultPurgeOnNoConsumers == nil {
		bdw.defaultPurgeOnNoConsumers = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultPurgeOnNoConsumers[addressName] = defaultPurgeOnNoConsumers
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultMaxConsumers(addressName string, defaultMaxConsumers int32) *BrokerDeploymentWrapper {
	if bdw.defaultMaxConsumers == nil {
		bdw.defaultMaxConsumers = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultMaxConsumers[addressName] = defaultMaxConsumers
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultQueueRoutingType(addressName string, defaultQueueRoutingType RoutingType) *BrokerDeploymentWrapper {
	if bdw.defaultQueueRoutingType == nil {
		bdw.defaultQueueRoutingType = map[string]RoutingType{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultQueueRoutingType[addressName] = defaultQueueRoutingType
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultAddressRoutingType(addressName string, defaultAddressRoutingType RoutingType) *BrokerDeploymentWrapper {
	if bdw.defaultAddressRoutingType == nil {
		bdw.defaultAddressRoutingType = map[string]RoutingType{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultAddressRoutingType[addressName] = defaultAddressRoutingType
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultConsumerWindowSize(addressName string, defaultConsumerWindowSize int32) *BrokerDeploymentWrapper {
	if bdw.defaultConsumerWindowSize == nil {
		bdw.defaultConsumerWindowSize = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultConsumerWindowSize[addressName] = defaultConsumerWindowSize
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultRingSize(addressName string, defaultRingSize int32) *BrokerDeploymentWrapper {
	if bdw.defaultRingSize == nil {
		bdw.defaultRingSize = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultRingSize[addressName] = defaultRingSize
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithDefaultRetroMessageCount(addressName string, defaultRetroMessageCount int32) *BrokerDeploymentWrapper {
	if bdw.defaultRetroMessageCount == nil {
		bdw.defaultRetroMessageCount = map[string]int32{}
	}
	bdw.addKnownAddress(addressName)
	bdw.defaultRetroMessageCount[addressName] = defaultRetroMessageCount
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithEnableMetrics(addressName string, enableMetrics bool) *BrokerDeploymentWrapper {
	if bdw.enableMetrics == nil {
		bdw.enableMetrics = map[string]bool{}
	}
	bdw.addKnownAddress(addressName)
	bdw.enableMetrics[addressName] = enableMetrics
	return bdw
}
