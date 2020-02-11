package basic

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"io/ioutil"
	"net/http"
	"time"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	//brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	brokerapi "github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v2alpha1"


)

var _ = ginkgo.Describe("DeploymentSingleBroker", func() {

	var (
		ctx1    *framework.ContextData
		artemis brokerapi.ActiveMQArtemis
		//brokerClient brokerclientset.Interface
	)

	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
	})

	ginkgo.It("Deploy single broker instance", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
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
		artemis.Spec.DeploymentPlan.Size=1

		for num,_ := range artemis.Spec.Acceptors {
			artemis.Spec.Acceptors[num].SSLEnabled=false
		}
		for num,_ := range artemis.Spec.Connectors {
			artemis.Spec.Connectors[num].SSLEnabled = false
		}
		artemis.Spec.DeploymentPlan.Image=test.BrokerImageName
		
		ctx1.Clients.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()

		//ctx1.Clients.KubeClient.AppsV1().StatefulSets(ctx1.Namespace).Create(&artemis)
		_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Create(&artemis)
		gomega.Expect(err).To(gomega.BeNil())
		err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient,ctx1.Namespace,"ex-aao-ss",1,time.Second*10,time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
	})

})
