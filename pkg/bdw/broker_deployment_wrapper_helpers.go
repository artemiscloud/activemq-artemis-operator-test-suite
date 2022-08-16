package bdw

/* This file contains non-deployment helper methods for BrokerDeploymentWrapper
 */

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	brokerbeta "github.com/artemiscloud/activemq-artemis-operator/api/v1beta1"
	brokerv1 "github.com/artemiscloud/activemq-artemis-operator/api/v2alpha1"
	brokerv3 "github.com/artemiscloud/activemq-artemis-operator/api/v2alpha3"
	"github.com/rh-messaging/shipshape/pkg/framework/log"

	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	//"github.com/rh-messaging/shipshape/pkg/framework/log"
	"strings"

	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (bdw *BrokerDeploymentWrapper) AddProperty(name, value string) {
	if bdw.properties == nil {
		bdw.properties = make(map[string]string)
	}
	bdw.properties[name] = value
}

func (bdw *BrokerDeploymentWrapper) GetStatefulSet() *appsv1.StatefulSet {
	statefulSet, err := bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Get(context.TODO(), bdw.name+"-ss", v1.GetOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "can't retrieve statefulset")
	return statefulSet
}

func (bdw *BrokerDeploymentWrapper) GetFile(podname, containername, filename string, restconfig rest.Config) (string, error) {
	log.Logf("restconfig: %v", restconfig)
	coreclient := bdw.ctx1.Clients.KubeClient.CoreV1()
	reader, outStream := io.Pipe()
	cmd := []string{"cat", filename}
	req := coreclient.RESTClient().
		Get().
		Namespace(bdw.ctx1.Namespace).
		Resource("pods").
		Name(podname).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containername,
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	restconfig.TLSClientConfig.Insecure = true
	exec, err := remotecommand.NewSPDYExecutor(&restconfig, "POST", req.URL())
	if err != nil {
		return "", err
	}
	go func() {
		defer outStream.Close()
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
			Tty:    false,
		})
		if err != nil {
			panic(err)
		}
	}()
	outC := ""
	var buf bytes.Buffer
	io.Copy(&buf, reader)
	outC = buf.String()
	return outC, err
}

func (bdw *BrokerDeploymentWrapper) CreateEmptySecurityCR() brokerbeta.ActiveMQArtemisSecurity {
	return brokerbeta.ActiveMQArtemisSecurity{}
}

func (bdw *BrokerDeploymentWrapper) CreateDefaultSecurityCR(name, username, pass string, roles []string) brokerbeta.ActiveMQArtemisSecurity {
	cr := bdw.CreateEmptySecurityCR()
	ploginmodule := brokerbeta.PropertiesLoginModuleType{}
	cr.Name = name
	ploginmodule.Name = name
	user := brokerbeta.UserType{}
	user.Name = username
	user.Password = &pass
	user.Roles = roles
	ploginmodule.Users = append(ploginmodule.Users, user)
	loginmodule := brokerbeta.LoginModulesType{}
	loginmodule.PropertiesLoginModules = append(loginmodule.PropertiesLoginModules, ploginmodule)
	cr.Spec.LoginModules = loginmodule
	secdom := brokerbeta.SecurityDomainsType{}
	secdom.BrokerDomain.Name = &name
	ref := brokerbeta.LoginModuleReferenceType{}
	suff := "sufficient"
	ref.Flag = &suff
	ref.Name = &name
	secdom.BrokerDomain.LoginModules = append(secdom.BrokerDomain.LoginModules, ref)
	cr.Spec.SecurityDomains = secdom
	return cr
}

func (bdw *BrokerDeploymentWrapper) GetPodList() *corev1.PodList {
	getopts := v1.GetOptions{}
	statefulSet, err := bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Get(context.TODO(), bdw.name+"-ss", getopts)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	listOptions := v1.ListOptions{LabelSelector: statefulSet.Name}
	pods, err := bdw.ctx1.Clients.KubeClient.CoreV1().Pods(bdw.ctx1.Namespace).List(context.TODO(), listOptions)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return pods
}

