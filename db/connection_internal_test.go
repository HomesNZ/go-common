package db

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PG", func() {
	Describe("logSafeConnectionString", func() {
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_USER", "user")
		os.Setenv("DB_HOST", "host")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_PASSWORD", "password")
		os.Setenv("DB_SSL_MODE", "disable")
		os.Setenv("DB_SEARCH_PATH", "public")

		Context("when there is a search path", func() {
			It("adds the search_path parameter", func() {
				db := PG{}

				actual := db.logSafeConnectionString("go-common")
				expected := "postgres://user:****@host:5432/test?sslmode=disable&application_name=go-common&search_path=public"

				Expect(actual).To(Equal(expected))
			})
		})

		Context("when there is no search path", func() {
			It("leaves the connection string as is", func() {
				os.Setenv("DB_SEARCH_PATH", "")

				db := PG{}

				actual := db.logSafeConnectionString("go-common")
				expected := "postgres://user:****@host:5432/test?sslmode=disable&application_name=go-common"

				Expect(actual).To(Equal(expected))
			})
		})

		Context("when there is a password", func() {
			It("replaces the password with stars", func() {
				db := PG{}

				actual := db.logSafeConnectionString("go-common")
				expected := "postgres://user:****@host:5432/test?sslmode=disable&application_name=go-common"

				Expect(actual).To(Equal(expected))
			})
		})

		Context("when there is no password", func() {
			It("leaves the connection string as is", func() {
				os.Setenv("DB_PASSWORD", "")

				db := PG{}

				actual := db.logSafeConnectionString("go-common")
				expected := db.connectionString("go-common")

				Expect(actual).To(Equal(expected))
			})
		})
	})
})
