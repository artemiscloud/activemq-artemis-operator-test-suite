package test_helpers

import (
	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CheckVolumeSize(ctx1 *framework.ContextData, expectedSize string) {
	kubeclient := ctx1.Clients.KubeClient
	opts := metav1.ListOptions{}
	pvcList, _ := kubeclient.CoreV1().PersistentVolumeClaims(ctx1.Namespace).List(opts)
	for _, pvc := range pvcList.Items {
		item := pvc.Spec.Resources.Requests["storage"]
		gomega.Expect(item.String()).To(gomega.Equal(expectedSize), "Actual volume size (%s) not equal to expected (%s)", item.String(), expectedSize)
	}
}
