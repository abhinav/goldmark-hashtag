package hashtag

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestRenderer_Resolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		dest string
		want string
	}{
		{
			desc: "no destination",
			want: `<span class="hashtag">#foo</span>`,
		},
		{
			desc: "has destination",
			dest: "/bar",
			want: `<span class="hashtag"><a href="/bar">#foo</a></span>`,
		},
		{
			desc: "destination with spaces",
			dest: "/bar baz",
			want: `<span class="hashtag"><a href="/bar%20baz">#foo</a></span>`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			r := goldmark.New().Renderer()
			r.AddOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(&Renderer{
						Resolver: constResolver{
							Dest: tt.dest,
						},
					}, 999),
				),
			)

			src := []byte("#foo")
			node := &Node{Tag: src[1:]}
			node.AppendChild(node,
				ast.NewTextSegment(text.NewSegment(0, len(src))))

			var buff bytes.Buffer
			w := bufio.NewWriter(&buff)

			require.NoError(t, r.Render(w, src, node))
			assert.Equal(t, tt.want, buff.String())
		})
	}
}

func TestRender_WrongNode(t *testing.T) {
	t.Parallel()

	src := []byte("#foo")

	var r Renderer
	_, err := r.Render(
		bufio.NewWriter(new(bytes.Buffer)),
		src,
		ast.NewTextSegment(text.NewSegment(0, len(src))),
		true, // enter
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected node *ast.Text")
}

func TestRender_ResolveError(t *testing.T) {
	t.Parallel()

	giveErr := errors.New("great sadness")

	r := renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{
				Resolver: constResolver{
					Err: giveErr,
				},
			}, 999),
		),
	)
	src := []byte("#foo")
	node := &Node{Tag: src[1:]}
	node.AppendChild(node,
		ast.NewTextSegment(text.NewSegment(0, len(src))))

	var buff bytes.Buffer
	w := bufio.NewWriter(&buff)

	err := r.Render(w, src, node)
	require.Error(t, err)
	assert.ErrorIs(t, err, giveErr)
}

type constResolver struct {
	Dest string
	Err  error
}

func (r constResolver) ResolveHashtag(*Node) ([]byte, error) {
	return []byte(r.Dest), r.Err
}
