package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEcoscript(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ecoscript suite")
}
