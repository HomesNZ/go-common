package geo_test

import (
	"encoding/json"

	. "github.com/HomesNZ/go-common/geo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Point", func() {
	Describe("MarshalJSON", func() {
		It("marshals to null", func() {
			p := &Point{}
			b, err := json.Marshal(p)
			Expect(err).ToNot(HaveOccurred())
			Expect(b).To(Equal([]byte(`null`)))
		})
		It("marshals a value", func() {
			p := NewPoint(1, 2, WGS84)
			b, err := json.Marshal(p)
			Expect(err).ToNot(HaveOccurred())
			Expect(b).To(Equal([]byte(`{"lat":2,"long":1}`)))
		})
	})

	Describe("UnmarshalJSON", func() {
		It("unmarshals null", func() {
			var p Point
			b := []byte(`null`)
			err := json.Unmarshal(b, &p)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.IsNull()).To(BeTrue())
		})
		It("unmarshals a value", func() {
			var p Point
			b := []byte(`{"lat":2,"long":1}`)
			err := json.Unmarshal(b, &p)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.Lat).To(Equal(2.0))
			Expect(p.Long).To(Equal(1.0))
			Expect(p.SRID).To(Equal(WGS84))
		})
	})
})
