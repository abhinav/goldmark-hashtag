package hashtag

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Parser is a Goldmark inline parser for parsing hashtag nodes.
//
// Hashtags start with # and an alphabet, followed by zero or more
// alphanumeric characters and the following symbols.
//
//  /_-
type Parser struct {
}

var _ parser.InlineParser = (*Parser)(nil)

var _hash = byte('#')

// Trigger reports characters that trigger this parser.
func (*Parser) Trigger() []byte {
	return []byte{_hash}
}

// Parse parses a hashtag node.
func (*Parser) Parse(parent ast.Node, block text.Reader, pctx parser.Context) ast.Node {
	line, seg := block.PeekLine()

	if len(line) == 0 || line[0] != _hash {
		return nil
	}
	line = line[1:]

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

	n := Node{
		Tag: block.Value(seg.WithStart(seg.Start + 1)), // omit the "#"
	}
	n.AppendChild(&n, ast.NewTextSegment(seg))
	block.Advance(seg.Len())
	return &n
}

func endOfHashtag(r rune) bool {
	return !(unicode.IsLetter(r) ||
		unicode.IsDigit(r) ||
		r == '_' || r == '-' || r == '/')
}
