package main

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/totherme/unstructured"
)

func main() {
	myYaml := `
top-level-list:
- "this first element is a string - perhaps containing metadata"
- name: first real element
  type: element-type-1
  payload: [1,2,3,4]
- name: second real element
  type: element-type-2
  payload:
    some: embedded structure
`

	myData, err := unstructured.ParseYAML(myYaml)
	if err != nil {
		panic("Couldn't parse my own yaml")
	}

	myPayloadData, err := myData.GetByPointer("/top-level-list/2/payload/some")
	if err != nil {
		panic("Couldn't address into my own yaml")
	}

	if !myPayloadData.IsString() {
		panic("I really thought that was a string...")
	}
	fmt.Println(myPayloadData.StringValue())

	myPayloadMap, err := myData.GetByPointer("/top-level-list/2/payload")
	if err != nil {
		panic("Couldn't address into my own yaml")
	}

	if !myPayloadMap.IsOb() {
		panic("I can't write into this object if it's not an object")
	}

	myPayloadMap.SetField("additional-key", []string{"some", "arbitrary", "data"})

	outputYaml, err := yaml.Marshal(myData.RawValue())
	if err != nil {
		panic("myData should definitely still be serializable")
	}
	fmt.Println(string(outputYaml))
}
