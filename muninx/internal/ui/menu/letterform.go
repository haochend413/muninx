package menu

import (
	"strings"

	slice "github.com/charmbracelet/x/exp/slice"
	"github.com/haochend413/lipgloss/v2"
)

// letterform renders one letter given an extra column budget to spend on its
// flexible strokes, on top of its natural (extra == 0) width. Each letter is
// built directly as three row strings — not by gluing together pre-rendered
// column blocks — so stretching is just repeating a row's middle character
// more times.
type letterform func(extra int) string

// letter pairs a letterform with how many extra columns it can absorb, so
// renderWord can plan a distribution before rendering anything.
type letter struct {
	render   letterform
	capacity int
}

func newLetter(render letterform, capacity int) letter {
	return letter{render: render, capacity: capacity}
}

// renderWord lays out letters left to right with spacing columns between
// them, stretching their flexible strokes just enough (and as evenly as
// their capacities allow) to reach targetWidth. If targetWidth is already
// met by the natural layout, nothing is stretched.
func renderWord(spacing int, targetWidth int, letters ...letter) string {
	if spacing < 0 {
		spacing = 0
	}

	rendered := make([]string, len(letters))
	for i, l := range letters {
		rendered[i] = l.render(0)
	}

	natural := 0
	for i, r := range rendered {
		natural += lipgloss.Width(r)
		if i > 0 {
			natural += spacing
		}
	}

	if needed := targetWidth - natural; needed > 0 {
		extras := distributeExtra(needed, letters)
		for i, l := range letters {
			if extras[i] > 0 {
				rendered[i] = l.render(extras[i])
			}
		}
	}

	if spacing > 0 {
		rendered = slice.Intersperse(rendered, strings.Repeat(" ", spacing))
	}

	return strings.TrimSpace(
		lipgloss.JoinHorizontal(lipgloss.Top, rendered...),
	)
}

// distributeExtra spreads needed extra columns across letters in proportion
// to their capacity (so big, flexible letters like U absorb more than tight
// ones like N), then hands out any leftover from flooring one column at a
// time, round-robin, to whichever letters still have room.
func distributeExtra(needed int, letters []letter) []int {
	extras := make([]int, len(letters))

	totalCapacity := 0
	for _, l := range letters {
		totalCapacity += l.capacity
	}
	if totalCapacity == 0 {
		return extras
	}
	if needed > totalCapacity {
		needed = totalCapacity
	}

	remaining := needed
	for i, l := range letters {
		share := min(needed*l.capacity/totalCapacity, l.capacity)
		extras[i] = share
		remaining -= share
	}

	for i := 0; remaining > 0; i = (i + 1) % len(letters) {
		if extras[i] < letters[i].capacity {
			extras[i]++
			remaining--
		}
	}

	return extras
}

// rows joins multiple lines into the multi-line string a letterform returns.
func rows(lines ...string) string {
	return strings.Join(lines, "\n")
}

// LetterM is a fixed 7x4 glyph; it doesn't use extra (capacity 0 in
// RenderMuninx).
func LetterM(extra int) string {
	return rows(
		"█▄   ▄█",
		"██▄ ▄██",
		"█ ▀█▀ █",
		"▀     ▀",
	)
}

// LetterU is a fixed 7x4 glyph; it doesn't use extra (capacity 0 in
// RenderMuninx).
func LetterU(extra int) string {
	return rows(
		"█     █",
		"█     █",
		"▀▄   ▄▀",
		"  ▀▀▀  ",
	)
}

// LetterN is a fixed 6x4 glyph; it doesn't use extra (capacity 0 in
// RenderMuninx).
func LetterN(extra int) string {
	return rows(
		"█▄   █",
		"█ █▄ █",
		"█  █▄█",
		"▀   ▀▀",
	)
}

// LetterI is a fixed 1x4 glyph; it doesn't use extra (capacity 0 in
// RenderMuninx).
func LetterI(extra int) string {
	return rows(
		"▀▀█▀▀",
		"  █ ",
		"  █  ",
		"▀▀▀▀▀",
	)
}

// LetterX is a fixed 5x4 glyph; it doesn't use extra (capacity 0 in
// RenderMuninx).
func LetterX(extra int) string {
	return rows(
		"█   █",
		" ▀▄▀ ",
		" ▄▀▄ ",
		"▀   ▀",
	)
}

// RenderMuninx renders the MUNINX wordmark, stretching its letters just
// enough to reach targetWidth (use 0 for the natural, unstretched size).
// None of the current letterforms have stretch capacity, so this renders at
// its fixed natural size regardless of targetWidth.
func RenderMuninx(targetWidth int) string {
	return renderWord(2, targetWidth,
		newLetter(LetterM, 0),
		newLetter(LetterU, 0),
		newLetter(LetterN, 0),
		newLetter(LetterI, 0),
		newLetter(LetterN, 0),
		newLetter(LetterX, 0),
	)
}
