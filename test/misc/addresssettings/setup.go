package addresssettings

import (
	"github.com/onsi/ginkgo"
	"gitlab.cee.redhat.com/msgqe/openshift-broker-suite-golang/test"
)

// Constants available for all test specs related with the One Interior topology
const (
	DeployName = "addresssetting"
	BaseName   = "broker-framework"
)

var (
	sw *test.SetupWrapper
)

// Create the Framework instance to be used oneinterior test

var _ = ginkgo.BeforeEach(func() {
	sw = &test.SetupWrapper{}
	sw.WithBaseName(BaseName).WithDeployName(DeployName)
	if !test.Config.NeedsLatestCR {
		ginkgo.Skip("Not supported on pre-0.17 operator")
	}
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
