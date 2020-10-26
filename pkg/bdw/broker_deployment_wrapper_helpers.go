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
	if acceptorType!= NoChangeAcceptor {
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

	return artemis
}

func (bdw *BrokerDeploymentWrapper) ConvertToV1(artemisOriginal *brokerv3.ActiveMQArtemis) *brokerv1.ActiveMQArtemis {
	artemisResult := &brokerv1.ActiveMQArtemis{}
	data, err := json.Marshal(artemisOriginal)
	if err!= nil {
		panic(err)
	}
	err = json.Unmarshal(data, artemisResult)
	if err!= nil {
		panic(err)
	}
	return artemisResult
}
