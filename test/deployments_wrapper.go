package test

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/onsi/gomega"
	brokerapi "github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"github.com/rh-messaging/shipshape/pkg/framework/operators"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"time"
)

// DeploymentWrapper takes care of deployment of Broker
type DeploymentWrapper struct {
	wait         bool
	brokerClient brokerclientset.Interface
	ctx1         *framework.ContextData
	customImage  string
	migration    bool
	persistence  bool
	name         string
	sslEnabled   bool
}

const (
//AmqpAcceptor =
)

// WithWait sets if shipshape would wait for completion
func (dw *DeploymentWrapper) WithWait(wait bool) *DeploymentWrapper {
	dw.wait = wait
	return dw
}

func (dw *DeploymentWrapper) WithName(name string) *DeploymentWrapper {
	dw.name = name
	return dw
}

// WithBrokerClient sets broker kubernetes client to use
func (dw *DeploymentWrapper) WithBrokerClient(brokerClient brokerclientset.Interface) *DeploymentWrapper {
	dw.brokerClient = brokerClient
	return dw
}

// WithContext sets shipshape context
func (dw *DeploymentWrapper) WithContext(ctx1 *framework.ContextData) *DeploymentWrapper {
	dw.ctx1 = ctx1
	return dw
}

// WithCustomImage wets Broker Image to be used
func (dw *DeploymentWrapper) WithCustomImage(image string) *DeploymentWrapper {
	dw.customImage = image
	return dw
}

// WithMigration Sets Migration parameter (controls message migration availability)
func (dw *DeploymentWrapper) WithMigration(migration bool) *DeploymentWrapper {
	dw.migration = migration
	return dw
}

// WithPersistence Sets Persistence parameter (controls creationf of PVCs)
func (dw *DeploymentWrapper) WithPersistence(persistence bool) *DeploymentWrapper {
	dw.persistence = persistence
	return dw
}

func (dw *DeploymentWrapper) WithSsl(ssl bool) *DeploymentWrapper {
	dw.sslEnabled = ssl
	return dw
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
	acceptorPorts = map[AcceptorType]int32{
		AmqpAcceptor:     5672,
		OpenwireAcceptor: 61616,
		CoreAcceptor:     61613,
	}
	acceptors = map[AcceptorType]*brokerapi.AcceptorType{
		AmqpAcceptor:     defaultAcceptor("amqp", acceptorPorts[AmqpAcceptor]),
		OpenwireAcceptor: defaultAcceptor("openwire", acceptorPorts[OpenwireAcceptor]),
		CoreAcceptor:     defaultAcceptor("core", acceptorPorts[CoreAcceptor]),
		MultiAcceptor:    defaultAcceptor("core,openwire,amqp", acceptorPorts[CoreAcceptor]),
	}
)

// Scale scales already deployed Broker
func (dw *DeploymentWrapper) Scale(result int) error {
	resourceVersion := int64(0)

	var err error
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Get(dw.name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	originalSize := artemisCreated.Spec.DeploymentPlan.Size
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	artemisCreated.Spec.DeploymentPlan.Size = int32(result)
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)
	artemisCreated.Spec.DeploymentPlan.MessageMigration = &dw.migration
	artemisCreated.Spec.DeploymentPlan.PersistenceEnabled = dw.persistence
	artemisCreated.Spec.Console.Expose = true
	_, err = dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	if dw.wait {
		log.Logf("Waiting for exactly %d instances.\n", result)
		err = framework.WaitForStatefulSet(dw.ctx1.Clients.KubeClient,
			dw.ctx1.Namespace,
			dw.name+"-ss",
			result,
			time.Second*10, time.Minute*time.Duration(5*max(result, int(originalSize))))
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		log.Logf("Not waiting for instances to spawn.\n")
	}
	return err
}

func (dw *DeploymentWrapper) DeployBrokersWithAcceptor(count int, acceptorType AcceptorType) error {
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
	artemis.Spec.DeploymentPlan.Size = int32(count)
	artemis.Spec.Acceptors = append(artemis.Spec.Acceptors, *acceptors[acceptorType])
	for num := range artemis.Spec.Acceptors {
		artemis.Spec.Acceptors[num].SSLEnabled = dw.sslEnabled
	}
	for num := range artemis.Spec.Connectors {
		artemis.Spec.Connectors[num].SSLEnabled = dw.sslEnabled
	}

	artemis.Spec.DeploymentPlan.MessageMigration = &dw.migration
	artemis.Spec.DeploymentPlan.PersistenceEnabled = dw.persistence
	artemis.Spec.AdminUser = Username
	artemis.Spec.AdminPassword = Password
	artemis.Spec.DeploymentPlan.Image = dw.customImage
	artemis.ObjectMeta.Name = dw.name

	artemis.Spec.Console.Expose = true

	//dw.ctx1.Clients.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()

	//ctx1.Clients.KubeClient.AppsV1().StatefulSets(ctx1.Namespace).Create(&artemis)
	_, err = dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Create(&artemis)
	gomega.Expect(err).To(gomega.BeNil())

	if dw.wait {
		log.Logf("Waiting for exactly %d instances.\n", count)
		err = framework.WaitForStatefulSet(dw.ctx1.Clients.KubeClient,
			dw.ctx1.Namespace,
			dw.name+"-ss",
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

// DeployBrokers actually deploys brokers defined by dw
func (dw *DeploymentWrapper) DeployBrokers(count int) error {
	return dw.DeployBrokersWithAcceptor(count, AmqpAcceptor)
}

// ChangeImage changes image used in Broker instance to a new one
func (dw *DeploymentWrapper) ChangeImage() error {
	resourceVersion := int64(0)
	var err error
	resourceVersion = resourceVersion + 5
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Get(dw.name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	countExpected := artemisCreated.Spec.DeploymentPlan.Size
	if countExpected == 0 {
		countExpected = 1
	}
	artemisCreated.Spec.DeploymentPlan.Image = dw.customImage
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)
	artemisCreated.Spec.Console.Expose = true

	_, err = dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	err = framework.WaitForStatefulSet(dw.ctx1.Clients.KubeClient,
		dw.ctx1.Namespace,
		dw.name+"-ss",
		int(countExpected),
		time.Second*10, time.Minute*time.Duration(5*countExpected))

	gomega.Expect(err).To(gomega.BeNil())

	return err
}

func PrepareOperator() operators.OperatorSetupBuilder {
	builder := operators.SupportedOperators[operators.OperatorTypeBroker]

	//Set image to parameter if one is supplied, otherwise use default from shipshape.
	if len(Config.OperatorImageName) != 0 {
		builder.WithImage(Config.OperatorImageName)
	}
	if Config.DownstreamBuild {
		builder.WithCommand("/home/amq-broker-operator/bin/entrypoint")
		builder.WithOperatorName("amq-broker-operator")
	}

	if Config.RepositoryPath != "" {
		// Try loading YAMLs from the repo.
		yamls, err := LoadYamls(Config.RepositoryPath)
		if err != nil {
			panic(err)
		} else {
			builder.WithYamls(yamls)
		}
	}

	if Config.AdminAvailable {
		//builder.()
	}
	return builder
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
