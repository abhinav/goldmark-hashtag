package hashtag_test

import (
	"bytes"
	"testing"

	hashtag "github.com/abhinav/goldmark-hashtag"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestIntegration_Default(t *testing.T) {
	t.Parallel()

	testutil.DoTestCaseFile(
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{})),
		"testdata/default.txt",
		t,
	)
}
func TestIntegration_Obsidian(t *testing.T) {
	t.Parallel()

	testutil.DoTestCaseFile(
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{
			Variant: hashtag.ObsidianVariant,
		})),
		"testdata/obsidian.txt",
		t,
	)
}
func TestIntegration_Resolver(t *testing.T) {
	t.Parallel()

	testutil.DoTestCaseFile(
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{
			Resolver: almostAlwaysResolver{},
		})),
		"testdata/resolver.txt",
		t,
	)
}

var _unknownTag = []byte("unknown")

// Resolves all tags except "unknown".
type almostAlwaysResolver struct{}

func (almostAlwaysResolver) ResolveHashtag(n *hashtag.Node) ([]byte, error) {
	if bytes.Equal(n.Tag, _unknownTag) {
		return nil, nil
	}

	return append([]byte("/tag/"), n.Tag...), nil
}
