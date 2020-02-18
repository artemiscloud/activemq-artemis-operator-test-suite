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
//	corev1 "k8s.io/api/core/v1"
//	appsv1 "k8s.io/api/apps/v1"
	"net/http"
	"strconv"
	"time"
)

func Scale(ctx1 *framework.ContextData,
	result int,
	brokerClient brokerclientset.Interface) error {
	resourceVersion := int64(0)
	var err error
	resourceVersion = resourceVersion + 5
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Get("ex-aao", v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	artemisCreated.Spec.DeploymentPlan.Size = int32(result)
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)

	_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	log.Logf("Waiting for exactly " + string(result) + " instances.")
	err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", result, time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())

	return err
}

func DeployBrokers(ctx1 *framework.ContextData,
	count int,
	brokerClient brokerclientset.Interface,
	brokerImage string) error {
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
	artemis.Spec.AdminUser = Username
	artemis.Spec.AdminPassword = Password
	artemis.Spec.DeploymentPlan.Image = brokerImage

	ctx1.Clients.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()

	//ctx1.Clients.KubeClient.AppsV1().StatefulSets(ctx1.Namespace).Create(&artemis)
	_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Create(&artemis)
	gomega.Expect(err).To(gomega.BeNil())

	err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", count, time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())

	fmt.Print("Waiting for 5 seconds\n")
	time.Sleep(time.Duration(5) * time.Second)
	fmt.Print("Done waiting\n")
	return err
}

func ChangeImage(ctx1 *framework.ContextData,
	brokerClient brokerclientset.Interface, newImage string) error {
	resourceVersion := int64(0)
	var err error
	resourceVersion = resourceVersion + 5
	// getting created artemis custom resource to overwrite the resourceVersion and params.
	artemisCreated, err := brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Get("ex-aao", v1.GetOptions{})
	gomega.Expect(err).To(gomega.BeNil())
	resourceVersion, err = strconv.ParseInt(string(artemisCreated.ObjectMeta.ResourceVersion), 10, 64)
	gomega.Expect(err).To(gomega.BeNil())
	countExpected := artemisCreated.Spec.DeploymentPlan.Size
	artemisCreated.Spec.DeploymentPlan.Image = newImage
	artemisCreated.ObjectMeta.ResourceVersion = strconv.FormatInt(int64(resourceVersion), 10)

	_, err = brokerClient.BrokerV2alpha1().ActiveMQArtemises(ctx1.Namespace).Update(artemisCreated)
	gomega.Expect(err).To(gomega.BeNil())
	err = framework.WaitForStatefulSet(ctx1.Clients.KubeClient, ctx1.Namespace, "ex-aao-ss", int(countExpected), time.Second*10, time.Minute*5)
	gomega.Expect(err).To(gomega.BeNil())

	return err
}