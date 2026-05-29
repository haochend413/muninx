package write

const (
	relatedPanelRatio = 40 // percent of window width given to the related-notes panel
	panelPad          = 4  // border (2) + padding (2) per panel
	verticalPad       = 2  // help line (2 rows: 1 top-padding + 1 content)
)

// Layout holds all computed dimensions for WriteView.
type Layout struct {
	WindowWidth  int
	WindowHeight int

	// Outer widths passed to lipgloss style.Width.
	RelatedWidth int
	EditorWidth  int

	// Inner dimensions passed to component SetWidth/SetHeight.
	ListWidth  int
	ListHeight int

	TextAreaWidth  int
	TextAreaHeight int

	// Height used for the lipgloss box Height call.
	InnerHeight int
}

func computeLayout(width, height int) Layout {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	relatedW := (width * relatedPanelRatio) / 100
	editorW := width - relatedW

	// outerH is the total Height() passed to the lipgloss box (includes border).
	// contentH is the inner height passed to component SetHeight (outer minus border 2).
	outerH := height - verticalPad
	if outerH < 3 {
		outerH = 3
	}
	contentH := outerH - 2
	if contentH < 1 {
		contentH = 1
	}

	listW := relatedW - panelPad
	if listW < 10 {
		listW = 10
	}

	taW := editorW - panelPad
	if taW < 10 {
		taW = 10
	}

	// textarea_vim.View() appends a statusbar row after the viewport, so
	// SetHeight(n) produces n+1 rendered rows. Subtract 1 so the textarea
	// fits inside the same contentH inner box height as the list.
	taH := contentH - 1
	if taH < 1 {
		taH = 1
	}

	return Layout{
		WindowWidth:    width,
		WindowHeight:   height,
		RelatedWidth:   relatedW,
		EditorWidth:    editorW,
		ListWidth:      listW,
		ListHeight:     contentH,
		TextAreaWidth:  taW,
		TextAreaHeight: taH,
		InnerHeight:    outerH,
	}
}
