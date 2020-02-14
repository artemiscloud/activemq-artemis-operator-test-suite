package basic

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"time"
	//brokerclientset "github.com/rh-messaging/activemq-artemis-operator/pkg/client/clientset/versioned"
	brokerapi "github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v2alpha1"
)

var _ = ginkgo.Describe("DeploymentScalingBroker", func() {

	var (
		ctx1    *framework.ContextData
		artemis brokerapi.ActiveMQArtemis
		resourceVersion int64
		//brokerClient brokerclientset.Interface
	)

	// Initialize after framework has been created
	ginkgo.JustBeforeEach(func() {
		ctx1 = Framework.GetFirstContext()
	})

	ginkgo.It("Deploy single broker instance and scale it to 4 replicas", func() {
		//ctx1.OperatorMap[operators.OperatorTypeBroker].Namespace()
		resp, err := http.Get("https://raw.githubusercontent.com/rh-messaging/activemq-artemis-operator/master/deploy/crs/broker_v2alpha1_activemqartemis_cr.yaml") //load yaml body from url
        resourceVersion=0
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
		//Ignoring artemis result, since it cant' be modified anyways.
		_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Create(&artemis)
		gomega.Expect(err).To(gomega.BeNil())
		err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient,ctx1.Namespace,"ex-aao-ss",1,time.Second*10,time.Minute*5)
		gomega.Expect(err).To(gomega.BeNil())
		resourceVersion = resourceVersion + 5;
		// getting created artemis custom resource to overwrite the resourceVersion and params.
		artemisCreated, err := brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Get("ex-aao",v1.GetOptions{})
		gomega.Expect(err).To(gomega.BeNil())
		resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion),10,64)
		gomega.Expect(err).To(gomega.BeNil())
        artemisCreated.Spec.DeploymentPlan.Size=4
        artemisCreated.ObjectMeta.ResourceVersion=strconv.FormatInt(int64(resourceVersion),10)

        _,  err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Update(artemisCreated)
        gomega.Expect(err).To(gomega.BeNil())
        err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", 4, time.Second*10, time.Minute*5)
        gomega.Expect(err).To(gomega.BeNil())
        
	})
    

})
