package hashtag

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"github.com/forPelevin/gomoji"
	"github.com/rivo/uniseg"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Parser is a Goldmark inline parser for parsing hashtag nodes.
//
// Hashtags start with "#". The list of other characters allowed in the hashtag
// is determined by variant. See the documentation for Variant for more
// details.
type Parser struct {
	ParserConfig
}

// ParserConfig is a set of configuration options for the Parser.
type ParserConfig struct {
	// Variant is the tags sytax to parse for.
	Variant Variant
}

// Variant represents one of the different flavours of hashtag syntax.
type Variant uint

const (
	// DefaultVariant is the default flavor of hashtag syntax supported by
	// this package.
	//
	// In this format, hashtags start with "#" and an alphabet, followed by
	// zero or more alphanumeric characters and the following symbols.
	//
	//   /_-
	DefaultVariant Variant = iota

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
	ObsidianVariant
)

var _ parser.InlineParser = (*Parser)(nil)

var _hash = byte('#')

// NewParser creates a new parser.InlineParser to parse hashtags.
func NewParser(opts ...Option) parser.InlineParser {
	p := &Parser{}
	for _, o := range opts {
		if _, ok := o.(*withObsidianTags); ok {
			p.Variant = ObsidianVariant
		}
	}
	return p
}

// Trigger reports characters that trigger this parser.
func (*Parser) Trigger() []byte {
	return []byte{_hash}
}

// Parse parses a hashtag node.
func (p *Parser) Parse(parent ast.Node, block text.Reader, pctx parser.Context) ast.Node {
	line, seg := block.PeekLine()

	if len(line) == 0 || line[0] != _hash {
		return nil
	}
	line = line[1:]

	switch p.Variant {
	default:
		// Hashtag must start with a letter.
		start, sz := utf8.DecodeRune(line)
		if !unicode.IsLetter(start) {
			return nil
		}
		line = line[sz:]

		// Truncate seg down to "#foo".
		if i := bytes.IndexFunc(line, endOfHashtag); i >= 0 {
			seg = seg.WithStop(seg.Start + i + 1 + sz) // + '#' + start
		} // else { line ends with a "#foo" so seg remains unchanged }

	case ObsidianVariant:
		if end := endOfObsidianHashtag(line); end >= 0 {
			line = line[:end]
			seg = seg.WithStop(seg.Start + end + 1)
		}
		// if cannot find something that's not a digit, it's not a valid tag
		if i := bytes.IndexFunc(line, func(r rune) bool { return !unicode.IsDigit(r) }); i == -1 {
			return nil
		}
	}

	n := Node{
		Tag: block.Value(seg.WithStart(seg.Start + 1)), // omit the "#"
	}
	n.AppendChild(&n, ast.NewTextSegment(seg))
	block.Advance(seg.Len())
	return &n
}

func endOfObsidianHashtag(line []byte) int {
	gr := uniseg.NewGraphemes(string(line))
	for gr.Next() {
		rs := gr.Runes()
		if (len(rs) == 1 && endOfHashtag(rs[0])) && !gomoji.ContainsEmoji(gr.Str()) {
			pos, _ := gr.Positions()
			return pos
		}
	}
	return -1
}

func endOfHashtag(r rune) bool {
	return !(unicode.IsLetter(r) ||
		unicode.IsDigit(r) ||
		r == '_' || r == '-' || r == '/')
}

const (
	optVariant parser.OptionName = "HashtagVariant"
)

// Option is a configuration option for the Parser.
type Option interface {
	parser.Option
}

type withObsidianTags struct{}

func (o *withObsidianTags) SetParserOption(c *parser.Config) {
	c.Options[optVariant] = ObsidianVariant
}

// WithObsidianTags allows to parse tags in the Obsidian syntax.
func WithObsidianTags() Option {
	return &withObsidianTags{}
}
