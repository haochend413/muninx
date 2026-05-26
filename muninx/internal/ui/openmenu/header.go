package openmenu

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/haochend413/lipgloss/v2"
)

// This package defines the header of the opening menu.

type Opts struct {
	MainTextColor color.Color // diagonal lines
	Width         int         // width of the rendered logo, used for truncation
	Height        int
}

func Render(base lipgloss.Style, o Opts) string {
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
