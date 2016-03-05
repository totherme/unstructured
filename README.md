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
        "beauty": true
      }`

json, err = nosj.Json(rawjson)
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