func (bdw *BrokerDeploymentWrapper) SetEnvVariable(name, value string) {
	getopts := v1.GetOptions{}
	statefulSet, err := bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Get(context.TODO(), bdw.name+"-ss", getopts)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	env := statefulSet.Spec.Template.Spec.Containers[0].Env
	statefulSet.Spec.Template.Spec.Containers[0].Env = append(env, corev1.EnvVar{Name: name, Value: value})
	_, err = bdw.ctx1.Clients.KubeClient.AppsV1().StatefulSets(bdw.ctx1.Namespace).Update(context.TODO(), statefulSet, v1.UpdateOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

//This expects to be ran on openshift.
func (bdw *BrokerDeploymentWrapper) GetExternalUrls(filter string, podNumber int) ([]string, error) {
	var result []string
	routes, _ := bdw.ctx1.Clients.OcpClient.RoutesClient.RouteV1().Routes(bdw.ctx1.Namespace).List(context.TODO(), v1.ListOptions{})
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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

//We always configure Artemis as if it is latest API version
func (bdw *BrokerDeploymentWrapper) ConfigureBroker(artemis *brokerbeta.ActiveMQArtemis, acceptorType AcceptorType) *brokerbeta.ActiveMQArtemis {
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
	if bdw.customImage != "" {
		artemis.Spec.DeploymentPlan.Image = bdw.customImage
	}
	if bdw.customInit != "" {
		artemis.Spec.DeploymentPlan.InitImage = bdw.customInit
	}
	artemis.ObjectMeta.Name = bdw.name
	artemis.Spec.Console.Expose = bdw.exposeConsole
	artemis.Spec.DeploymentPlan.Storage.Size = bdw.storageSize
	var processed []string
	addressSettingsArray := []brokerbeta.AddressSettingType{}
	for _, addressName := range bdw.knownAddresses {
		if !contains(processed, addressName) {
			processed = append(processed, addressName)
			addressSettingsItem := bdw.fillAddressSetting(addressName)
			addressSettingsArray = append(addressSettingsArray, addressSettingsItem)
		}
	}

	artemis.Spec.DeploymentPlan.Resources.Limits = getResourceList(bdw.ResourcesLimits.cpu, bdw.ResourcesLimits.mem)
	artemis.Spec.DeploymentPlan.Resources.Requests = getResourceList(bdw.ResourcesRequests.cpu, bdw.ResourcesRequests.mem)
	artemis.Spec.AddressSettings.AddressSetting = addressSettingsArray
	if len(bdw.properties) != 0 {
		for it := range bdw.properties {
			prop := fmt.Sprintf("%s=%s", it, bdw.properties[it])
			artemis.Spec.BrokerProperties = append(artemis.Spec.BrokerProperties, prop)
		}
	}
	if len(bdw.version) != 0 {
		artemis.Spec.Version = bdw.version
		artemis.Spec.Upgrades.Enabled = true
		artemis.Spec.Upgrades.Minor = true
	}

	return artemis
}

func getResourceList(cpu, memory string) corev1.ResourceList {
	res := corev1.ResourceList{}
	if cpu != "" {
		res[corev1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[corev1.ResourceMemory] = resource.MustParse(memory)
	}
	return res
}

func (bdw *BrokerDeploymentWrapper) fillAddressSetting(addressName string) brokerbeta.AddressSettingType {
	maxSizeBytes := bdw.maxSizeBytes[addressName]
	deadLetterAddress := bdw.deadLetterAddress[addressName]
	autoCreateDeadResources := bdw.autoCreateDeadLetterResources[addressName]
	dlqPrefix := bdw.dlqPrefix[addressName]
	dlqSuffix := bdw.dlqSuffix[addressName]
	expiryAddress := bdw.expiryAddress[addressName]
	autoCreateExpiryResources := bdw.autoCreateExpiryResources[addressName]
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
	defaultGroupRebalancePause := bdw.defaultGroupRebalancePauseDispatch[addressName]
	defaultGroupBuckets := bdw.defaultGroupBuckets[addressName]
	defaultGroupFirstKey := bdw.defaultGroupFirstKey[addressName]
	defaultConsumerBeforeDispatch := bdw.defaultConsumerBeforeDispatch[addressName]
	defaultDelayBeforeDispatch := bdw.defaultDelayBeforeDispatch[addressName]
	redistributionDelay := bdw.redistributionDelay[addressName]
	sendToDLAOnNoRoute := bdw.sendToDLAOnNoRoute[addressName]
	slowConsumerThreshold := bdw.slowConsumerThreshold[addressName]
	slowConsumerPolicy := bdw.slowConsumerPolicy[addressName].String()
	slowConsumerCheckPeriod := bdw.slowConsumerCheckPeriod[addressName]
	if slowConsumerCheckPeriod == 0 {
		slowConsumerCheckPeriod++
	}
	if slowConsumerThreshold == 0 {
		slowConsumerThreshold = -1
	}
	autoCreateJMSQueues := bdw.autoCreateJmsQueues[addressName]
	autoCreateJmsTopics := bdw.autoCreateJmsTopics[addressName]
	autoDeleteJmsQueues := bdw.autoDeleteJmsQueues[addressName]
	autoDeleteJmsTopics := bdw.autoDeleteJmsTopics[addressName]
	autoCreateQueues := bdw.autoCreateQueues[addressName]
	autoDeleteQueues := bdw.autoDeleteQueues[addressName]
	autoDeleteCreatedQueues := bdw.autoDeleteCreatedQueues[addressName]
	autoDeleteQueuesDelay := bdw.autoDeleteQueuesDelay[addressName]
	autoDeleteQueuesMessageCount := bdw.audoDeleteQueuesMessageCount[addressName]
	configDeleteQueues := bdw.configDeleteQueues[addressName].String()
	//configDeleteAddresses := bdw.configDeleteAddresses[addressName].String()
	autoDeleteAddresses := bdw.autoDeleteAddresses[addressName]
	autoDeleteAddressDelay := bdw.autoDeleteAddressesDelay[addressName]
	autoCreateAddresses := bdw.autoCreateAddresses[addressName]
	managementBrowsePageSize := bdw.managementBrowsePageSize[addressName]
	defaultPurgeOnNoConsumers := bdw.defaultPurgeOnNoConsumers[addressName]
	defaultMaxConsumers := bdw.defaultMaxConsumers[addressName]
	defaultQueueRoutingType := bdw.defaultQueueRoutingType[addressName].String()
	defaultAddressRoutingType := bdw.defaultAddressRoutingType[addressName].String()
	defaultConsumerWindowSize := bdw.defaultConsumerWindowSize[addressName]
	defaultRingSize := bdw.defaultRingSize[addressName]
	retroactiveMessageCount := bdw.defaultRetroMessageCount[addressName]
	enableMetrics := bdw.enableMetrics[addressName]

	return brokerbeta.AddressSettingType{
		DeadLetterAddress:                  &deadLetterAddress,
		AutoCreateDeadLetterResources:      &autoCreateDeadResources,
		DeadLetterQueuePrefix:              &dlqPrefix,
		DeadLetterQueueSuffix:              &dlqSuffix,
		ExpiryAddress:                      &expiryAddress,
		AutoCreateExpiryResources:          &autoCreateExpiryResources,
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
		DefaultGroupRebalancePauseDispatch: &defaultGroupRebalancePause,
		DefaultGroupBuckets:                &defaultGroupBuckets,
		DefaultGroupFirstKey:               &defaultGroupFirstKey,
		DefaultConsumersBeforeDispatch:     &defaultConsumerBeforeDispatch,
		DefaultDelayBeforeDispatch:         &defaultDelayBeforeDispatch,
		RedistributionDelay:                &redistributionDelay,
		SendToDlaOnNoRoute:                 &sendToDLAOnNoRoute,
		SlowConsumerThreshold:              &slowConsumerThreshold,
		SlowConsumerPolicy:                 &slowConsumerPolicy,
		SlowConsumerCheckPeriod:            &slowConsumerCheckPeriod,
		AutoCreateJmsQueues:                &autoCreateJMSQueues,
		AutoDeleteJmsQueues:                &autoDeleteJmsQueues,
		AutoCreateJmsTopics:                &autoCreateJmsTopics,
		AutoDeleteJmsTopics:                &autoDeleteJmsTopics,
		AutoCreateQueues:                   &autoCreateQueues,
		AutoDeleteQueues:                   &autoDeleteQueues,
		AutoDeleteCreatedQueues:            &autoDeleteCreatedQueues,
		AutoDeleteQueuesDelay:              &autoDeleteQueuesDelay,
		AutoDeleteQueuesMessageCount:       &autoDeleteQueuesMessageCount,
		ConfigDeleteQueues:                 &configDeleteQueues,
		AutoCreateAddresses:                &autoCreateAddresses,
		AutoDeleteAddresses:                &autoDeleteAddresses,
		AutoDeleteAddressesDelay:           &autoDeleteAddressDelay,
		ConfigDeleteAddresses:              nil, //This particular setting is broken in generator
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
	bdw.WithAddressFullPolicy(addressName, PagePolicy).
		WithAddressPolicy(addressName, PagePolicy).
		WithAddressSize(addressName, DEFAULT_SIZE_BIG).
		WithAutoDeleteAddresses(addressName, false).
		WithAudoDeleteQueuesMessageCount(addressName, DEFAULT_COUNT).
		WithAutoCreateAddresses(addressName, true).
		WithAutoCreateDeadLetterResources(addressName, false).
		WithAutoCreateExpiryResources(addressName, false).
		WithAutoCreateJmsQueues(addressName, true). // deprecated?
		WithAutoCreateJmsTopics(addressName, true). // deprecated?
		WithAutoCreateQueues(addressName, true).
		WithAutoDeleteAddressesDelay(addressName, 0).
		WithAutoDeleteCreatedQueues(addressName, false).
		WithAutoDeleteJmsQueues(addressName, true).
		WithAutoDeleteJmsTopics(addressName, true).
		WithAutoDeleteQueues(addressName, true).
		WithAutoDeleteQueuesDelay(addressName, 0).
		WithConfigDeleteAddresses(addressName, Off).
		WithConfigDeleteQueues(addressName, Off).
		WithDeadLetterAddress(addressName, DEFAULT_DEAD_ADDRESS).
		WithDefaultAddressRoutingType(addressName, Multicast).
		WithDefaultConsumerBeforeDispatch(addressName, 0).
		WithDefaultConsumerWindowSize(addressName, 1048576).
		WithDefaultDelayBeforeDispatch(addressName, -1).
		WithDefaultExclusiveQueue(addressName, false).
		WithDefaultGroupBuckets(addressName, -1).
		WithDefaultGroupFirstKey(addressName, DEFAULT_KEY).
		WithDefaultGroupRebalance(addressName, false).
		WithDefaultGroupRebalancePauseDispatch(addressName, false).
		WithDefaultLastValueKey(addressName, DEFAULT_KEY).
		WithDefaultLastValueQueue(addressName, false).
		WithDefaultMaxConsumers(addressName, -1).
		WithDefaultNonDestructive(addressName, false).
		WithDefaultPurgeOnNoConsumers(addressName, false).
		WithDefaultQueueRoutingType(addressName, Multicast).
		WithDefaultRetroMessageCount(addressName, 0).
		WithDefaultRingSize(addressName, -1).
		WithDlqPrefix(addressName, DEFAULT_PREFIX).
		WithDlqSuffix(addressName, DEFAULT_SUFFIX).
		WithEnableMetrics(addressName, true).
		WithExpiryAddress(addressName, DEFAULT_EXPIRY_ADDRESS).
		WithExpiryPrefix(addressName, DEFAULT_PREFIX).
		WithExpirySuffix(addressName, DEFAULT_SUFFIX).
		WithLastValueQueue(addressName, true).
		WithManagementBrowsePageSize(addressName, 1000).
		WithMaxExpiryDelay(addressName, -1).
		WithMaxRedeliveryAttempts(addressName, DEFAULT_COUNT).
		WithMaxRedeliveryDelay(addressName, -1).
		WithMaxSizeBytes(addressName, DEFAULT_SIZE_SMALL).
		WithMaxSizeBytesRejectThreshold(addressName, 10000).
		WithMessageCounterHistoryDayLimit(addressName, 0).
		WithMinExpiryDelay(addressName, DEFAULT_DELAY).
		WithPageMaxCacheSize(addressName, 20000000).
		WithPageSizeBytes(addressName, "10485760").
		WithRedeliveryCollisionsAvoidance(addressName, "0").
		WithRedeliveryDelay(addressName, 0).
		WithRedeliveryDelayMult(addressName, "1").
		WithRedistributionDelay(addressName, -1).
		WithSendToDLAOnNoRoute(addressName, false).
		WithSlowConsumerCheckPeriod(addressName, 5).
		WithSlowConsumerPolicy(addressName, Notify).
		WithSlowConsumerThreshold(addressName, -1)

	return bdw
}

const (
	DEFAULT_DELAY          = 1000
	DEFAULT_COUNT          = 100
	DEFAULT_SIZE_BIG       = "2G"
	DEFAULT_SIZE_SMALL     = "10K"
	DEFAULT_DEAD_ADDRESS   = ""
	DEFAULT_EXPIRY_ADDRESS = ""
	DEFAULT_KEY            = ""
	DEFAULT_SUFFIX         = ""
	DEFAULT_PREFIX         = ""
	DEFAULT_PERIOD         = DEFAULT_DELAY
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
