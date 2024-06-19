package hashtag_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/hashtag"
	"gopkg.in/yaml.v3"
)

func TestIntegration_Default(t *testing.T) {
	t.Parallel()

	testIntegration(t, "default.yaml",
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{})))
}

func TestIntegration_Obsidian(t *testing.T) {
	t.Parallel()

	testIntegration(t, "obsidian.yaml",
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{
			Variant: hashtag.ObsidianVariant,
		})))
}

func TestIntegration_Resolver(t *testing.T) {
	t.Parallel()

	testIntegration(t, "resolver.yaml",
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{
			Resolver: almostAlwaysResolver{},
		})))
}

func TestIntegration_Attributes(t *testing.T) {
	t.Parallel()

	testIntegration(t, "attributes.yaml",
		goldmark.New(goldmark.WithExtensions(&hashtag.Extender{
			Resolver: almostAlwaysResolver{},
			Attributes: []hashtag.Attribute{
				{
					Attr:  "class",
					Value: "p-category",
				},
			},
		})))
}

func testIntegration(t *testing.T, file string, md goldmark.Markdown) {
	testsdata, err := os.ReadFile(filepath.Join("testdata", file))
	require.NoError(t, err)

	var tests []struct {
		Desc string `yaml:"desc"`
		Give string `yaml:"give"`
		Want string `yaml:"want"`
	}
	require.NoError(t, yaml.Unmarshal(testsdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &buf))
			require.Equal(t, tt.Want, buf.String())
		})
	}
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
