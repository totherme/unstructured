package nosj_test

import (
	"github.com/totherme/nosj"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Json", func() {
	Context("when I look at some valid JSON", func() {
		var err error
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
							"beauty": true
						}`

			json, err = nosj.Json(rawjson)

		})

		It("should parse the json successfully", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should tell me that precisely the keys name, othernames and life exist", func() {
			Expect(json.HasKey("name")).To(BeTrue(), "the name key should exist")
			Expect(json.HasKey("othernames")).To(BeTrue(), "the othernames key should exist")
			Expect(json.HasKey("life")).To(BeTrue(), "the life key should exist")
			Expect(json.HasKey("wat?")).To(BeFalse(), "the wat key should not exist")
		})

		It("should give accurate type reports about the keys", func() {
			Expect(json.IsString("name")).To(BeTrue(), "name is a string")
			Expect(json.IsString("othernames")).To(BeFalse(), "othernames is not a string")
			Expect(json.IsList("othernames")).To(BeTrue(), "othernames is a list")
			Expect(json.IsList("life")).To(BeFalse(), "life is not a list")
			Expect(json.IsNum("life")).To(BeTrue(), "life is a number")
			Expect(json.IsNum("name")).To(BeFalse(), "name is not a number")
			Expect(json.IsOb("things")).To(BeTrue(), "things is an object")
			Expect(json.IsOb("beauty")).To(BeFalse(), "beauty is not an object")
			Expect(json.IsBool("beauty")).To(BeTrue(), "beauty is a bool")
			Expect(json.IsBool("things")).To(BeFalse(), "things is not a bool")
		})
	})

	Context("when I look at some invalid JSON", func() {
		It("returns a helpful error", func() {
			_, err := nosj.Json("this isn't even slightly json")
			Expect(err).To(MatchError(ContainSubstring("parse error")))
		})
	})
})
