package hashtag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestVariantSpan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give         string // input string
		wantDefault  string // empty for invalid
		wantObsidian string // empty for invalid
	}{
		{give: "", wantDefault: "", wantObsidian: ""},
		{give: "1a", wantDefault: "", wantObsidian: "1a"},
		{give: "a", wantDefault: "a", wantObsidian: "a"},
		{give: "a1", wantDefault: "a1", wantObsidian: "a1"},
		{give: "foo", wantDefault: "foo", wantObsidian: "foo"},
		{give: "foo bar", wantDefault: "foo", wantObsidian: "foo"},
		{give: "éabc d", wantDefault: "éabc", wantObsidian: "éabc"},
		{give: "✅/🚧", wantDefault: "", wantObsidian: "✅/🚧"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()

			t.Run("default", func(t *testing.T) {
				t.Parallel()

				var got string
				if idx := DefaultVariant.span([]byte(tt.give)); idx >= 0 {
					got = tt.give[:idx]
				}
				assert.Equal(t, tt.wantDefault, got)
			})

			t.Run("obsidian", func(t *testing.T) {
				t.Parallel()

				var got string
				if idx := ObsidianVariant.span([]byte(tt.give)); idx >= 0 {
					got = tt.give[:idx]
				}
				assert.Equal(t, tt.wantObsidian, got)
			})
		})
	}
}

func TestParser(t *testing.T) {
	t.Parallel()

	type node struct {
		Tag  string
		Body string
	}

	tests := []struct {
		desc      string
		give      string
		want      *node
		remaining string
		variant   Variant
	}{
		{
			desc:      "empty",
			give:      "",
			remaining: "",
		},
		{
			desc:      "not hash",
			give:      "foo",
			remaining: "foo",
		},
		{
			desc:      "hash alone",
			give:      "# foo",
			remaining: "# foo",
		},
		{
			desc:      "starts with number",
			give:      "#1foo",
			remaining: "#1foo",
		},
		{
			desc: "simple tag",
			give: "#foo bar",
			want: &node{
				Tag:  "foo",
				Body: "#foo",
			},
			remaining: " bar",
		},
		{
			desc: "line ends with tag",
			give: "#foo\n#bar",
			want: &node{
				Tag:  "foo",
				Body: "#foo",
			},
			remaining: "\n#bar",
		},
		{
			desc: "hyphen",
			give: "#foo-bar",
			want: &node{
				Tag:  "foo-bar",
				Body: "#foo-bar",
			},
		},
		{
			desc: "underscore",
			give: "#foo_bar",
			want: &node{
				Tag:  "foo_bar",
				Body: "#foo_bar",
			},
		},
		{
			desc: "slash",
			give: "#foo/bar",
			want: &node{
				Tag:  "foo/bar",
				Body: "#foo/bar",
			},
		},
		{
			desc:      "obsidian hash alone",
			give:      "# foo",
			remaining: "# foo",
			variant:   ObsidianVariant,
		},
		{
			desc: "obsidian start",
			give: "#123tag",
			want: &node{
				Tag:  "123tag",
				Body: "#123tag",
			},
			variant: ObsidianVariant,
		},
		{
			desc: "obsidian deny symbols",
			give: "#tag%tag",
			want: &node{
				Tag:  "tag",
				Body: "#tag",
			},
			remaining: "%tag",
			variant:   ObsidianVariant,
		},
		{
			desc: "obsidian accept underscore",
			give: "#asd_123",
			want: &node{
				Tag:  "asd_123",
				Body: "#asd_123",
			},
			variant: ObsidianVariant,
		},
		{
			desc: "obsidian accept dash",
			give: "#asd-123",
			want: &node{
				Tag:  "asd-123",
				Body: "#asd-123",
			},
			variant: ObsidianVariant,
		},
		{
			desc: "obsidian accept forward slash",
			give: "#asd/123",
			want: &node{
				Tag:  "asd/123",
				Body: "#asd/123",
			},
			variant: ObsidianVariant,
		},
		{
			desc:      "obsidian not all digits",
			give:      "#123",
			remaining: "#123",
			variant:   ObsidianVariant,
		},
		{
			desc: "obsidian ends line",
			give: "#foo\nbar",
			want: &node{
				Tag:  "foo",
				Body: "#foo",
			},
			remaining: "\nbar",
			variant:   ObsidianVariant,
		},
		{
			desc: "obsidian digits and symbol",
			give: "#321/123",
			want: &node{
				Tag:  "321/123",
				Body: "#321/123",
			},
			variant: ObsidianVariant,
		},
		{
			desc: "obsidian accept emojis",
			give: "#✅/🚧",
			want: &node{
				Tag:  "✅/🚧",
				Body: "#✅/🚧",
			},
			variant: ObsidianVariant,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			src := []byte(tt.give)
			rdr := text.NewReader(src)

			p := Parser{Variant: tt.variant}
			got := p.Parse(nil /* parent */, rdr, parser.NewContext())

			if tt.want != nil {
				require.IsType(t, &Node{}, got)
				got := got.(*Node)
				assert.Equal(t, *tt.want, node{
					Tag:  string(got.Tag),
					Body: string(got.Text(src)),
				})
			} else {
				assert.Nil(t, got)
			}

			_, pos := rdr.Position()
			assert.Equal(t, tt.remaining, string(src[pos.Start:]))
		})
	}
}
