package menu

const (
	tableIDWidth   = 6
	tableTimeWidth = 16
	tableColumnGap = 2 // blank spacer column rendered between Content and Last Edited
	// header(4) + inputBox(3) + statusbar(1) = 8
	// (tableBox renders with all border sides disabled, so it adds 0.)
	verticalOverhead = 8
)

// Layout holds all computed dimensions for MenuView.
type Layout struct {
	WindowWidth  int
	WindowHeight int

	TableIDWidth      int
	TableContentWidth int
	TableColumnGap    int
	TableTimeWidth    int
	TableHeight       int
	TableWidth        int // inner table width (= WindowWidth - FocusedStyle frame 2)

	InputWidth int // outer width of the input box
}

func computeLayout(width, height int) Layout {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	// tableW is the table's own rendered width. FocusedStyle renders with no
	// border (disabled in view.go) and Padding(0, 1), so it adds 2 (1 each side).
	tableW := width - 2
	contentW := tableW - tableIDWidth - tableTimeWidth - tableColumnGap
	if contentW < 10 {
		contentW = 10
	}

	tableH := height - verticalOverhead
	if tableH < 3 {
		tableH = 3
	}

	return Layout{
		WindowWidth:       width,
		WindowHeight:      height,
		TableIDWidth:      tableIDWidth,
		TableContentWidth: contentW,
		TableColumnGap:    tableColumnGap,
		TableTimeWidth:    tableTimeWidth,
		TableHeight:       tableH,
		TableWidth:        tableW,
		InputWidth:        width,
	}
}
