package openmenu

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	slice "github.com/charmbracelet/x/exp/slice"
	"github.com/haochend413/lipgloss/v2"
)

type letterform func(bool) string

func renderWord(spacing int, stretchIndex int, letterforms ...letterform) string {
	if spacing < 0 {
		spacing = 0
	}

	rendered := make([]string, len(letterforms))
	for i, letter := range letterforms {
		rendered[i] = letter(i == stretchIndex)
	}

	if spacing > 0 {
		rendered = slice.Intersperse(rendered, strings.Repeat(" ", spacing))
	}

	return strings.TrimSpace(
		lipgloss.JoinHorizontal(lipgloss.Top, rendered...),
	)
}

func LetterM(stretch bool) string {
	left := heredoc.Doc(`
		█
		█
		▀
	`)

	innerLeft := heredoc.Doc(`
		▄
		█
		 
	`)

	innerRight := heredoc.Doc(`
		▄
		█
		 
	`)

	right := heredoc.Doc(`
		█
		█
		▀
	`)

	return joinLetterform(
		left,
		stretchLetterformPart(innerLeft, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 4,
			maxStretch: 8,
		}),
		stretchLetterformPart(innerRight, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 4,
			maxStretch: 8,
		}),
		right,
	)
}

func LetterU(stretch bool) string {
	side := heredoc.Doc(`
		█
		█
		▀
	`)

	middle := heredoc.Doc(`


		▀
	`)

	return joinLetterform(
		side,
		stretchLetterformPart(middle, letterformProps{
			stretch:    stretch,
			width:      5,
			minStretch: 10,
			maxStretch: 18,
		}),
		side,
	)
}

func LetterN(stretch bool) string {
	left := heredoc.Doc(`
		█
		█
		▀
	`)

	middle := heredoc.Doc(`
		▄
		█
		▀
	`)

	right := heredoc.Doc(`
		█
		█
		▀
	`)

	return joinLetterform(
		left,
		stretchLetterformPart(middle, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 6,
			maxStretch: 12,
		}),
		right,
	)
}

func LetterI(stretch bool) string {
	left := heredoc.Doc(`
		▀
		 
		▀
	`)

	center := heredoc.Doc(`
		█
		█
		▀
	`)

	right := heredoc.Doc(`
		▀
		 
		▀
	`)

	return joinLetterform(
		stretchLetterformPart(left, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 4,
			maxStretch: 8,
		}),
		center,
		stretchLetterformPart(right, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 4,
			maxStretch: 8,
		}),
	)
}

func LetterX(stretch bool) string {
	left := heredoc.Doc(`
		█
		 
		▀
	`)

	middle := heredoc.Doc(`
		▄
		█
		▀
	`)

	right := heredoc.Doc(`
		█
		 
		▀
	`)

	return joinLetterform(
		left,
		stretchLetterformPart(middle, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 6,
			maxStretch: 12,
		}),
		right,
	)
}

func RenderMuninx() string {
	return renderWord(2, -1, LetterM, LetterU, LetterN, LetterI, LetterN, LetterX)
}

func joinLetterform(parts ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

type letterformProps struct {
	width      int
	minStretch int
	maxStretch int
	stretch    bool
}

func stretchLetterformPart(s string, p letterformProps) string {
	if p.maxStretch < p.minStretch {
		p.minStretch, p.maxStretch = p.maxStretch, p.minStretch
	}

	n := p.width
	if p.stretch {
		n = p.minStretch
	}

	parts := make([]string, n)
	for i := range parts {
		parts[i] = s
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
