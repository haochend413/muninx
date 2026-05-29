package menu

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/haochend413/lipgloss/v2"
)

// HeaderOpts controls how the ASCII logo header is rendered.
type HeaderOpts struct {
	MainTextColor color.Color
	Width         int
	Height        int
}

func renderHeader(base lipgloss.Style, o HeaderOpts) string {
	text := RenderMuninx()

	style := base.Foreground(o.MainTextColor).Border(lipgloss.BlockBorder(), false, false, false, false)

	if o.Width > 0 {
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		text = strings.Join(lines, "\n")
	}

	return style.Render(text)
}
