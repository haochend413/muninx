package styles

import "github.com/haochend413/bubbles/v2/table"

// Layout carries all derived dimensions used by the UI.
type Layout struct {
	WindowWidth           int
	WindowHeight          int
	ViewMainContentHeight int
	MainContentHeight     int

	TableWidth  int
	EditorWidth int

	NotesBaseHeight int

	TextAreaHeight    int
	ChangeTableHeight int

	NoteColumns   []table.Column
	ChangeColumns []table.Column
}

const (
	borderOverhead          = 8
	tableWidthRatio         = 0.35
	idWidthRatio            = 0.08
	timeWidthRatio          = 0.22
	flagWidthRatio          = 0.15
	contentWidthRatio       = 0.58
	changeTypeWidthRatio    = 0.10
	changeIDWidthRatio      = 0.10
	changeTimeWidthRatio    = 0.25
	changeDescWidthRatio    = 0.40
	statusReservedHeight    = 5
	viewReservedHeight      = 3
	tableMathOffset         = 3
	notesBaseMultiplier     = 8
	notesBaseOffset         = 1
	notesFocusedExtraHeight = 6
	notesNormalExtraHeight  = 2
)

func ComputeLayout(windowWidth, windowHeight int) Layout {
	if windowWidth < 1 {
		windowWidth = 1
	}
	if windowHeight < 1 {
		windowHeight = 1
	}

	availableWidth := windowWidth - borderOverhead
	if availableWidth < 1 {
		availableWidth = 1
	}

	tableWidth := int(float64(availableWidth) * tableWidthRatio)
	if tableWidth < 1 {
		tableWidth = 1
	}

	editorWidth := availableWidth - tableWidth
	if editorWidth < 1 {
		editorWidth = 1
	}

	idWidth := max(4, int(float64(tableWidth)*idWidthRatio))
	timeWidth := max(8, int(float64(tableWidth)*timeWidthRatio))
	flagWidth := max(4, int(float64(tableWidth)*flagWidthRatio))
	contentWidth := max(10, int(float64(tableWidth)*contentWidthRatio))

	mainContentHeight := windowHeight - statusReservedHeight
	viewMainContentHeight := windowHeight - viewReservedHeight
	tableHeight := max(3, (mainContentHeight-tableMathOffset)/10)
	notesBaseHeight := tableHeight*notesBaseMultiplier + notesBaseOffset

	textareaHeight := max(5, mainContentHeight) - 1
	changeTableHeight := max(5, int(float64(mainContentHeight)*0.3)) - 3

	noteColumns := []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Content", Width: contentWidth},
		{Title: "Flags", Width: flagWidth},
	}

	changeColumns := []table.Column{
		{Title: "Type", Width: max(6, int(float64(editorWidth)*changeTypeWidthRatio))},
		{Title: "ID", Width: max(4, int(float64(editorWidth)*changeIDWidthRatio))},
		{Title: "Time", Width: max(12, int(float64(editorWidth)*changeTimeWidthRatio))},
		{Title: "Description", Width: max(15, int(float64(editorWidth)*changeDescWidthRatio))},
	}

	return Layout{
		WindowWidth:           windowWidth,
		WindowHeight:          windowHeight,
		ViewMainContentHeight: viewMainContentHeight,
		MainContentHeight:     mainContentHeight,
		TableWidth:            tableWidth,
		EditorWidth:           editorWidth,
		NotesBaseHeight:       notesBaseHeight,
		TextAreaHeight:        textareaHeight,
		ChangeTableHeight:     changeTableHeight,
		NoteColumns:           noteColumns,
		ChangeColumns:         changeColumns,
	}
}

func NotesFocusedHeight(layout Layout) int {
	return layout.NotesBaseHeight + notesFocusedExtraHeight
}

func NotesNormalHeight(layout Layout) int {
	return layout.NotesBaseHeight + notesNormalExtraHeight
}
