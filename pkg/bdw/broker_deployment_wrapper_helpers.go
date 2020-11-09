package bdw

/* This file contains non-deployment helper methods for BrokerDeploymentWrapper
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	brokerv1 "github.com/artemiscloud/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	brokerv3 "github.com/artemiscloud/activemq-artemis-operator/pkg/apis/broker/v2alpha3"
	"github.com/onsi/gomega"
	//"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (bdw *BrokerDeploymentWrapper) GetPodList() *corev1.PodList {
	getopts := v1.GetOptions{}
	statefulSet, err := bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Get(bdw.name+"-ss", getopts)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	listOptions := v1.ListOptions{LabelSelector: statefulSet.Name}
	pods, err := bdw.ctx1.Clients.KubeClient.CoreV1().Pods(bdw.ctx1.Namespace).List(listOptions)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return pods
}

func (bdw *BrokerDeploymentWrapper) SetEnvVariable(name, value string) {
	getopts := v1.GetOptions{}
	statefulSet, err := bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Get(bdw.name+"-ss", getopts)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	env := statefulSet.Spec.Template.Spec.Containers[0].Env
	statefulSet.Spec.Template.Spec.Containers[0].Env = append(env, corev1.EnvVar{Name: name, Value: value})
	_, err = bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Update(statefulSet)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

//This expects to be ran on openshift.
func (bdw *BrokerDeploymentWrapper) GetExternalUrls(filter string, podNumber int) ([]string, error) {
	var result []string
	routes, _ := bdw.ctx1.Clients.OcpClient.RoutesClient.RouteV1().Routes(bdw.ctx1.Namespace).List(v1.ListOptions{})
	for _, route := range routes.Items {
		url := route.Spec.Host
		if strings.Contains(url, filter) && strings.Contains(url, fmt.Sprintf("-%d-svc", podNumber)) {
			result = append(result, url)
		}
	}
	if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("no URLs were returned with filter: %s for pod number %d", filter, podNumber))
	}
	return result, nil
}

//We always configure Artemis as if it is latest API version
func (bdw *BrokerDeploymentWrapper) ConfigureBroker(artemis *brokerv3.ActiveMQArtemis, acceptorType AcceptorType) *brokerv3.ActiveMQArtemis {
	artemis.Spec.DeploymentPlan.Size = int32(bdw.deploymentSize)
	if acceptorType != NoChangeAcceptor {
		artemis.Spec.Acceptors = append(artemis.Spec.Acceptors, *acceptors[acceptorType])
	}
	for num := range artemis.Spec.Acceptors {
		artemis.Spec.Acceptors[num].SSLEnabled = bdw.sslEnabled
	}
	for num := range artemis.Spec.Connectors {
		artemis.Spec.Connectors[num].SSLEnabled = bdw.sslEnabled
	}
	artemis.Spec.DeploymentPlan.MessageMigration = &bdw.migration
	artemis.Spec.DeploymentPlan.PersistenceEnabled = bdw.persistence
	artemis.Spec.AdminUser = test.Username
	artemis.Spec.AdminPassword = test.Password
	artemis.Spec.DeploymentPlan.Image = bdw.customImage
	artemis.ObjectMeta.Name = bdw.name
	artemis.Spec.Console.Expose = bdw.exposeConsole
	artemis.Spec.DeploymentPlan.Storage.Size = bdw.storageSize

	addressSettingsArray := []brokerv3.AddressSettingType{}
	for _, addressName := range bdw.knownAddresses {
		addressSettingsItem := bdw.fillAddressSetting(addressName)
		addressSettingsArray = append(addressSettingsArray, addressSettingsItem)
	}

	artemis.Spec.AddressSettings.AddressSetting = addressSettingsArray
	return artemis
}

func (bdw *BrokerDeploymentWrapper) fillAddressSetting(addressName string) brokerv3.AddressSettingType {
	maxSizeBytes := bdw.maxSizeBytes[addressName]
	deadLetterAddress := bdw.deadLetterAddress[addressName]
	autoCreateDeadResources := bdw.autoCreateDeadLetterResources[addressName]
	dlqPrefix := bdw.dlqPrefix[addressName]
	dlqSuffix := bdw.dlqSuffix[addressName]
	expiryAddress := bdw.expiryAddress[addressName]
	//autoCreateExpiryResources := bdw.autoCreateDeadLetterResources[addressName]
	expiryPrefix := bdw.expiryPrefix[addressName]
	expirySuffix := bdw.expirySuffix[addressName]
	expiryDelay := bdw.expiryDelay[addressName]
	minExpiryDelay := bdw.minExpiryDelay[addressName]
	maxExpiryDelay := bdw.maxExpiryDelay[addressName]
	redeliverDelay := bdw.redeliveryDelay[addressName]
	redeliveryDelayMax := bdw.maxRedeliveryDelay[addressName]
	redeliveryAttemptsMax := bdw.maxRedeliveryAttempts[addressName]
	redeliveryDelayMult := bdw.redeliveryDelayMult[addressName]
	redeliveryDelayAvoidance := bdw.redeliveryCollisionsAvoidance[addressName]
	maxRejectThreshold := bdw.maxSizeBytesRejectThreshold[addressName]
	//pageSizeBytes := bdw.pageSizeBytes[addressName]
	pageMaxCacheSize := bdw.pageMaxCacheSize[addressName]
	addressFullPolicy := bdw.addressFullPolicy[addressName].String()
	messageCounterHistoryLimit := bdw.messageCounterHistoryDayLimit[addressName]
	lastValueQueue := bdw.lastValueQueue[addressName]
	defaultLastValueQueue := bdw.defaultLastValueQueue[addressName]
	defaultLastValueKey := bdw.defaultLastValueKey[addressName]
	defaultNonDestructive := bdw.defaultNonDestructive[addressName]
	defaultExclusiveQueue := bdw.defaultExclusiveQueue[addressName]
	defaultGroupRebalance := bdw.defaultGroupRebalance[addressName]
	//defaultGroupRebalancePause := bdw.defaultGroupRebalancePauseDispatch[addressName]
	defaultGroupBuckets := bdw.defaultGroupBuckets[addressName]
	defaultGroupFirstKey := bdw.defaultGroupFirstKey[addressName]
	defaultConsumerBeforeDispatch := bdw.defaultConsumerBeforeDispatch[addressName]
	defaultDelayBeforeDispatch := bdw.defaultDelayBeforeDispatch[addressName]
	redistributionDelay := bdw.redistributionDelay[addressName]
	sendToDLAOnNoRoute := bdw.sendToDLAOnNoRoute[addressName]
	slowConsumerThreshold := bdw.slowConsumerThreshold[addressName]
	slowConsumerPolicy := bdw.slowConsumerPolicy[addressName].String()
	slowConsumerCheckPeriod := bdw.slowConsumerCheckPeriod[addressName]
	//autoCreateJMSQueues := bdw.autoCreateJmsQueues[addressName]
	//autoCreateJmsTopics := bdw.autoCreateJmsTopics[addressName]
	//autoDeleteJmsQueues := bdw.autoDeleteJmsQueues[addressName]
	//autoDeleteJmsTopics := bdw.autoDeleteJmsTopics[addressName]
	//autoCreateQueues := bdw.autoCreateQueues[addressName]
	//autoDeleteQueues := bdw.autoDeleteQueues[addressName]
	//autoDeleteCreatedQueues := bdw.autoDeleteCreatedQueues[addressName]
	//autoDeleteQueuesDelay := bdw.autoDeleteQueuesDelay[addressName]
	//autoDeleteQueuesMessageCount := bdw.audoDeleteQueuesMessageCount[addressName]
	//configDeleteQueues := bdw.configDeleteQueues[addressName].String()
	//configDeleteAddresses := bdw.configDeleteAddresses[addressName].String()
	//autoDeleteAddresses := bdw.autoDeleteAddresses[addressName]
	//autoDeleteAddressDelay := bdw.autoDeleteAddressesDelay[addressName]
	//autoCreateAddresses := bdw.autoCreateAddresses[addressName]
	managementBrowsePageSize := bdw.managementBrowsePageSize[addressName]
	defaultPurgeOnNoConsumers := bdw.defaultPurgeOnNoConsumers[addressName]
	defaultMaxConsumers := bdw.defaultMaxConsumers[addressName]
	defaultQueueRoutingType := bdw.defaultQueueRoutingType[addressName].String()
	defaultAddressRoutingType := bdw.defaultAddressRoutingType[addressName].String()
	defaultConsumerWindowSize := bdw.defaultConsumerWindowSize[addressName]
	defaultRingSize := bdw.defaultRingSize[addressName]
	retroactiveMessageCount := bdw.defaultRetroMessageCount[addressName]
	enableMetrics := bdw.enableMetrics[addressName]

	return brokerv3.AddressSettingType{
		DeadLetterAddress:                  &deadLetterAddress,
		AutoCreateDeadLetterResources:      &autoCreateDeadResources,
		DeadLetterQueuePrefix:              &dlqPrefix,
		DeadLetterQueueSuffix:              &dlqSuffix,
		ExpiryAddress:                      &expiryAddress,
		//AutoCreateExpiryResources:          &autoCreateExpiryResources,
		ExpiryQueuePrefix:                  &expiryPrefix,
		ExpiryQueueSuffix:                  &expirySuffix,
		ExpiryDelay:                        &expiryDelay,
		MinExpiryDelay:                     &minExpiryDelay,
		MaxExpiryDelay:                     &maxExpiryDelay,
		RedeliveryDelay:                    &redeliverDelay,
		RedeliveryDelayMultiplier:          &redeliveryDelayMult,
		RedeliveryCollisionAvoidanceFactor: &redeliveryDelayAvoidance,
		MaxRedeliveryDelay:                 &redeliveryDelayMax,
		MaxDeliveryAttempts:                &redeliveryAttemptsMax,
		MaxSizeBytes:                       &maxSizeBytes,
		MaxSizeBytesRejectThreshold:        &maxRejectThreshold,
		//PageSizeBytes:                      &pageSizeBytes, //TODO: this is bugged in operator/crd
		PageMaxCacheSize:                   &pageMaxCacheSize,
		AddressFullPolicy:                  &addressFullPolicy,
		MessageCounterHistoryDayLimit:      &messageCounterHistoryLimit,
		LastValueQueue:                     &lastValueQueue,
		DefaultLastValueQueue:              &defaultLastValueQueue,
		DefaultLastValueKey:                &defaultLastValueKey,
		DefaultNonDestructive:              &defaultNonDestructive,
		DefaultExclusiveQueue:              &defaultExclusiveQueue,
		DefaultGroupRebalance:              &defaultGroupRebalance,
		//DefaultGroupRebalancePauseDispatch: &defaultGroupRebalancePause,
		DefaultGroupBuckets:                &defaultGroupBuckets,
		DefaultGroupFirstKey:               &defaultGroupFirstKey,
		DefaultConsumersBeforeDispatch:     &defaultConsumerBeforeDispatch,
		DefaultDelayBeforeDispatch:         &defaultDelayBeforeDispatch,
		RedistributionDelay:                &redistributionDelay,
		SendToDlaOnNoRoute:                 &sendToDLAOnNoRoute,
		SlowConsumerThreshold:              &slowConsumerThreshold,
		SlowConsumerPolicy:                 &slowConsumerPolicy,
		SlowConsumerCheckPeriod:            &slowConsumerCheckPeriod,
		//AutoCreateJmsQueues:                &autoCreateJMSQueues,
		//AutoDeleteJmsQueues:                &autoDeleteJmsQueues,
		//AutoCreateJmsTopics:                &autoCreateJmsTopics,
		//AutoDeleteJmsTopics:                &autoDeleteJmsTopics,
		//AutoCreateQueues:                   &autoCreateQueues,
		//AutoDeleteQueues:                   &autoDeleteQueues,
		//AutoDeleteCreatedQueues:            &autoDeleteCreatedQueues,
		//AutoDeleteQueuesDelay:              &autoDeleteQueuesDelay,
		//AutoDeleteQueuesMessageCount:       &autoDeleteQueuesMessageCount,
		//ConfigDeleteQueues:                 &configDeleteQueues,
		//AutoCreateAddresses:                &autoCreateAddresses,
		//AutoDeleteAddresses:                &autoDeleteAddresses,
		//AutoDeleteAddressesDelay:           &autoDeleteAddressDelay,
		//ConfigDeleteAddresses:              nil,
		ManagementBrowsePageSize:           &managementBrowsePageSize,
		DefaultPurgeOnNoConsumers:          &defaultPurgeOnNoConsumers,
		DefaultMaxConsumers:                &defaultMaxConsumers,
		DefaultQueueRoutingType:            &defaultQueueRoutingType,
		DefaultAddressRoutingType:          &defaultAddressRoutingType,
		DefaultConsumerWindowSize:          &defaultConsumerWindowSize,
		DefaultRingSize:                    &defaultRingSize,
		RetroactiveMessageCount:            &retroactiveMessageCount,
		EnableMetrics:                      &enableMetrics,
		Match:                              addressName,
	}
}

func (bdw *BrokerDeploymentWrapper) SetUpDefaultAddressSettings(addressName string) *BrokerDeploymentWrapper {
        bdw.WithAddressFullPolicy(addressName, FailPolicy).
        WithAddressPolicy(addressName, FailPolicy).
        WithAddressSize(addressName, DEFAULT_SIZE_BIG).
        WithAudoDeleteAddresses(addressName, false).
        WithAudoDeleteQueuesMessageCount(addressName, DEFAULT_COUNT).
        WithAutoCreateAddresses(addressName, true).
        WithAutoCreateDeadLetterResources(addressName, true).
        WithAutoCreateExpiryResources(addressName, true).
        WithAutoCreateJmsQueues(addressName, true).
        WithAutoCreateJmsTopics(addressName, true).
        WithAutoCreateQueues(addressName, true).
        WithAutoDeleteAddressesDelay(addressName, DEFAULT_DELAY).
        WithAutoDeleteCreatedQueues(addressName, false).
        WithAutoDeleteJmsQueues(addressName, false).
        WithAutoDeleteJmsTopics(addressName, false).
        WithAutoDeleteQueues(addressName, false).
        WithAutoDeleteQueuesDelay(addressName, DEFAULT_DELAY).
        WithConfigDeleteAddresses(addressName, Off).
        WithConfigDeleteQueues(addressName, Off).
        WithDeadLetterAddress(addressName, DEFAULT_DEAD_ADDRESS).
        WithDefaultAddressRoutingType(addressName, Multicast).
        WithDefaultConsumerBeforeDispatch(addressName, DEFAULT_COUNT).
        WithDefaultConsumerWindowSize(addressName, DEFAULT_COUNT).
        WithDefaultDelayBeforeDispatch(addressName, DEFAULT_DELAY).
        WithDefaultExclusiveQueue(addressName, false).
        WithDefaultGroupBuckets(addressName, DEFAULT_COUNT).
        WithDefaultGroupFirstKey(addressName, DEFAULT_KEY).
        WithDefaultGroupRebalance(addressName, false).
        WithDefaultGroupRebalancePauseDispatch(addressName, true).
        WithDefaultLastValueKey(addressName, DEFAULT_KEY).
        WithDefaultLastValueQueue(addressName, false).
        WithDefaultMaxConsumers(addressName, DEFAULT_COUNT).
        WithDefaultNonDestructive(addressName, true).
        WithDefaultPurgeOnNoConsumers(addressName, false).
        WithDefaultQueueRoutingType(addressName, Multicast).
        WithDefaultRetroMessageCount(addressName, DEFAULT_COUNT).
        WithDefaultRingSize(addressName, DEFAULT_COUNT).
        WithDlqPrefix(addressName, DEFAULT_PREFIX).
        WithDlqSuffix(addressName, DEFAULT_SUFFIX).
        WithEnableMetrics(addressName, false).
        WithExpiryAddress(addressName, DEFAULT_EXPIRY_ADDRESS).
        WithExpiryPrefix(addressName, DEFAULT_PREFIX).
        WithExpirySuffix(addressName, DEFAULT_SUFFIX).
        WithLastValueQueue(addressName, true).
        WithManagementBrowsePageSize(addressName, 1000).
        WithMaxExpiryDelay(addressName, DEFAULT_DELAY).
        WithMaxRedeliveryAttempts(addressName, DEFAULT_COUNT).
        WithMaxRedeliveryDelay(addressName, DEFAULT_DELAY).
        WithMaxSizeBytes(addressName, DEFAULT_SIZE_SMALL).
        WithMaxSizeBytesRejectThreshold(addressName, 10000).
        WithMessageCounterHistoryDayLimit(addressName, DEFAULT_COUNT).
        WithMinExpiryDelay(addressName, DEFAULT_DELAY).
        WithPageMaxCacheSize(addressName, 20000000).
        WithPageSizeBytes(addressName, DEFAULT_SIZE_SMALL).
        WithRedeliveryCollisionsAvoidance(addressName, DEFAULT_COUNT).
        WithRedeliveryDelay(addressName, DEFAULT_DELAY).
        WithRedeliveryDelayMult(addressName, 2).
        WithRedistributionDelay(addressName, DEFAULT_DELAY).
        WithSendToDLAOnNoRoute(addressName, false).
        WithSlowConsumerCheckPeriod(addressName, DEFAULT_PERIOD).
        WithSlowConsumerPolicy(addressName, Notify).
        WithSlowConsumerThreshold(addressName, DEFAULT_DELAY)
        
        return bdw
}


const (
    DEFAULT_DELAY = 1000
    DEFAULT_COUNT = 100
    DEFAULT_SIZE_BIG = "2G"
    DEFAULT_SIZE_SMALL = "10K"
    DEFAULT_DEAD_ADDRESS = "DLQ"
    DEFAULT_EXPIRY_ADDRESS = "expiry"
    DEFAULT_KEY = "abc"
    DEFAULT_SUFFIX = "suffix"
    DEFAULT_PREFIX = "prefix"
    DEFAULT_PERIOD = DEFAULT_DELAY
)

func (bdw *BrokerDeploymentWrapper) ConvertToV1(artemisOriginal *brokerv3.ActiveMQArtemis) *brokerv1.ActiveMQArtemis {
	artemisResult := &brokerv1.ActiveMQArtemis{}
	data, err := json.Marshal(artemisOriginal)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, artemisResult)
	if err != nil {
		panic(err)
	}
	return artemisResult
}
