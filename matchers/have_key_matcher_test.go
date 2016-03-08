package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/totherme/nosj"
	"github.com/totherme/nosj/matchers"
)

var _ = Describe("HaveJSONKeyMatcher", func() {
	Describe("Match", func() {

		Context("When we give it a JSON object", func() {
			DescribeTable("the matcher matches iff HasKey returns true", func(key string) {
				rawjson := `{"name": "fred",
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

				testjson, err := nosj.ParseJSON(rawjson)
				Expect(err).NotTo(HaveOccurred())

				matcher := matchers.HaveJSONKey(key)

				Expect(matcher.Match(testjson)).To(Equal(testjson.HasKey(key)))
			},

				Entry("a string key", "name"),
				Entry("a number key", "life"),
				Entry("a list key", "othernames"),
				Entry("a boolean key", "beauty"),
				Entry("an object key", "things"),
				Entry("a null key", "not"),
				Entry("an absent key", "badgers"),
			)
		})
		Context("when we give it a non-json object", func() {
			It("returns a helpful error message", func() {
				matcher := matchers.HaveJSONKey("key")
				_, err := matcher.Match(`{"you":"might almost think this would work"}`)
				Expect(err).To(MatchError(ContainSubstring("not a JSON object. Have you done nosj.ParseJSON(...)?")))
			})
		})
	})
})
