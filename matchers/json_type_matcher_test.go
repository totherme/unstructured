package matchers_test

import (
	"github.com/totherme/nosj"
	"github.com/totherme/nosj/matchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type testData struct {
	Matcher types.GomegaMatcher
	Typ     string
}

var _ = Describe("The JSON type matchers", func() {
	var rawjson string
	var matcherSet []testData

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

		matcherSet = []testData{
			{
				Matcher: matchers.BeAnObject(),
				Typ:     nosj.JSONOb,
			},
			{
				Matcher: matchers.BeAString(),
				Typ:     nosj.JSONString,
			},
			{
				Matcher: matchers.BeAList(),
				Typ:     nosj.JSONList,
			},
			{
				Matcher: matchers.BeANum(),
				Typ:     nosj.JSONNum,
			},
			{
				Matcher: matchers.BeABool(),
				Typ:     nosj.JSONBool,
			},
			{
				Matcher: matchers.BeANull(),
				Typ:     nosj.JSONNull,
			},
		}
	})

	Context("when we're given a JSON struct", func() {
		DescribeTable("each matcher matches iff IsOfType returns true for its type", func(key string) {
			testjson, err := nosj.ParseJSON(rawjson)
			Expect(err).NotTo(HaveOccurred())

			field := testjson.F(key)
			for _, td := range matcherSet {
				Expect(td.Matcher.Match(field)).To(Equal(field.IsOfType(td.Typ)))
			}
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
			for _, td := range matcherSet {
				_, err := td.Matcher.Match(4)
				Expect(err).To(MatchError(ContainSubstring("not a JSON")))
			}
		})
	})

	Describe("FailureMessage", func() {
		Context("when we get a JSON struct", func() {
			DescribeTable("it gives the /actual/ json type we're looking at", func(key string, typ string) {
				testjson, err := nosj.ParseJSON(rawjson)
				Expect(err).NotTo(HaveOccurred())

				field := testjson.F(key)

				for _, td := range matcherSet {
					Expect(td.Matcher.FailureMessage(field)).To(ContainSubstring("expected a JSON %s", td.Typ))
					Expect(td.Matcher.FailureMessage(field)).To(ContainSubstring("got a JSON %s", typ))
				}
			},
				Entry("an object", "things", "object"),
				Entry("a string", "name", "string"),
				Entry("a number", "life", "number"),
				Entry("a list", "othernames", "list"),
				Entry("a boolean", "beauty", "bool"),
				Entry("null", "not", "null"),
			)
		})

		Context("when we get some other type of struct", func() {
			It("mentions the type of the struct we /did/ get", func() {
				for _, td := range matcherSet {
					Expect(td.Matcher.FailureMessage(12)).To(ContainSubstring("int"))
				}
			})
		})
	})

	Describe("NegatedFailureMessage", func() {
		It("tells us we got a JSON object", func() {
			json, err := nosj.ParseJSON(`{}`)
			Expect(err).NotTo(HaveOccurred())
			for _, td := range matcherSet {
				Expect(td.Matcher.NegatedFailureMessage(json)).To(ContainSubstring("got a JSON %s", td.Typ))
			}
		})
	})
})
