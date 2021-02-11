package bdw

/* This file contains deployment-related helper methods for BrokerDeploymentWrapper
 */

import (
	"encoding/json"
	"fmt"
	brokerv3 "github.com/artemiscloud/activemq-artemis-operator/pkg/apis/broker/v2alpha3"
	"github.com/fgiorgetti/qpid-dispatch-go-tests/pkg/framework/log"
	"github.com/ghodss/yaml"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"time"
)

func (bdw *BrokerDeploymentWrapper) DeployBrokersWithAcceptor(count int, acceptorType AcceptorType) error {
	bdw.deploymentSize = count
	artemis := &brokerv3.ActiveMQArtemis{}
	resp, err := http.Get("https://raw.githubusercontent.com/activemq-artemis-operator/blob/master/deploy/crs/broker_activemqartemis_cr.yaml") //load yaml body from url
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
	artemis = bdw.ConfigureBroker(artemis, acceptorType)
	return bdw.CreateBroker(artemis, count)
}

func (bdw *BrokerDeploymentWrapper) CreateBroker(artemis *brokerv3.ActiveMQArtemis, count int) error {
	var err error
    log.Logf("Timeout: %s", bdw.GetTimeout(count))
	if bdw.isLtsDeployment {
		artemisConverted := bdw.ConvertToV1(artemis)
		_, err = bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Create(artemisConverted)
	} else {
		_, err = bdw.brokerClient.BrokerV2alpha3().ActiveMQArtemises(bdw.ctx1.Namespace).Create(artemis)
	}

	gomega.Expect(err).To(gomega.BeNil())
	if bdw.wait {
		log.Logf("Waiting for exactly %d instances.\n", count)
		err = framework.WaitForStatefulSet(bdw.ctx1.Clients.KubeClient,
			bdw.ctx1.Namespace,
			bdw.name+"-ss",
			count,
			time.Second*10,bdw.GetTimeout(count))
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

func (bdw *BrokerDeploymentWrapper) Update() error {
	resourceVersion := int64(0)

	var err error
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := bdw.brokerClient.BrokerV2alpha3().ActiveMQArtemises(bdw.ctx1.Namespace).Get(bdw.name, v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	originalSize := artemisCreated.Spec.DeploymentPlan.Size
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)

	bdw.ConfigureBroker(artemisCreated, NoChangeAcceptor)

	if bdw.isLtsDeployment {
		artemisConverted := bdw.ConvertToV1(artemisCreated)
		_, err = bdw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(bdw.ctx1.Namespace).Update(artemisConverted)

	} else {
		_, err = bdw.brokerClient.BrokerV2alpha3().ActiveMQArtemises(bdw.ctx1.Namespace).Update(artemisCreated)
	}
	gomega.Expect(err).To(gomega.BeNil())
	if bdw.wait {
		log.Logf("Waiting for exactly %d instances.\n", bdw.deploymentSize)
		bdw.WaitForBrokerSet(bdw.deploymentSize, int(originalSize))
	} else {
		log.Logf("Not waiting for instances to spawn.\n")
	}
	return err
}

// Scale scales already deployed Broker
func (bdw *BrokerDeploymentWrapper) Scale(result int) error {
	bdw.deploymentSize = result
	return bdw.Update()
}

// ChangeImage changes image used in Broker instance to a new one
func (bdw *BrokerDeploymentWrapper) ChangeImage() error {
	return bdw.Update()
}
