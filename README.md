# Unstructured - for when you can't fit your data into a struct

[![GoDoc](https://godoc.org/github.com/totherme/unstructured?status.svg)](https://godoc.org/github.com/totherme/unstructured)

## Why?

Go is awesome at [de-]serialising structured data. We can do things like:

```go
type MyType struct {
  Key1 string `yaml:"key1"`
  Key2 int    `yaml:"key2"`
}

...

myYaml = `
key1: "Alice"
key2: 65537
`

var myVar MyType
err := yaml.Unmarshal([]byte(myYaml), &myVar)
```

This is fantastic so long as:
1. You know the exact structure of your data at compile time.
1. That structure can be mapped to a go struct type.

For an example of point (2) failing, consider the following yaml:

```yaml
top-level-list:
- "this first element is a string - perhaps containing metadata"
- name: first real element
  type: element-type-1
  payload: [1,2,3,4]
- name: second real element
  type: element-type-2
  payload: 
    some: embedded structure
```

I do not know of a way to define a struct type in go which will accept this
YAML when I attempt to deserialize it.  If we're designing our systems
end-to-end in go, we will almost certainly avoid producing YAML structures like
this. However, sometimes we may have to interface with other systems written
perhaps in more dynamically typed languages, where this sort of thing is more
natural.

For an example of point (1) failing, consider the schema of a [bosh
manifest](https://bosh.io/docs/manifest-v2.html). The type of the properties
block of an instance group is described as a "hash". In practice, these hashes
can be quite complex structures, depending on how configurable the jobs in that
particular instance group are. Unfortunately, the schema for each individual
properties block is defined by the authors of the jobs in the instance group.
If we are writing a tool in go that manages bosh manifests in general, then we
have no way to be more specific than the type `map[string]interface{}`.

## What?

When working with data which cannot be described at compile time in the golang
type language, we have no choice but to either leave go and work in some other
language, or to work without the safety net of our type system. This library
attempts to make managing unstructured data in an untyped way less unpleasant
than it might otherwise be.

## How?

This library leans on the excellent
[gojsonpointer](https://github.com/xeipuuv/gojsonpointer) library to allow us
to address deep into JSON and YAML structures and:
- retrieve data -- if it exists
- write data -- if the parent we're writing into exists

This allows us to handle the data above with [code something
like](examples/usage.go):

```go
	myData, err := unstructured.ParseYAML(myYaml)
	myPayloadData, err := myData.GetByPointer("/top-level-list/2/payload/some")
	fmt.Println(myPayloadData.StringValue())
	myPayloadMap, err := myData.GetByPointer("/top-level-list/2/payload")
	myPayloadMap.SetField("additional-key", []string{"some", "arbitrary", "data"})
```

However, do see the "Gotchas" section below. The unchecked `err` assignments in
that code block aren't the only dangerous bits. Click through to the full
example to see the fully-checked version.

We also provide a number of [gomega](https://onsi.github.io/gomega) matchers in
case you want to inspect semi-structured data in your tests. You can see these
used [here](examples/usage_test.go)

## Gotchas

Since we're deliberately working around go's type system, there's quite a high
risk that our program might panic if we're not careful. All known panics are
documented in the API documentation for each method, but here are the main
points:

- If you try to get a real data value of the wrong type from some
  unstructured.Data, then you'll get a panic. For example `d.StringValue()`
  will panic if the data in `d` is actually a list. You can always check the
  type with methods like `d.IsString()`.
- If you use the unsafe field accessor methods to get a field that doesn't
  exist, then you'll get a panic. For example
  `d.F("parent-field").F("sub-field")` will panic unless `d` is a map with key
  `"parent-field"` which contains a map with key `"sub-field"`. You can check
  these things with methods like `IsOb` and `HasKey`, but you should probably
  just use `GetByPointer` instead. The unsafe accessors are only really useful
  for writing terse chained statements in tests, when you're almost entirely
  certain that the path exists, and a panic would be a correctly failing test
  anyway.
- If you try to set a field of an unstructured.Data which does not represent a
  map. For example, `d.SetField("field-name", "value")` will panic if `d`
  actually represents a string. Use `IsOb` to check first.
