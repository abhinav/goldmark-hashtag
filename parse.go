package hashtag

import "go.abhg.dev/goldmark/hashtag"

// Parser is a Goldmark inline parser for parsing hashtag nodes.
//
// Hashtags start with "#". The list of other characters allowed in the hashtag
// is determined by variant. See the documentation for Variant for more
// details.
type Parser = hashtag.Parser

// Variant represents one of the different flavours of hashtag syntax.
type Variant = hashtag.Variant

const (
	// DefaultVariant is the default flavor of hashtag syntax supported by
	// this package.
	//
	// In this format, hashtags start with "#" and an alphabet, followed by
	// zero or more alphanumeric characters and the following symbols.
	//
	//   /_-
	DefaultVariant = hashtag.DefaultVariant

	// ObsidianVariant is a flavor of the hashtag syntax that aims to be
	// compatible with Obsidian (https://obsidian.md/).
	//
	// In this format, hashtags start with "#" followed by alphabets,
	// numbers, emoji, or any of the following symbols.
	//
	//   /_-
	//
	// Hashtags cannot be entirely numeric and must contain at least one
	// non-numeric character.
	//
	// See also https://help.obsidian.md/How+to/Working+with+tags.
	ObsidianVariant = hashtag.ObsidianVariant
)
