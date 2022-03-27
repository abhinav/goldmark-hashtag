package hashtag

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extender extends a goldmark Markdown object with support for parsing and
// rendering hashtags.
//
// Install it on your Markdown object upon creation.
//
//  goldmark.New(
//    goldmark.WithExtensions(
//      // ...
//      &hashtag.Extender{...},
//    ),
//    // ...
//  )
//
// Provide a Resolver to render tags as links that point to a specific
// destination.
type Extender struct {
	// Resolver specifies destination links for hashtags, if any.
	//
	// Defaults to no links.
	Resolver Resolver

	// Options is a list of options for this extension.
	Options []Option
}

var _ goldmark.Extender = (*Extender)(nil)

// NewExtender creates a new goldmark.Extender to extend goldmark Markdown
// object with support for parsing and rendering hashtags.
//
// Install it on your Markdown object upon creation.
//
//  goldmark.New(
//    goldmark.WithExtensions(
//      // ...
//      &hashtag.NewExtender(...),
//    ),
//    // ...
//  )
//
// See Extender for more.
func NewExtender(resolver Resolver, opts ...Option) goldmark.Extender {
	return &Extender{
		Resolver: resolver,
		Options:  opts,
	}
}

// Extend extends the provided goldmark Markdown object with support for
// hashtags.
func (e *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewParser(e.Options...), 999),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{
				Resolver: e.Resolver,
			}, 999),
		),
	)
}
