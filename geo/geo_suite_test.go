package geo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Geo Suite")
}
