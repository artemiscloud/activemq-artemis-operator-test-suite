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
}

// WithWait sets if shipshape would wait for completion
func (dw DeploymentWrapper) WithWait(wait bool) DeploymentWrapper {
	dw.wait = wait
	return dw
}

// WithBrokerClient sets broker kubernetes client to use
func (dw DeploymentWrapper) WithBrokerClient(brokerClient brokerclientset.Interface) DeploymentWrapper {
	dw.brokerClient = brokerClient
	return dw
}

// WithContext sets shipshape context
func (dw DeploymentWrapper) WithContext(ctx1 *framework.ContextData) DeploymentWrapper {
	dw.ctx1 = ctx1
	return dw
}

// WithCustomImage wets Broker Image to be used
func (dw DeploymentWrapper) WithCustomImage(image string) DeploymentWrapper {
	dw.customImage = image
	return dw
}

// Scale scales already deployed Broker
func (dw DeploymentWrapper) Scale(result int) error {
	resourceVersion := int64(0)
	var err error
	resourceVersion = resourceVersion + 5
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Get("ex-aao", v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	artemisCreated.Spec.DeploymentPlan.Size = int32(result)
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)

	_, err = dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	if dw.wait {
		log.Logf("Waiting for exactly " + string(result) + " instances.\n")
		err = framework.WaitForStatefulSet(dw.ctx1.Clients.KubeClient, dw.ctx1.Namespace, "ex-aao-ss", result, time.Second*10, time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		log.Logf("Not waiting for instances to spawn.\n")
	}
	return err
}

// DeployBrokers actually deploys brokers defined by dw
func (dw DeploymentWrapper) DeployBrokers(count int) error {
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
	_ = json.Unmarshal(jsonBody, &artemis)
	if err != nil {
		panic(err)
	}

	log.Logf("modifying acceptors")
	artemis.Spec.DeploymentPlan.Size = int32(count)

	for num := range artemis.Spec.Acceptors {
		artemis.Spec.Acceptors[num].SSLEnabled = false
	}
	for num := range artemis.Spec.Connectors {
		artemis.Spec.Connectors[num].SSLEnabled = false
	}
	artemis.Spec.AdminUser = Username
	artemis.Spec.AdminPassword = Password
	artemis.Spec.DeploymentPlan.Image = dw.customImage
	artemis.ObjectMeta.Name = "ex-aao"

	//dw.ctx1.Clients.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()

	//ctx1.Clients.KubeClient.AppsV1().StatefulSets(ctx1.Namespace).Create(&artemis)
	_, err = dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Create(&artemis)
	gomega.Expect(err).To(gomega.BeNil())

	if dw.wait {
		log.Logf("Waiting for exactly " + string(count) + " instances.\n")
		err = framework.WaitForStatefulSet(dw.ctx1.Clients.KubeClient, dw.ctx1.Namespace, "ex-aao-ss", count, time.Second*10, time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
	} else {
		log.Logf("Not waiting for instances to spawn.\n")
	}
	fmt.Print("Waiting for 5 seconds\n")
	time.Sleep(time.Duration(5) * time.Second)
	fmt.Print("Done waiting\n")
	return err
}

// ChangeImage changes image used in Broker instance to a new one
func (dw DeploymentWrapper) ChangeImage() error {
	resourceVersion := int64(0)
	var err error
	resourceVersion = resourceVersion + 5
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Get("ex-aao", v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	countExpected := artemisCreated.Spec.DeploymentPlan.Size
	artemisCreated.Spec.DeploymentPlan.Image = dw.customImage
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)

	_, err = dw.brokerClient.BrokerV2alpha1().ActiveMQArtemises(dw.ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	err = framework.WaitForStatefulSet(dw.ctx1.Clients.KubeClient, dw.ctx1.Namespace, "ex-aao-ss", int(countExpected), time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())

	return err
}
