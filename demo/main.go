// demo implements a WASM module that can be used to format markdown
// with the goldmark-hashtag extension.
package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
)

func main() {
	js.Global().Set("tagList", js.FuncOf(func(this js.Value, args []js.Value) any {
		var req request
		req.Decode(args[0])

		return tagList(&req).Encode()
	}))

	select {}
}

type request struct {
	Markdown string
	Variant  string
}

func (r *request) Decode(v js.Value) {
	r.Markdown = v.Get("markdown").String()
	r.Variant = v.Get("variant").String()
}

type response struct {
	HTML string
	Tags []string
}

func (r *response) Encode() js.Value {
	tags := make([]any, len(r.Tags))
	for i, tag := range r.Tags {
		tags[i] = tag
	}

	return js.ValueOf(map[string]any{
		"html": r.HTML,
		"tags": tags,
	})
}

func tagList(r *request) *response {
	var variant hashtag.Variant
	switch r.Variant {
	case "", "default":
		variant = hashtag.DefaultVariant
	case "obsidian":
		variant = hashtag.ObsidianVariant
	default:
		panic(fmt.Sprintf("invalid variant: %q", r.Variant))
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			&hashtag.Extender{
				Variant:  variant,
				Resolver: _fakeResolver,
			},
		),
	)

	doc := md.Parser().Parse(text.NewReader([]byte(r.Markdown)))

	var tags []string
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == hashtag.Kind {
			tags = append(tags, "#"+string(n.(*hashtag.Node).Tag))
		}
		return ast.WalkContinue, nil
	})

	var buff bytes.Buffer
	md.Renderer().Render(&buff, []byte(r.Markdown), doc)

	return &response{
		HTML: buff.String(),
		Tags: tags,
	}
}

type fakeResolver struct{}

var _fakeResolver hashtag.Resolver = fakeResolver{}

func (fakeResolver) ResolveHashtag(n *hashtag.Node) ([]byte, error) {
	return []byte(fmt.Sprintf("#/tags/%s", n.Tag)), nil
}
