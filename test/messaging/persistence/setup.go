package persistence

import (
	"github.com/artemiscloud/activemq-artemis-operator-test-suite/test"
	"github.com/onsi/ginkgo"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "persistence"
	BaseName   = "broker-framework"
)

var (
	sw *test.SetupWrapper
)

// Create the Framework instance to be used oneinterior test

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
