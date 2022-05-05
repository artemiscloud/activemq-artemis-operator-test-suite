package resources

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
)

// Constants available for all test specs related with the One Interior topology.

const (
	DeployName = "resources"
	BaseName   = "broker"
)

var (
	sw *test.SetupWrapper
)

var _ = ginkgo.BeforeEach(func() {
	sw = &test.SetupWrapper{}
	sw.WithBaseName(BaseName).WithDeployName(DeployName)
	if !test.Config.NeedsLatestCR {
		ginkgo.Skip("Not supported on pre-0.17 operator")
	}
	sw.BeforeEach()
}, 60)
