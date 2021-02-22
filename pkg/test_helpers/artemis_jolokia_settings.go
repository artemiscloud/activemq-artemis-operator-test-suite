package test_helpers

type JolokiaBrokerSettings struct {
	Request   interface{} `json:"request"`
	Value     string      `json:"value"`
	Timestamp float64     `json:"timestamp"`
	Status    int         `json:"status"`
}

type Value struct {
	DLA                                string  `json:"DLA"`
	ExpiryAddress                      string  `json:"expiryAddress"`
	ExpiryDelay                        int     `json:"expiryDelay"`
	MinExpiryDelay                     int     `json:"minExpiryDelay"`
	MaxExpiryDelay                     int     `json:"maxExpiryDelay"`
	MaxDeliveryAttempts                int     `json:"maxDeliveryAttempts"`
	PageCacheMaxSize                   int     `json:"pageCacheMaxSize"`
	MaxSizeBytes                       int     `json:"maxSizeBytes"`
	PageSizeBytes                      int     `json:"pageSizeBytes"`
	RedeliveryDelay                    int     `json:"redeliveryDelay"`
	RedeliveryMultiplier               float32 `json:"redeliveryMultiplier"`
	MaxRedeliveryDelay                 int     `json:"maxRedeliveryDelay"`
	RedistributionDelay                int     `json:"redistributionDelay"`
	LastValueQueue                     bool    `json:"lastValueQueue"`
	SendToDLAOnNoRoute                 bool    `json:"sendToDLAOnNoRoute"`
	AddressFullMessagePolicy           string  `json:"addressFullMessagePolicy"`
	SlowConsumerThreshold              int     `json:"slowConsumerThreshold"`
	SlowConsumerCheckPeriod            int     `json:"slowConsumerCheckPeriod"`
	SlowConsumerPolicy                 string  `json:"slowConsumerPolicy"`
	AutoCreateJmsQueues                bool    `json:"autoCreateJmsQueues"`
	AutoDeleteJmsQueues                bool    `json:"autoDeleteJmsQueues"`
	AutoCreateJmsTopics                bool    `json:"autoCreateJmsTopics"`
	AutoDeleteJmsTopics                bool    `json:"autoDeleteJmsTopics"`
	AutoCreateQueues                   bool    `json:"autoCreateQueues"`
	AutoDeleteQueues                   bool    `json:"autoDeleteQueues"`
	AutoCreateAddresses                bool    `json:"autoCreateAddresses"`
	AutoDeleteAddresses                bool    `json:"autoDeleteAddresses"`
	ConfigDeleteQueues                 string  `json:"configDeleteQueues"`
	ConfigDeleteAddresses              string  `json:"configDeleteAddresses"`
	MaxSizeBytesRejectThreshold        int     `json:"maxSizeBytesRejectThreshold"`
	DefaultLastValueKey                string  `json:"defaultLastValueKey"`
	DefaultNonDestructive              bool    `json:"defaultNonDestructive"`
	DefaultExclusiveQueue              bool    `json:"defaultExclusiveQueue"`
	DefaultGroupRebalance              bool    `json:"defaultGroupRebalance"`
	DefaultGroupRebalancePauseDispatch bool    `json:"defaultGroupRebalancePauseDispatch"`
	DefaultGroupBuckets                int     `json:"defaultGroupBuckets"`
	DefaultGroupFirstKey               string  `json:"defaultGroupFirstKey"`
	DefaultMaxConsumers                int     `json:"defaultMaxConsumers"`
	DefaultPurgeOnNoConsumers          bool    `json:"defaultPurgeOnNoConsumers"`
	DefaultConsumersBeforeDispatch     int     `json:"defaultConsumersBeforeDispatch"`
	DefaultDelayBeforeDispatch         int     `json:"defaultDelayBeforeDispatch"`
	DefaultQueueRoutingType            string  `json:"defaultQueueRoutingType"`
	DefaultAddressRoutingType          string  `json:"defaultAddressRoutingType"`
	DefaultConsumerWindowSize          int     `json:"defaultConsumerWindowSize"`
	DefaultRingSize                    int     `json:"defaultRingSize"`
	AutoDeleteCreatedQueues            bool    `json:"autoDeleteCreatedQueues"`
	AutoDeleteQueuesDelay              int     `json:"autoDeleteQueuesDelay"`
	AutoDeleteQueuesMessageCount       int     `json:"autoDeleteQueuesMessageCount"`
	AutoDeleteAddressesDelay           int     `json:"autoDeleteAddressesDelay"`
	RedeliveryCollisionAvoidanceFactor float32 `json:"redeliveryCollisionAvoidanceFactor"`
	RetroactiveMessageCount            int     `json:"retroactiveMessageCount"`
	AutoCreateDeadLetterResources      bool    `json:"autoCreateDeadLetterResources"`
	DeadLetterQueuePrefix              string  `json:"deadLetterQueuePrefix"`
	DeadLetterQueueSuffix              string  `json:"deadLetterQueueSuffix"`
	AutoCreateExpiryResources          bool    `json:"autoCreateExpiryResources"`
	ExpiryQueuePrefix                  string  `json:"expiryQueuePrefix"`
	ExpiryQueueSuffix                  string  `json:"expiryQueueSuffix"`
	EnableMetrics                      bool    `json:"enableMetrics"`
}
