package db_test

import (
	"encoding/json"

	"github.com/HomesNZ/go-common/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PG Array Test", func() {
	Context("converting row value", func() {
		It("returns the correctly formatted array", func() {
			value := `{"one, sdf",string_1," ",test@homes-dev.co.nz," GO IS THE BEST!!! ",abcdef:83bf98cc-fec9-4e77-b4cf-99f9fb6655fa-0NH:zxcvzxc:wers:vxdfw-asdf-asdf,test@nothing.tld,test@nothing2.tld}`
			actual := db.ParseArray(value)
			expected := db.PGArray{"one, sdf", "string_1", " ", "test@homes-dev.co.nz", " GO IS THE BEST!!! ", `abcdef:83bf98cc-fec9-4e77-b4cf-99f9fb6655fa-0NH:zxcvzxc:wers:vxdfw-asdf-asdf`, "test@nothing.tld", "test@nothing2.tld"}

			Expect(actual).To(Equal(expected))

		})
	})
	Context("converting a splice/arry to string for Postgres", func() {
		It("is a string", func() {

			value := db.PGArray{"one, sdf", "{one}", "string_1", " ", "test@homes-dev.co.nz", " GO IS THE BEST!!! ", `abcdef:83bf98cc-fec9-4e77-b4cf-99f9fb6655fa-0NH:zxcvzxc:wers:vxdfw-asdf-asdf`, "test@nothing.tld", "test@nothing2.tld"}

			actual := db.CreateStringArray(value)
			expected := `{"one, sdf","{one}","string_1"," ","test@homes-dev.co.nz"," GO IS THE BEST!!! ","abcdef:83bf98cc-fec9-4e77-b4cf-99f9fb6655fa-0NH:zxcvzxc:wers:vxdfw-asdf-asdf","test@nothing.tld","test@nothing2.tld"}`

			Expect(actual).To(Equal(expected))
		})
	})

	It("marshals to JSON as an array of strings", func() {
		a := db.PGArray{"a", "b", "c"}
		b, err := json.Marshal(a)
		Expect(err).ToNot(HaveOccurred())
		Expect(b).To(Equal([]byte(`["a","b","c"]`)))
	})
	It("unmarshals an array of strings from JSON", func() {
		j := []byte(`["a","b","c"]`)
		a := db.PGArray{}
		err := json.Unmarshal(j, &a)
		Expect(err).ToNot(HaveOccurred())
		Expect(a).To(Equal(db.PGArray{"a", "b", "c"}))
	})
	It("unmarshals an empty array from JSON", func() {
		j := []byte(`[]`)
		a := db.PGArray{}
		err := json.Unmarshal(j, &a)
		Expect(err).ToNot(HaveOccurred())
		Expect(a).To(Equal(db.PGArray{}))
	})
	It("unmarshals a null value from JSON", func() {
		j := []byte(`null`)
		a := db.PGArray{}
		var b db.PGArray
		err := json.Unmarshal(j, &a)
		Expect(err).ToNot(HaveOccurred())
		Expect(a).To(Equal(b))
	})
})
