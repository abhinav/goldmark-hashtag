package hashtag

import "go.abhg.dev/goldmark/hashtag"

// Resolver resolves hashtags to pages they should link to.
type Resolver = hashtag.Resolver

// Renderer renders hashtag nodes into HTML, optionally linking them to
// specific pages.
//
//	#foo
//
// Renders as the following by default.
//
//	<span class="hashtag">#foo</span>
//
// Supply a Resolver that returns a non-empty destination to render it like
// the following.
//
//	<span class="hashtag"><a href="...">#foo</a></span>
type Renderer = hashtag.Renderer
