package basic_test

import (
	"broker-suite/test"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBasic(t *testing.T) {

	RegisterFailHandler(Fail)
	test.Initialize(t, "basic", "Basic Suite")

	RunSpecs(t, "Basic Suite")

}
