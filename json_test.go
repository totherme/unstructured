package nosj_test

import (
	"reflect"

	"github.com/totherme/nosj"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON", func() {
	Context("when my JSON represents an object", func() {
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

			json, err = nosj.ParseJSON(rawjson)

		})

		It("parses the json successfully", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me that my json represents an object", func() {
			Expect(json.IsOb()).To(BeTrue(), "this json represents an object")
			simpleObJson, err := nosj.ParseJSON(`{"string":1 , "otherstring":2}`)
			Expect(err).NotTo(HaveOccurred())
			Expect(simpleObJson.IsOb()).To(BeTrue(), "this json represents an object")
		})

		It("tells me it doesn't represent anything else", func() {
			Expect(json.IsString()).To(BeFalse(), "not a string")
			Expect(json.IsNum()).To(BeFalse(), "not a number")
			Expect(json.IsBool()).To(BeFalse(), "not a bool")
			Expect(json.IsList()).To(BeFalse(), "not a list")
			Expect(json.IsNull()).To(BeFalse(), "not null")
		})

		It("tells me that extant keys exist, and others do not", func() {
			Expect(json.HasKey("name")).To(BeTrue(), "the name key should exist")
			Expect(json.HasKey("othernames")).To(BeTrue(), "the othernames key should exist")
			Expect(json.HasKey("life")).To(BeTrue(), "the life key should exist")
			Expect(json.HasKey("wat?")).To(BeFalse(), "the wat key should not exist")
		})

		It("can get an extant key", func() {
			newJson := json.GetField("things")
			Expect(newJson.IsOb()).To(BeTrue(), "the inner object is also an object")
			Expect(newJson.HasKey("more")).To(BeTrue(), "the inner object has the 'more' field")
		})

		It("can chain extant keys", func() {
			Expect(json.GetField("things").GetField("more").StringValue()).To(Equal("things"))
		})

		DescribeTable("F does the same thing as GetField", func(key string) {
			Expect(json.F(key)).To(Equal(json.GetField(key)))
		},
			Entry("existing object key", "things"),
			Entry("existing string key", "name"),
			Entry("existing list key", "othernames"),
			Entry("existing number key", "life"),
			Entry("existing boolean key", "beauty"),
		)

		Context("when I try to get a key that doesn't exist", func() {
			It("panics", func() {
				Expect(func() { json.GetField("oh noe!") }).To(Panic())
				Expect(func() { json.F("oh noe!") }).To(Panic())
			})
		})

		Context("when I try to do non-objectey things with it", func() {
			It("panics", func() {
				Expect(func() { json.StringValue() }).To(Panic())
				Expect(func() { json.NumValue() }).To(Panic())
				Expect(func() { json.BoolValue() }).To(Panic())
				Expect(func() { json.ListValue() }).To(Panic())
			})
		})
	})

	Context("when my json represents a string", func() {
		var json nosj.JSON
		var err error
		BeforeEach(func() {
			json, err = nosj.ParseJSON(`"this is a string"`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me my json represents a string", func() {
			Expect(json.IsString()).To(BeTrue(), "this json represents a string")
		})

		It("tells me it doesn't represent anything else", func() {
			Expect(json.IsOb()).To(BeFalse(), "not an object")
			Expect(json.IsNum()).To(BeFalse(), "not a number")
			Expect(json.IsBool()).To(BeFalse(), "not a bool")
			Expect(json.IsList()).To(BeFalse(), "not a list")
			Expect(json.IsNull()).To(BeFalse(), "not null")
		})

		It("can get that string", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(json.StringValue()).To(Equal("this is a string"))
		})

		Context("when I try to do non-string things", func() {
			It("panics", func() {
				Expect(func() { json.HasKey("wat?") }).To(Panic())
				Expect(func() { json.GetField("oh noe!") }).To(Panic())
				Expect(func() { json.NumValue() }).To(Panic())
				Expect(func() { json.BoolValue() }).To(Panic())
				Expect(func() { json.ListValue() }).To(Panic())
			})
		})
	})

	Context("when my json represents a number", func() {
		var json nosj.JSON
		var err error

		BeforeEach(func() {
			json, err = nosj.ParseJSON(`3.141`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me my json represents a number", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(json.IsNum()).To(BeTrue(), "this is a number")
		})

		It("tells me it doesn't represent anything else", func() {
			Expect(json.IsOb()).To(BeFalse(), "not an object")
			Expect(json.IsString()).To(BeFalse(), "not a string")
			Expect(json.IsBool()).To(BeFalse(), "not a bool")
			Expect(json.IsList()).To(BeFalse(), "not a list")
			Expect(json.IsNull()).To(BeFalse(), "not null")
		})

		It("can get that number", func() {
			Expect(json.NumValue()).To(BeNumerically("==", 3.141))
		})

		Context("when I try to do non-number things", func() {
			It("panics", func() {
				Expect(func() { json.HasKey("wat?") }).To(Panic())
				Expect(func() { json.GetField("oh noe!") }).To(Panic())
				Expect(func() { json.StringValue() }).To(Panic())
				Expect(func() { json.BoolValue() }).To(Panic())
				Expect(func() { json.ListValue() }).To(Panic())
			})
		})
	})

	Context("when my json represents a bool", func() {
		var json nosj.JSON
		var err error

		BeforeEach(func() {
			json, err = nosj.ParseJSON(`true`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me my json represents a bool", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(json.IsBool()).To(BeTrue(), "this is a bool")
		})

		It("tells me it doesn't represent anything else", func() {
			Expect(json.IsOb()).To(BeFalse(), "not an object")
			Expect(json.IsString()).To(BeFalse(), "not a string")
			Expect(json.IsNum()).To(BeFalse(), "not a number")
			Expect(json.IsList()).To(BeFalse(), "not a list")
			Expect(json.IsNull()).To(BeFalse(), "not null")
		})

		It("can get that bool", func() {
			Expect(json.BoolValue()).To(BeTrue(), "actually should be the value 'true'")
		})

		Context("when I try to do non-bool things", func() {
			It("panics", func() {
				Expect(func() { json.HasKey("wat?") }).To(Panic())
				Expect(func() { json.GetField("oh noe!") }).To(Panic())
				Expect(func() { json.StringValue() }).To(Panic())
				Expect(func() { json.NumValue() }).To(Panic())
				Expect(func() { json.ListValue() }).To(Panic())
			})
		})
	})

	Context("when my json represents a list", func() {
		var json nosj.JSON
		var err error

		BeforeEach(func() {
			json, err = nosj.ParseJSON(`[true, 32, {"this":"that"}]`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me my json represents a list", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(json.IsList()).To(BeTrue(), "this is a list")
		})

		It("tells me it doesn't represent anything else", func() {
			Expect(json.IsOb()).To(BeFalse(), "not an object")
			Expect(json.IsString()).To(BeFalse(), "not a string")
			Expect(json.IsNum()).To(BeFalse(), "not a number")
			Expect(json.IsBool()).To(BeFalse(), "not a bool")
			Expect(json.IsNull()).To(BeFalse(), "not null")
		})

		It("can get that list", func() {
			Expect(json.ListValue()).To(HaveLen(3))
			Expect(reflect.TypeOf(json.ListValue()[0])).To(Equal(reflect.TypeOf(json)))
			Expect(json.ListValue()[0].IsBool()).To(BeTrue())
			Expect(json.ListValue()[1].IsNum()).To(BeTrue())
		})

		Context("when I try to do non-list things", func() {
			It("panics", func() {
				Expect(func() { json.HasKey("wat?") }).To(Panic())
				Expect(func() { json.GetField("oh noe!") }).To(Panic())
				Expect(func() { json.StringValue() }).To(Panic())
				Expect(func() { json.NumValue() }).To(Panic())
				Expect(func() { json.BoolValue() }).To(Panic())
			})
		})
	})

	Context("when my json represents null", func() {
		var json nosj.JSON
		var err error

		BeforeEach(func() {
			json, err = nosj.ParseJSON(`null`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me my json represents null", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(json.IsNull()).To(BeTrue(), "this is null")
		})

		It("tells me it doesn't represent anything else", func() {
			Expect(json.IsOb()).To(BeFalse(), "not an object")
			Expect(json.IsString()).To(BeFalse(), "not a string")
			Expect(json.IsNum()).To(BeFalse(), "not a number")
			Expect(json.IsBool()).To(BeFalse(), "not a bool")
			Expect(json.IsList()).To(BeFalse(), "not a list")
		})

		Context("when I try to do ...well... things", func() {
			It("panics", func() {
				Expect(func() { json.HasKey("wat?") }).To(Panic())
				Expect(func() { json.GetField("oh noe!") }).To(Panic())
				Expect(func() { json.StringValue() }).To(Panic())
				Expect(func() { json.NumValue() }).To(Panic())
				Expect(func() { json.BoolValue() }).To(Panic())
				Expect(func() { json.ListValue() }).To(Panic())
			})
		})
	})

	Context("when I look at some invalid JSON", func() {
		It("returns a helpful error", func() {
			_, err := nosj.ParseJSON("this isn't even slightly json")
			Expect(err).To(MatchError(ContainSubstring("parse error")))
		})
	})
})
