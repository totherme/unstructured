package matchers_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/totherme/nosj"
	"github.com/totherme/nosj/matchers"
)

var _ = Describe("HaveJSONKeyMatcher", func() {
	var json nosj.JSON
	BeforeEach(func() {
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
		var err error
		json, err = nosj.ParseJSON(rawjson)
		Expect(err).NotTo(HaveOccurred())

	})
	Describe("Match", func() {
		Context("When we give it a JSON object", func() {
			DescribeTable("the matcher matches iff HasKey returns true", func(key string) {

				matcher := matchers.HaveJSONKey(key)
				Expect(matcher.Match(json)).To(Equal(json.HasKey(key)))
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

	Describe("FailureMessage", func() {
		It("should tell us what key we expected to find", func() {
			Expect(matchers.HaveJSONKey("my-key").FailureMessage("actual-object")).
				To(ContainSubstring("expected 'actual-object' to be a nosj.JSON object with key 'my-key'"))
			Expect(matchers.HaveJSONKey("my-key").FailureMessage(42)).
				To(ContainSubstring("expected '42' to be a nosj.JSON object with key 'my-key'"))
		})

		Context("when the input is large", func() {
			It("truncates the string representation", func() {
				Expect(matchers.HaveJSONKey("absent-key").FailureMessage(json)).
					To(ContainSubstring("'{nosj:map[beauty:true life:42 name:fred othernames...'"))
			})
		})

		Context("when the input is exactly as large as we're willing to print", func() {
			It("prints it all, without elipses", func() {
				stringOfLength50 := strings.Repeat("a", 50)
				failureMessage := matchers.HaveJSONKey("absent-key").FailureMessage(stringOfLength50)
				Expect(failureMessage).To(ContainSubstring(fmt.Sprintf("'%s'", stringOfLength50)))
				Expect(failureMessage).NotTo(ContainSubstring("..."))
			})
		})
	})

	Describe("NegatedFailureMessage", func() {
		var (
			shortJson      nosj.JSON
			jsonOfLength50 nosj.JSON
		)

		BeforeEach(func() {
			var err error
			shortJson, err = nosj.ParseJSON(`{"key":"val"}`)
			Expect(err).NotTo(HaveOccurred())
			stringOfLen34 := strings.Repeat("a", 34)
			jsonOfLength50, err = nosj.ParseJSON(fmt.Sprintf(`{"key":"%s"}`, stringOfLen34))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should tell us what key we expected not to find", func() {
			Expect(matchers.HaveJSONKey("key").NegatedFailureMessage(shortJson)).
				To(ContainSubstring("expected '{nosj:map[key:val]}' not to contain the key 'key'"))
		})

		Context("when the input is large", func() {
			It("truncates the string representation", func() {
				Expect(matchers.HaveJSONKey("beauty").NegatedFailureMessage(json)).
					To(ContainSubstring("'{nosj:map[beauty:true life:42 name:fred othernames...'"))
			})
		})

		Context("when the input is exactly as large as we're willing to print", func() {
			It("prints it all, without elipses", func() {
				failureMessage := matchers.HaveJSONKey("key").NegatedFailureMessage(jsonOfLength50)
				Expect(failureMessage).To(ContainSubstring(fmt.Sprintf("'%+v'", jsonOfLength50)))
				Expect(failureMessage).NotTo(ContainSubstring("..."))
			})
		})
	})
})
