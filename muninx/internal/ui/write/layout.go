package write

import "github.com/haochend413/muninx/internal/ui/statusbar"

const (
	leftPanelRatio = 55 // percent of window width for the textarea panel
	minLeftWidth   = 20
	minRightWidth  = 20
)

// Layout holds computed dimensions for the write view's side-by-side panels.
type Layout struct {
	WindowWidth   int
	WindowHeight  int
	ContentHeight int // WindowHeight minus the bottom status bar
	TextAreaWidth int // passed to textarea.SetWidth; fills the left panel
	RelatedWidth  int // passed to viewport.SetWidth; fills the right panel
}

func computeLayout(width, height int) Layout {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	leftW := width * leftPanelRatio / 100
	rightW := width - leftW
	if leftW < minLeftWidth {
		leftW = minLeftWidth
	}
	if rightW < minRightWidth {
		rightW = minRightWidth
	}
	contentHeight := height - statusbar.DefaultHeight
	if contentHeight < 1 {
		contentHeight = 1
	}
	return Layout{
		WindowWidth:   width,
		WindowHeight:  height,
		ContentHeight: contentHeight,
		TextAreaWidth: leftW,
		RelatedWidth:  rightW,
	}
}
