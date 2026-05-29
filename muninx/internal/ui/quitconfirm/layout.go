package quitconfirm

// Layout holds terminal dimensions for QuitConfirmView.
// The quit prompt is static text so no component sizing is needed,
// but the layout keeps the pattern consistent with other view models.
type Layout struct {
	WindowWidth  int
	WindowHeight int
}

func computeLayout(width, height int) Layout {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return Layout{WindowWidth: width, WindowHeight: height}
}
