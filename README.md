# NOSJ - Tiny Simple JSON Reflection in Go

This is to make testing JSON output of whatever-you're-doing a tiny bit easier.

Import [nosj](http://github.com/totherme/nosj),
[ginkgo](http://github.com/onsi/ginkgo) and
[gomega](http://github.com/onsi/gomega):

```go
import (
	"github.com/totherme/nosj"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)
```

Get yourself some JSON:

```go
rawjson :=
  `{"name": "fred",
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

json, err = nosj.ParseJSON(rawjson)
Expect(err).NotTo(HaveOccurred())
```

Test that it looks the way you expect:

```go
It("has a 'life' field with a number in", func() {
  Expect(json.HasKey("life")).To(BeTrue())
  Expect(json.GetField("life").IsNum()).To(BeTrue())
})

It("contains deep things", func() {
  Expect(json.GetField("things").GetField("more").StringValue()).To(Equal("things"))
})
```

If you get bored typing `GetField`, then try `F`

```go
It("contains deep things", func() {
  Expect(json.F("things").F("more").StringValue()).To(Equal("things"))
})
```

A disadvantage of using `GetField` is that it panics rather than raising an 
error if the field does not exist. This makes chaining possible, but can make 
debugging difficult. 

You can avoid this problem with json pointers. Json pointers provide an
 alternative way to access information deep within a json structure:

```go
It("contains deep things", func() {
  Expect(json.HasPointer("/things/more")).To(BeTrue(), "the pointer should exist")
    
  got, err := json.GetByPointer("/things/more")
  Expect(err).NotTo(HaveOccurred())
  Expect(got.StringValue()).To(Equal("things"))
})
```
`GetByPointer` returns a helpful error if the pointer does not exist in the json. 