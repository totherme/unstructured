package matchers_test

import (
	"github.com/totherme/nosj"
	"github.com/totherme/nosj/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeAString", func() {
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
		DescribeTable("BeAString matches iff IsString returns true", func(key string) {
			testjson, err := nosj.ParseJSON(rawjson)
			Expect(err).NotTo(HaveOccurred())

			field := testjson.F(key)
			Expect(matchers.BeAString().Match(field)).To(Equal(field.IsString()))
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
			_, err := matchers.BeAString().Match(4)
			Expect(err).To(MatchError(ContainSubstring("not a JSON")))
		})
	})

	Describe("FailureMessage", func() {
		Context("when we get a JSON struct", func() {
			DescribeTable("it gives the /actual/ json type we're looking at", func(key string, typ string) {
				testjson, err := nosj.ParseJSON(rawjson)
				Expect(err).NotTo(HaveOccurred())

				field := testjson.F(key)

				Expect(matchers.BeAString().FailureMessage(field)).To(ContainSubstring(typ))
			},
				Entry("an object", "things", "object"),
				Entry("a number", "life", "number"),
				Entry("a list", "othernames", "list"),
				Entry("a boolean", "beauty", "bool"),
				Entry("null", "not", "null"),
			)
		})

		Context("when we get some other type of struct", func() {
			It("mentions the type of the struct we /did/ get", func() {
				Expect(matchers.BeAString().FailureMessage(12)).To(ContainSubstring("int"))
			})
		})
	})

	Describe("NegatedFailureMessage", func() {
		It("tells us we got a JSON string", func() {
			json, err := nosj.ParseJSON(`{}`)
			Expect(err).NotTo(HaveOccurred())
			Expect(matchers.BeAString().NegatedFailureMessage(json)).To(ContainSubstring("got a JSON string"))
		})
	})
})
