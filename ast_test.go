package hashtag

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// /nolint:paralleltest // hijacks os.Stdout
func TestNodeDump(t *testing.T) {
	stdoutPath := filepath.Join(t.TempDir(), "stdout")
	stdout, err := os.Create(stdoutPath)
	require.NoError(t, err)

	defer func(stdout *os.File) { os.Stdout = stdout }(os.Stdout)
	os.Stdout = stdout

	src := []byte("#foo")
	node := &Node{Tag: src[1:]}
	node.AppendChild(node, ast.NewTextSegment(text.NewSegment(0, len(src))))

	node.Dump(src, 0)

	require.NoError(t, stdout.Close())

	got, err := os.ReadFile(stdoutPath)
	require.NoError(t, err)
	assert.Equal(t, strings.Join([]string{
		"Hashtag {",
		`    Tag: foo`,
		`    Text: "#foo"`,
		"}",
		"",
	}, "\n"), string(got))
}
