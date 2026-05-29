package menu

const (
	tableIDWidth   = 6
	tableTimeWidth = 16
	tablePad       = 8 // internal table column padding (bubbles table overhead)
	// header(3) + tableBox border(2) + inputBox(3) + help(2) = 10
	verticalOverhead = 10
)

// Layout holds all computed dimensions for MenuView.
type Layout struct {
	WindowWidth  int
	WindowHeight int

	TableIDWidth      int
	TableContentWidth int
	TableTimeWidth    int
	TableHeight       int
	TableWidth        int // inner table width (= WindowWidth - FocusedStyle frame 4)

	InputWidth int // outer width of the input box
}

func computeLayout(width, height int) Layout {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	// tableW is the table's own rendered width (FocusedStyle adds 4 around it).
	tableW := width - 4
	contentW := tableW - tableIDWidth - tableTimeWidth - tablePad
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
		TableTimeWidth:    tableTimeWidth,
		TableHeight:       tableH,
		TableWidth:        tableW,
		InputWidth:        width,
	}
}
