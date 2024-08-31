package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMdToHtml(t *testing.T) {
	md := `test\ntest\n`
	html := MdToHtml([]byte(md))

	assert.Equal(t, "<p>test\\ntest\\n</p>\n", string(html))
}
