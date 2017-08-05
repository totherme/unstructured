package unstructured_test

import (
	"reflect"

	"github.com/totherme/unstructured"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data", func() {
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

	It("looks the same, whether it's JSON or YAML", func() {
		rawyaml := `
name: "fred"
othernames:
- "alice"
- "bob"
- "ezekiel"
life: 42
things:
  more: "things"
beauty: true
not: null
`
		json, err := unstructured.ParseJSON(rawjson)
		Expect(err).NotTo(HaveOccurred())
		yaml, err := unstructured.ParseYAML(rawyaml)
		Expect(err).NotTo(HaveOccurred())
		Expect(json).To(BeEquivalentTo(yaml))

	})

	Context("when my Data represents an object", func() {
		var err error
		var json unstructured.Data
		BeforeEach(func() {
			json, err = unstructured.ParseJSON(rawjson)
		})

		It("parses the json successfully", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me that my json represents an object", func() {
			Expect(json.IsOb()).To(BeTrue(), "this json represents an object")
			simpleObJson, err := unstructured.ParseJSON(`{"string":1 , "otherstring":2}`)
			Expect(err).NotTo(HaveOccurred())
			Expect(simpleObJson.IsOb()).To(BeTrue(), "this json represents an object")
		})

		It("can get that object", func() {
			obVal := json.ObValue()
			Expect(obVal).To(HaveLen(6))
			Expect(obVal).To(HaveKey("name"))
			Expect(obVal).To(HaveKey("othernames"))
			Expect(obVal).To(HaveKey("life"))
			Expect(obVal).To(HaveKey("things"))
			Expect(obVal).To(HaveKey("beauty"))
			Expect(obVal).To(HaveKey("not"))
		})

		It("can update fields on that object", func() {
			json.SetField("name", "david")
			Expect(json.F("name").StringValue()).To(Equal("david"))
		})

		It("can add fields to that object", func() {
			json.SetField("newfield", "new value")
			Expect(json.F("newfield").StringValue()).To(Equal("new value"))
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

		It("tells me that extant pointers exist, and others do not", func() {
			Expect(json.HasPointer("/name")).To(BeTrue(), "the pointer should exist")
			Expect(json.HasPointer("/life")).To(BeTrue(), "the pointer should exist")
			Expect(json.HasPointer("/wat?")).To(BeFalse(), "the pointer should not exist")
			Expect(json.HasPointer("/things/more")).To(BeTrue(), "the pointer should exist")
			Expect(json.HasPointer("/not/there")).To(BeFalse(), "the pointer should not exist")
		})

		Context("when we pass an invalid pointer", func() {
			It("returns a helpful error message", func() {
				_, err = json.HasPointer("invalid/pointer")
				Expect(err).To(MatchError(ContainSubstring("JSON pointer must be empty or start with a \"/\"")))
			})
		})

		It("can get an extant key", func() {
			newJson := json.GetField("things")
			Expect(newJson.IsOb()).To(BeTrue(), "the inner object is also an object")
			Expect(newJson.HasKey("more")).To(BeTrue(), "the inner object has the 'more' field")
		})

		It("can chain extant keys", func() {
			Expect(json.GetField("things").GetField("more").StringValue()).To(Equal("things"))
		})

		It("can get by pointer", func() {
			got, err := json.GetByPointer("/things/more")
			Expect(err).NotTo(HaveOccurred())
			Expect(got.StringValue()).To(Equal("things"))
		})

		DescribeTable("F and GetByPointer both mirror GetField for single-level paths", func(key string) {
			Expect(json.F(key)).To(Equal(json.GetField(key)))
			Expect(json.GetByPointer("/" + key)).To(Equal(json.GetField(key)))
		},
			Entry("existing object key", "things"),
			Entry("existing string key", "name"),
			Entry("existing list key", "othernames"),
			Entry("existing number key", "life"),
			Entry("existing boolean key", "beauty"),
			Entry("existing null key", "not"),
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
				Expect(json.SetElem(0, "some-value")).To(MatchError(ContainSubstring("not a list")))
			})
		})

		Describe("error handling of GetByPointer", func() {
			Context("when we pass a pointer that is invalid", func() {
				It("returns a helpful error message", func() {
					_, err = json.GetByPointer("not/starting/with/slash")
					Expect(err).To(MatchError(ContainSubstring("JSON pointer must be empty or start with a \"/\"")))
				})
			})
			Context("when we pass a pointer to a non-existing key", func() {
				It("returns a helpful error message", func() {
					_, err = json.GetByPointer("/not/there")
					Expect(err).To(MatchError(ContainSubstring("Invalid token reference")))
				})
			})
		})

		It("has a raw value equal to the parsed JSON", func() {
			Expect(json.RawValue()).To(HaveLen(6))
			Expect(json.RawValue()).To(HaveKey("name"))
			Expect(json.RawValue().(map[string]interface{})["name"]).To(Equal("fred"))
			Expect(json.RawValue()).To(HaveKey("othernames"))
			Expect(json.RawValue().(map[string]interface{})["othernames"]).To(HaveLen(3))
			Expect(json.RawValue().(map[string]interface{})["othernames"]).To(ContainElement("alice"))
			Expect(json.RawValue().(map[string]interface{})["othernames"]).To(ContainElement("bob"))
			Expect(json.RawValue().(map[string]interface{})["othernames"]).To(ContainElement("ezekiel"))
			Expect(json.RawValue()).To(HaveKey("life"))
			Expect(json.RawValue().(map[string]interface{})["life"]).To(Equal(42.0))
			Expect(json.RawValue()).To(HaveKey("things"))
			Expect(json.RawValue().(map[string]interface{})["things"]).To(Equal(map[string]interface{}{
				"more": "things",
			}))
			Expect(json.RawValue()).To(HaveKey("beauty"))
			Expect(json.RawValue().(map[string]interface{})["beauty"]).To(Equal(true))
			Expect(json.RawValue()).To(HaveKey("not"))
			Expect(json.RawValue().(map[string]interface{})["not"]).To(BeNil())
		})
	})

	Context("when my json represents a string", func() {
		var json unstructured.Data
		var err error
		BeforeEach(func() {
			json, err = unstructured.ParseJSON(`"this is a string"`)
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
			Expect(json.StringValue()).To(Equal("this is a string"))
			Expect(json.RawValue()).To(Equal("this is a string"))
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
		var json unstructured.Data
		var err error

		BeforeEach(func() {
			json, err = unstructured.ParseJSON(`3.141`)
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
			Expect(json.RawValue()).To(BeNumerically("==", 3.141))
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
		var json unstructured.Data
		var err error

		BeforeEach(func() {
			json, err = unstructured.ParseJSON(`true`)
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
			Expect(json.RawValue()).To(BeTrue(), "actually should be the value 'true'")
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
		var json unstructured.Data
		var err error

		BeforeEach(func() {
			json, err = unstructured.ParseJSON(`[true, 32, {"this":"that"}]`)
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

		It("can get that list in raw form", func() {
			Expect(json.RawValue()).To(HaveLen(3))
			Expect(reflect.TypeOf(json.RawValue())).To(Equal(reflect.TypeOf([]interface{}{})))
			Expect(reflect.TypeOf(json.RawValue().([]interface{})[0])).To(Equal(reflect.TypeOf(true)))
			Expect(json.RawValue().([]interface{})[0]).To(BeTrue())
			Expect(json.RawValue().([]interface{})[1]).To(Equal(32.0))
		})

		It("can set items in that list", func() {
			err := json.SetElem(1, "badgers")
			Expect(err).NotTo(HaveOccurred())
			Expect(json.ListValue()[1].StringValue()).To(Equal("badgers"))
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
		var json unstructured.Data
		var err error

		BeforeEach(func() {
			json, err = unstructured.ParseJSON(`null`)
			Expect(err).NotTo(HaveOccurred())
		})

		It("tells me my json represents null", func() {
			Expect(json.IsNull()).To(BeTrue(), "this is null")
		})

		It("has raw value nil", func() {
			Expect(json.RawValue()).To(BeNil())
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
			_, err := unstructured.ParseJSON("this isn't even slightly json")
			Expect(err).To(MatchError(ContainSubstring("parse error")))
		})
	})

	Describe("the IsOfType convenience method", func() {
		var json unstructured.Data
		BeforeEach(func() {
			var err error
			json, err = unstructured.ParseJSON(rawjson)
			Expect(err).NotTo(HaveOccurred())
		})
		DescribeTable("IsOfType does the same as the individual type methods", func(key string) {
			field := json.F(key)
			Expect(field.IsOfType(unstructured.DataOb)).To(Equal(field.IsOb()))
			Expect(field.IsOfType(unstructured.DataString)).To(Equal(field.IsString()))
			Expect(field.IsOfType(unstructured.DataList)).To(Equal(field.IsList()))
			Expect(field.IsOfType(unstructured.DataNum)).To(Equal(field.IsNum()))
			Expect(field.IsOfType(unstructured.DataBool)).To(Equal(field.IsBool()))
			Expect(field.IsOfType(unstructured.DataNull)).To(Equal(field.IsNull()))
		},
			Entry("an object key", "things"),
			Entry("an string key", "name"),
			Entry("an list key", "othernames"),
			Entry("an number key", "life"),
			Entry("an boolean key", "beauty"),
			Entry("a null key", "not"),
		)

		Context("when we give a string that isn't a Data type", func() {
			It("panics", func() {
				Expect(func() { json.IsOfType("badgers") }).To(Panic())
			})
		})
	})
})
