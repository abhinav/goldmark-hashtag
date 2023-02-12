package hashtag

import "go.abhg.dev/goldmark/hashtag"

// Extender extends a goldmark Markdown object with support for parsing and
// rendering hashtags.
//
// Install it on your Markdown object upon creation.
//
//	goldmark.New(
//	  goldmark.WithExtensions(
//	    // ...
//	    &hashtag.Extender{...},
//	  ),
//	  // ...
//	)
//
// Provide a Resolver to render tags as links that point to a specific
// destination.
type Extender = hashtag.Extender
