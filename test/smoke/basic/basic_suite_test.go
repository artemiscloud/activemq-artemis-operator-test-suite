package basic

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/onsi/gomega"
	brokerapi "github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/ginkgowrapper"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {

	gomega.RegisterFailHandler(ginkgowrapper.Fail)
	test.Initialize(t, "basic", "Basic Suite")
}

func DeployAndScale(ctx1 *framework.ContextData,
	initial int,
	result int) error {
	resourceVersion := int64(0)
	err := DeployBrokers(ctx1, initial)
	gomega.Expect(err).To(gomega.BeNil())
	err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", 1, time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion = resourceVersion + 5
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Get("ex-aao", v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	artemisCreated.Spec.DeploymentPlan.Size = 1
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)

	_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", 4, time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())

	return err
}

func DeployBrokers(ctx1 *framework.ContextData, count int) error {
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

	for num, _ := range artemis.Spec.Acceptors {
		artemis.Spec.Acceptors[num].SSLEnabled = false
	}
	for num, _ := range artemis.Spec.Connectors {
		artemis.Spec.Connectors[num].SSLEnabled = false
	}
	artemis.Spec.DeploymentPlan.Image = test.BrokerImageName

	ctx1.Clients.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()

	//ctx1.Clients.KubeClient.AppsV1().StatefulSets(ctx1.Namespace).Create(&artemis)
	_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Create(&artemis)
	gomega.Expect(err).To(gomega.BeNil())
	err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", 1, time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())

	return err
}
