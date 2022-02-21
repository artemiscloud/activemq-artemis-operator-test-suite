package test

import (
	"index/suffixarray"
	"regexp"
	"strings"
	"time"

	"github.com/onsi/gomega"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
)

func FormUrl(protocol, DeployName, number, subdomain, namespace, domain, address, port string) string {
	return protocol + "://" + DeployName + "-ss-" + number + "." + DeployName + subdomain + "." + namespace + "." + domain + ":" + port +
		"/" + address
}

func WaitForDrainerRemovalSlow(sw *SetupWrapper, count int, timeout time.Duration, retries int) bool {
	expectedLog := "Deleting drain pod"
	loop := 0
	r := regexp.MustCompile(expectedLog)
	label := Config.BrokerName + "-operator"
	operatorPodName, err := sw.Framework.GetFirstContext().GetPodName(label)
	log.Logf("loading logs from pod %s", operatorPodName)
	gomega.Expect(err).To(gomega.BeNil())
	for loop < retries {
		if loop%10 == 0 {
			log.Logf("(still) waiting for drainer completion")
		}
		operatorLog, _ := sw.Framework.GetFirstContext().GetLogs(operatorPodName)
		if strings.Contains(operatorLog, expectedLog) {
			index := suffixarray.New([]byte(operatorLog))
			results := index.FindAllIndex(r, -1)
			if len(results) == count {
				return true
			}
		}
		time.Sleep(timeout)
		loop++
	}
	return false
}

// WaitForDrainerRemoval would check logs for amount of drainer finished messages.
// Wait for up to 60 seconds * count in intervals of 10 seconds
// Returns true when found all drainers expected, false otherwise
func WaitForDrainerRemoval(sw *SetupWrapper, count int) bool {
	return WaitForDrainerRemovalSlow(sw, count, time.Second*time.Duration(10), count*6)
}
