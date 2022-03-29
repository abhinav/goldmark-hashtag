[![Go Reference](https://pkg.go.dev/badge/github.com/abhinav/goldmark-hashtag.svg)](https://pkg.go.dev/github.com/abhinav/goldmark-hashtag)
[![Go](https://github.com/abhinav/goldmark-hashtag/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-hashtag/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-hashtag/branch/main/graph/badge.svg?token=w6jkI2SQ9u)](https://codecov.io/gh/abhinav/goldmark-hashtag)

goldmark-hashtag is an extension for the [goldmark] Markdown parser that adds
support for tagging documents with `#hashtag`s.

  [goldmark]: http://github.com/yuin/goldmark

# Usage

To use goldmark-hashtag, import the `hashtag` package.

```go
import hashtag "github.com/abhinav/goldmark-hashtag"
```

Then include the `hashtag.Extender` in the list of extensions you build your
[`goldmark.Markdown`] with.

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  &hashtag.Extender{},
  // ...
).Convert(src, out)
```

This alone has little effect besides adding `<span class="hashtag">...</span>`
around hashtags in your Markdown document.

Supply a `hashtag.Resolver` to the `hashtag.Extender` to render hashtags as
links:

```go
goldmark.New(
  &hashtag.Extender{
    Resolver: hashtagResolver,
  },
  // ...
).Convert(src, out)
```

## Syntax

Hashtags must always begin with a "#".
The characters that follow that depend on the variant you have chosen.
goldmark-hashtag supports the following variants:

- *Default*: Hashtags must begin with a letter, and may contain letters,
  numbers and any of the following symbols: `/_-`.
  goldmark-hashtag uses this variant if you do not specify one.
- *Obsidian*: Hashtags can begin with and contain letters, numbers, emoji, and
  any of the following symbols: `/_-`, but must not contain only numbers.

You can specify the variant by setting the `Variant` property of the
`hashtag.Extender`.

```go
&hashtag.Extender{
  // ...
  Variant: hashtag.ObsidianVariant,
}
```

## Inspection

To collect all hashtags from a Markdown document, use Goldmark's [`ast.Walk`]
function after parsing the document.

  [`ast.Walk`]: https://pkg.go.dev/github.com/yuin/goldmark/ast#Walk

For example, the following will populate the `hashtags` map with all hashtags
found in the document.

```go
markdown := goldmark.New(
  &hashtag.Extender{
    Resolver: hashtagResolver,
  },
  // ...
)

// Parse the Markdown document.
doc := markdown.Parser().Parse(text.NewReader(src))

// List the tags.
hashtags := make(map[string]struct{})
ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
  if n, ok := node.(*hashtag.Node); ok && enter {
    hashtags[string(n.Tag)] = struct{}{}
  }
  return ast.WalkContinue, nil
})
```
