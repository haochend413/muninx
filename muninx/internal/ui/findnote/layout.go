package findnote

import bTable "github.com/haochend413/bubbles/v2/table"

const (
	overlayMarginX    = 4
	overlayMarginY    = 3
	overlayMinWidth   = 60
	overlayMinHeight  = 10
	tableMinHeight    = 3
	viewportMinWidth  = 10
	viewportMaxWidth  = 45
)

// Layout holds all computed dimensions for the FindNote overlay.
type Layout struct {
	WindowWidth  int
	WindowHeight int

	OverlayWidth  int
	OverlayHeight int

	TableInnerWidth int
	TableHeight     int
	ViewportWidth   int

	TableColumns []bTable.Column
}

func computeLayout(width, height int) Layout {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	overlayW := width - 2*overlayMarginX
	if overlayW < overlayMinWidth {
		overlayW = overlayMinWidth
	}
	overlayH := height - 2*overlayMarginY
	if overlayH < overlayMinHeight {
		overlayH = overlayMinHeight
	}

	tableBoxW := tableInnerW + 4
	vpW := overlayW - 3*tableBoxW - 4
	if vpW > viewportMaxWidth {
		vpW = viewportMaxWidth
	}
	if vpW < viewportMinWidth {
		vpW = viewportMinWidth
	}
	tableH := overlayH - 4
	if tableH < tableMinHeight {
		tableH = tableMinHeight
	}

	cols := []bTable.Column{
		{Title: "ID", Width: 3},
		{Title: "Name", Width: tableInnerW - 4},
	}

	return Layout{
		WindowWidth:     width,
		WindowHeight:    height,
		OverlayWidth:    overlayW,
		OverlayHeight:   overlayH,
		TableInnerWidth: tableInnerW,
		TableHeight:     tableH,
		ViewportWidth:   vpW,
		TableColumns:    cols,
	}
}
