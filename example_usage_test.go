package nosj_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/totherme/nosj"
	. "github.com/totherme/nosj/gnosj"
)

var _ = Describe("Some things you can test with nosj", func() {
	Context("when I have some JSON", func() {
		var rawjson string = `
		{ "id": 1,
			"type": "hr-data",
			"employees": [
				{ "id": 1,
					"type": "employee",
					"name": "Alex",
					"salary-band": 12,
					"department": "finance",
					"profile": {
						"favourite-animal": "cat",
						"special-skill": "szechuan cookery",
						"number-of-reports": 4
					}
				},
				{ "id": 2,
					"type": "employee",
					"name": "Sue",
					"salary-band": 8,
					"department": "logistics",
					"profile": {
						"favourite-animal": "dog",
						"special-skill": "rock climbing",
						"number-of-reports": 12
					}
				},
				{ "id": 3,
					"type": "employee",
					"name": "Hilary",
					"salary-band": 14,
					"department": "engineering",
					"profile": {
						"favourite-animal": "hedgehog",
						"special-skill": "archery",
						"number-of-reports": 2
					}
				}
			]
		}`
		var json nosj.JSON

		BeforeEach(func() {
			var err error
			json, err = nosj.ParseJSON(rawjson)
			Expect(err).NotTo(HaveOccurred())
		})

		It("contains three employees", func() {
			employees, err := json.GetByPointer("/employees")
			Expect(err).NotTo(HaveOccurred())
			Expect(employees).To(BeAList())
			Expect(employees.ListValue()).To(HaveLen(3))
		})

		Describe("the first employee", func() {
			It("is great at cooking", func() {
				skill, err := json.GetByPointer("/employees/0/profile/special-skill")
				Expect(err).NotTo(HaveOccurred())
				Expect(skill).To(BeAString())
				Expect(skill.StringValue()).To(Equal("szechuan cookery"))
			})
		})
	})
})
