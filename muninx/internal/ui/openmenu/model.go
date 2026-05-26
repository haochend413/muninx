package openmenu

import (
	"github.com/haochend413/bubbles/v2/textinput"
	"github.com/haochend413/muninx/internal/ui/table"
)

type MenuModel struct {
	IndexSelect table.Model
	Input       textinput.Model
}

func NewMenuModel() {
	cols := []table.Column{
		{Title: "ID", Width: 1},
		{Title: "Content", Width: 1},
		{Title: "Last Updated", Width: 1},
	}

	_ = table.New(
		table.WithColumns(cols),
		table.WithWidth(1),
	)
}
