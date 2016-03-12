package matchers_test

import (
	"github.com/totherme/nosj/matchers"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/totherme/nosj"
	"strings"
)

var _ = Describe("HavePointerMatcher", func() {
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
			DescribeTable("the matcher matches iff HasPointer returns true", func(p string) {

				var matcher types.GomegaMatcher = matchers.HaveJSONPointer(p)
				hasp, err := json.HasPointer(p)
				Expect(err).NotTo(HaveOccurred())
				Expect(matcher.Match(json)).To(Equal(hasp))
			},
				Entry("a string pointer", "/name"),
				Entry("a number pointer", "/life"),
				Entry("a list pointer", "/othernames"),
				Entry("a boolean pointer", "/beauty"),
				Entry("an object pointer", "/things"),
				Entry("a pointer to null", "/not"),
				Entry("an absent pointer", "/badgers"),
				Entry("a long pointer", "/things/more"),
			)
		})

		Context("when we give it a non-json object", func() {
			It("returns a helpful error message", func() {
				matcher := matchers.HaveJSONPointer("/perfectly/valid")
				_, err := matcher.Match(`{"you":"might almost think this would work"}`)
				Expect(err).To(MatchError(ContainSubstring("not a JSON object. Have you done nosj.ParseJSON(...)?")))
			})
		})
		Context("when we give it an invalid pointer", func() {
			It("returns a helpful error message", func() {
				matcher := matchers.HaveJSONPointer("not/a/valid/pointer")
				_, err := matcher.Match(json)
				Expect(err).To(MatchError(ContainSubstring("JSON pointer must be empty or start with a \"/\"")))
			})
		})
	})
	Describe("FailureMessage", func() {
		It("should tell us what pointer we expected to find", func() {
			Expect(matchers.HaveJSONPointer("/my/pointer").FailureMessage("actual-object")).
				To(ContainSubstring("expected 'actual-object' to be a nosj.JSON object with pointer '/my/pointer'"))
		})
		Context("when the input has a long string representation", func() {
			It("truncates that representation", func() {
				Expect(len(matchers.HaveJSONPointer("/pointer").FailureMessage(json))).To(BeNumerically("<", 115))
				Expect(matchers.HaveJSONPointer("/pointer").FailureMessage(json)).
					To(ContainSubstring("..."))
				Expect(matchers.HaveJSONPointer("/pointer").FailureMessage(json)).
					To(ContainSubstring("{nosj:map"))
			})
		})

		Context("when the input's string representation is exactly as large as we're willing to print", func() {
			It("prints it all, without elipses", func() {
				stringOfLength50 := strings.Repeat("a", 50)
				failureMessage := matchers.HaveJSONPointer("/pointer").FailureMessage(stringOfLength50)
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
			jsonOfLength50, err = nosj.ParseJSON(fmt.Sprintf(`{"key":"%s"}`, strings.Repeat("a", 34)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should tell us what pointer we expected not to find", func() {
			Expect(matchers.HaveJSONPointer("/key").NegatedFailureMessage(shortJson)).
				To(ContainSubstring("expected '{nosj:map[key:val]}' not to contain the pointer '/key'"))
		})

		Context("when the input has a long string representation", func() {
			It("truncates that representation", func() {
				Expect(len(matchers.HaveJSONPointer("/beauty").NegatedFailureMessage(json))).To(BeNumerically("<", 102))
				Expect(matchers.HaveJSONPointer("/beauty").NegatedFailureMessage(json)).
					To(ContainSubstring("..."))
				Expect(matchers.HaveJSONPointer("/beauty").NegatedFailureMessage(json)).
					To(ContainSubstring("{nosj:map"))
			})
		})

		Context("when the input is exactly as large as we're willing to print", func() {
			It("prints it all, without elipses", func() {
				failureMessage := matchers.HaveJSONPointer("/key").NegatedFailureMessage(jsonOfLength50)
				Expect(failureMessage).To(ContainSubstring(fmt.Sprintf("'%+v'", jsonOfLength50)))
				Expect(failureMessage).NotTo(ContainSubstring("..."))
			})
		})
	})
})
