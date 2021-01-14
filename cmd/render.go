package cmd

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

// render takes a table writer and returns a string from requested renderer
// currently only supports the default terminal output or Markdown
func render(t table.Writer) string {
	if markdownOutput {
		return t.RenderMarkdown()
	}
	return t.Render()
}
