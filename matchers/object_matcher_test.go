package matchers_test

import (
	"github.com/totherme/nosj"
	"github.com/totherme/nosj/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeAnObject", func() {
	var rawjson string

	BeforeEach(func() {
		rawjson = `{"name": "fred",
							"othernames": [
								"alice",
								"bob",
								"ezekiel"
							],
							"life": 42,
							"things": {
								"more": "things"
							},
							"beauty": true,
							"not": null
						}`
	})

	Context("when we're given a JSON struct", func() {
		DescribeTable("BeAnObject matches iff IsOb returns true", func(key string) {
			testjson, err := nosj.ParseJSON(rawjson)
			Expect(err).NotTo(HaveOccurred())

			field := testjson.F(key)
			Expect(matchers.BeAnObject().Match(field)).To(Equal(field.IsOb()))
		},

			Entry("a string", "name"),
			Entry("a number", "life"),
			Entry("a list", "othernames"),
			Entry("a boolean", "beauty"),
			Entry("an object", "things"),
			Entry("null", "not"),
		)
	})

	Context("when we're not given a json struct", func() {
		It("fails", func() {
			_, err := matchers.BeAnObject().Match(4)
			Expect(err).To(MatchError(ContainSubstring("not a JSON")))
		})
	})

	Describe("FailureMessage", func() {
		Context("when we get a JSON object", func() {
			DescribeTable("it gives the /actual/ json type we're looking at", func(key string, typ string) {
				testjson, err := nosj.ParseJSON(rawjson)
				Expect(err).NotTo(HaveOccurred())

				field := testjson.F(key)

				Expect(matchers.BeAnObject().FailureMessage(field)).To(ContainSubstring(typ))
			},
				Entry("a string", "name", "string"),
				Entry("a number", "life", "number"),
				Entry("a list", "othernames", "list"),
				Entry("a boolean", "beauty", "bool"),
				Entry("null", "not", "null"),
			)
		})

		Context("when we get some other type of object", func() {
			It("mentions the type of the object we /did/ get", func() {
				Expect(matchers.BeAnObject().FailureMessage(12)).To(ContainSubstring("int"))
			})
		})
	})

	Describe("NegatedFailureMessage", func() {
		It("tells us we got a JSON object", func() {
			json, err := nosj.ParseJSON(`{}`)
			Expect(err).NotTo(HaveOccurred())
			Expect(matchers.BeAnObject().NegatedFailureMessage(json)).To(ContainSubstring("got a JSON object"))
		})
	})
})
