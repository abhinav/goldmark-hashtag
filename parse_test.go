package hashtag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			src := []byte(tt.give)
			rdr := text.NewReader(src)

			var p Parser
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
