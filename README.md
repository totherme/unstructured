# NOSJ - Tiny Simple JSON Reflection in Go

This is to make testing JSON output of whatever-you're-doing a tiny bit easier.

[![GoDoc](https://godoc.org/github.com/totherme/nosj?status.svg)](https://godoc.org/github.com/totherme/nosj)

## Quick start

For full details of the below examples, see [this file](example_usage_test.go).

First: import [nosj](http://github.com/totherme/nosj),
[ginkgo](http://github.com/onsi/ginkgo) and
[gomega](http://github.com/onsi/gomega):

```go
import (
	"github.com/totherme/nosj"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)
```

Next, get yourself some nosj.JSON:

```go
rawjson := `{...}`

myjson, err = nosj.ParseJSON(rawjson)
Expect(err).NotTo(HaveOccurred())
```

Test!

```go
It("contains three employees", func() {
	employees, err := myjson.GetByPointer("/employees")
	Expect(err).NotTo(HaveOccurred())
	Expect(employees).To(BeAList())
	Expect(employees.ListValue()).To(HaveLen(3))
})

Describe("the first employee", func() {
	It("is great at cooking", func() {
		skill, err := myjson.GetByPointer("/employees/0/profile/special-skill")
		Expect(err).NotTo(HaveOccurred())
		Expect(skill).To(BeAString())
		Expect(skill.StringValue()).To(Equal("szechuan cookery"))
	})
})
```

In addition to nosj's JSON-wrangling, these tests make use of the
[ginkgo](http://github.com/onsi/ginkgo) BDD DSL, which provides constructs like
`Describe` and `It`; and the [gomega](http://github.com/onsi/gomega) matcher
library, which gives us `Expect`.

## Getting nosj.JSON Values

There are two approaches you can take to pulling values out of a nosj.JSON
object. Firstly, if you like [JSON
pointers](https://tools.ietf.org/html/rfc6901), you can use them:

```go
morejson, err := myjson.GetByPointer("/path/to/property")
```

The other approach is to get by one key at a time. For this, you can use
`GetField` or its alias `F`:

```go
morejson := myjson.GetField("path").GetField("to").GetField("property")
// or equivalently:
morejson := myjson.F("path").F("to").F("property")
```

Note that all these methods will return another nosj.JSON object to the
variable `morejson`.  To do anything interesting, you'll probably want to get a
golang value out of the nosj.JSON object, as documented in 'Getting Golang Values'
below. Finally, notice that while `GetByPointer`'s return type includes an
error (e.g. when the pointer does not exist in this nosj.JSON), the return type of
`GetField` does not. This allows `GetField` to be chained (and thus makes up
for single-field names being significantly less expressive than JSON pointers),
but does mean that `GetField` will panic in those occasions where
`GetByPointer` would return a helpful error message.


## Testing nosj.JSON Values

nosj.JSON values may represent data of a variety of golang types. To discover these
types, you can use methods such as:

```go
isJSONBool := myjson.IsBool()
isJSONNum := myjson.IsNum()
isJSONString := myjson.IsString()
isJSONOb := myjson.IsOb()
isJSONList := myjson.IsList()
isJSONNull := myjson.IsNull()
```

If your nosj.JSON object represents a JSON object (as opposed to, for example,
a JSON String), you may also want to check if that object has a field of a
particular name:

```go
hasField := myjson.HasKey("particular-name")
```

## Getting Golang Values

If you are sure what JSON type is represented by your nosj.JSON object, you can
get a value of the appropriate golang type like so:

```go
myGolangBool := myjson.BoolValue()
myGolangNum := myjson.NumValue()
myGolangString := myjson.StringValue()
myGolangOb := myjson.ObValue()
myGolangList := myjson.ListValue()
```

## Matchers

If you're using [ginkgo](http://github.com/onsi/ginkgo) and
[gomega](http://github.com/onsi/gomega), you might prefer to use a
[matcher](gnosj) rather than test your JSON directly:

```go
import (
	"github.com/totherme/nosj"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gnosj matchers", func() {
  It("does the same thing as our test methods", func() {
    // Type matchers
    Expect(myjson.IsBool()).To(BeTrue())
    Expect(myjson).To(BeABool())

    Expect(myjson.IsNum()).To(BeTrue())
    Expect(myjson).To(BeANum())

    Expect(myjson.IsString()).To(BeTrue())
    Expect(myjson).To(BeAString())

    Expect(myjson.IsOb()).To(BeTrue())
    Expect(myjson).To(BeAnObject())

    Expect(myjson.IsList()).To(BeTrue())
    Expect(myjson).To(BeAList())

    Expect(myjson.IsNull()).To(BeTrue())
    Expect(myjson).To(BeANull())

    // Access matchers
    Expect(myjson.HasKey("my-key")).To(BeTrue())
    Expect(myjson).To(HaveJSONKey("my-key"))

    Expect(myjson.HasPointer("/my/pointer")).To(BeTrue())
    Expect(myjson).To(HaveJSONPointer("/my/pointer"))
  })
})
```

Each pair of expectations in the above block is composed of two equivalent
statements.
