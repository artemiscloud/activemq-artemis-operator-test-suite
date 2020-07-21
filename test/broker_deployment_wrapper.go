package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/onsi/gomega"
	brokerapi "github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// BrokerDeploymentWrapper takes care of deployment of Broker
type BrokerDeploymentWrapper struct {
	wait           bool
	brokerClient   brokerclientset.Interface
	ctx1           *framework.ContextData
	customImage    string
	migration      bool
	persistence    bool
	name           string
	sslEnabled     bool
	exposeConsole  bool
	deploymentSize int
}

func defaultAcceptor(protocol string, port int32) *brokerapi.AcceptorType {
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

type AcceptorType int

const (
	AmqpAcceptor AcceptorType = iota
	CoreAcceptor
	OpenwireAcceptor
	MultiAcceptor
)

var (
	AcceptorPorts = map[AcceptorType]int32{
		AmqpAcceptor:     5672,
		OpenwireAcceptor: 61613,
		CoreAcceptor:     61616,
	}
	acceptors = map[AcceptorType]*brokerapi.AcceptorType{
		AmqpAcceptor:     defaultAcceptor(AMQP, AcceptorPorts[AmqpAcceptor]),
		OpenwireAcceptor: defaultAcceptor(OPENWIRE, AcceptorPorts[OpenwireAcceptor]),
		CoreAcceptor:     defaultAcceptor(CORE, AcceptorPorts[CoreAcceptor]),
		MultiAcceptor:    defaultAcceptor(fmt.Sprintf("%s,%s,%s", AMQP, OPENWIRE, CORE), AcceptorPorts[CoreAcceptor]),
	}
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

// WithWait sets if shipshape would wait for completion
func (bdw *BrokerDeploymentWrapper) WithWait(wait bool) *BrokerDeploymentWrapper {
	bdw.wait = wait
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithName(name string) *BrokerDeploymentWrapper {
	bdw.name = name
	return bdw
}

// WithBrokerClient sets broker kubernetes client to use
func (bdw *BrokerDeploymentWrapper) WithBrokerClient(brokerClient brokerclientset.Interface) *BrokerDeploymentWrapper {
	bdw.brokerClient = brokerClient
	return bdw
}

// WithContext sets shipshape context
func (bdw *BrokerDeploymentWrapper) WithContext(ctx1 *framework.ContextData) *BrokerDeploymentWrapper {
	bdw.ctx1 = ctx1
	return bdw
}

// WithCustomImage wets Broker Image to be used
func (bdw *BrokerDeploymentWrapper) WithCustomImage(image string) *BrokerDeploymentWrapper {
	bdw.customImage = image
	return bdw
}

// WithMigration Sets Migration parameter (controls message migration availability)
func (bdw *BrokerDeploymentWrapper) WithMigration(migration bool) *BrokerDeploymentWrapper {
	bdw.migration = migration
	return bdw
}

// WithPersistence Sets Persistence parameter (controls creationf of PVCs)
func (bdw *BrokerDeploymentWrapper) WithPersistence(persistence bool) *BrokerDeploymentWrapper {
	bdw.persistence = persistence
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithSsl(ssl bool) *BrokerDeploymentWrapper {
	bdw.sslEnabled = ssl
	return bdw
}

func getAcceptor(
	name string,
	port int32,
	protocols string,
	sslEnabled bool,
	sslSecret string,
	enabledCipherSuites string,
	enabledProtocols string,
	needsClientAuth bool,
	wantClientAuth bool,
	verifyHost bool,
	sslProvider string,
	sniHost string,
	expose bool,
	anycastPrefix string,
	multicastPrefix string,
	connectionsAllowed int) *brokerapi.AcceptorType {
	acceptor := &brokerapi.AcceptorType{
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

// Scale scales already deployed Broker
func (bdw *BrokerDeploymentWrapper) Scale(result int) error {
	bdw.deploymentSize = result
	return bdw.Update()
}

func (bdw *BrokerDeploymentWrapper) WaitForBrokerSet(result int, originalSize int) {
	err := framework.WaitForStatefulSet(bdw.ctx1.Clients.KubeClient,
		bdw.ctx1.Namespace,
		bdw.name+"-ss",
		result,
		time.Second*10, time.Minute*time.Duration(5*max(result, originalSize)))
	gomega.Expect(err).To(gomega.BeNil())
}

func (bdw *BrokerDeploymentWrapper) DeployBrokersWithAcceptor(count int, acceptorType AcceptorType) error {
	bdw.deploymentSize = count
	artemis := brokerapi.ActiveMQArtemis{}
	resp, err := http.Get("https://raw.githubusercontent.com/rh-messaging/activemq-artemis-operator/master/deploy/crs/broker_v2alpha1_activemqartemis_cr.yaml") //load yaml body from url
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	jsonBody, err := yaml.YAMLToJSON(body)
	err = json.Unmarshal(jsonBody, &artemis)
	if err != nil {
		panic(err)
	}

	log.Logf("modifying acceptors")
	artemis.Spec.DeploymentPlan.Size = int32(bdw.deploymentSize)
	artemis.Spec.Acceptors = append(artemis.Spec.Acceptors, *acceptors[acceptorType])
	for num := range artemis.Spec.Acceptors {
		artemis.Spec.Acceptors[num].SSLEnabled = bdw.sslEnabled
	}
	for num := range artemis.Spec.Connectors {
		artemis.Spec.Connectors[num].SSLEnabled = bdw.sslEnabled
	}
	artemis.Spec.DeploymentPlan.MessageMigration = &bdw.migration
	artemis.Spec.DeploymentPlan.PersistenceEnabled = bdw.persistence
	artemis.Spec.AdminUser = Username
	artemis.Spec.AdminPassword = Password
	artemis.Spec.DeploymentPlan.Image = bdw.customImage
	artemis.ObjectMeta.Name = bdw.name
	artemis.Spec.Console.Expose = bdw.exposeConsole
	//bdw.ctx1.Clients.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	//ctx1.Clients.KubeClient.AppsV1().StatefulSets(ctx1.Namespace).Create(&artemis)
	_, err = bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Create(&artemis)
	gomega.Expect(err).To(gomega.BeNil())
	if bdw.wait {
		log.Logf("Waiting for exactly %d instances.\n", count)
		err = framework.WaitForStatefulSet(bdw.ctx1.Clients.KubeClient,
			bdw.ctx1.Namespace,
			bdw.name+"-ss",
			count,
			time.Second*10, time.Minute*time.Duration(5*count))
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		log.Logf("Not waiting for instances to spawn.\n")
	}
	fmt.Print("Waiting for 5 seconds\n")
	time.Sleep(time.Duration(5) * time.Second)
	fmt.Print("Done waiting\n")
	return err
}

// DeployBrokers actually deploys brokers defined by bdw
func (bdw *BrokerDeploymentWrapper) DeployBrokers(count int) error {
	return bdw.DeployBrokersWithAcceptor(count, AmqpAcceptor)
}

func (bdw *BrokerDeploymentWrapper) VerifyImage(target string) error {
	artemisCreated, err := bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Get(bdw.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	gomega.Expect(artemisCreated.Spec.DeploymentPlan.Image).To(gomega.Equal(target))
	return nil
}

// ChangeImage changes image used in Broker instance to a new one
func (bdw *BrokerDeploymentWrapper) ChangeImage() error {
	return bdw.Update()
}

func (bdw *BrokerDeploymentWrapper) SetConsoleExposure(expose bool) {
	bdw.exposeConsole = expose
}

func (bdw *BrokerDeploymentWrapper) Update() error {
	resourceVersion := int64(0)

	var err error
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Get(bdw.name, v1.GetOptions{})

	gomega.Expect(err).To(gomega.BeNil())
	originalSize := artemisCreated.Spec.DeploymentPlan.Size
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	artemisCreated.Spec.DeploymentPlan.Size = int32(bdw.deploymentSize)
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)
	artemisCreated.Spec.DeploymentPlan.MessageMigration = &bdw.migration
	artemisCreated.Spec.DeploymentPlan.PersistenceEnabled = bdw.persistence
	artemisCreated.Spec.Console.Expose = bdw.exposeConsole
	_, err = bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	if bdw.wait {
		log.Logf("Waiting for exactly %d instances.\n", bdw.deploymentSize)
		bdw.WaitForBrokerSet(bdw.deploymentSize, int(originalSize))
	} else {
		log.Logf("Not waiting for instances to spawn.\n")
	}
	return err
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

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
