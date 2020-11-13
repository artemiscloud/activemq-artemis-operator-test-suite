package statistics

import (
	"github.com/onsi/ginkgo"
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "configuration"
	BaseName   = "broker-framework"
)

var (
	sw *test.SetupWrapper
)

// This needs to be configured on per-test basis
var _ = ginkgo.BeforeEach(func() {
	sw = &test.SetupWrapper{}
	sw.WithBaseName(BaseName).WithDeployName(DeployName)
	sw.BeforeEach()
}, 60)

var _ = ginkgo.JustBeforeEach(func() {
	sw.JustBeforeEach()
})

// After each test completes, run cleanup actions to save resources (otherwise resources will remain till
// all specs from this suite are done.
var _ = ginkgo.AfterEach(func() {
	sw.AfterEach()
})
