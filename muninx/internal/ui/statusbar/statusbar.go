// Package providing a statusbar for terminal apps.
package statusbar

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/lipgloss/v2"
)

type Partition int

// options for new input model;
type BarOptions func(m *Model)

const (
	Left Partition = iota
	Right
)

type Elem struct {
	tag     string //come with create, no edit;
	Content string
	Width   int
	BgColor *string
	FgColor *string
}

func colorPtr(c string) *string {
	return &c
}

func (e *Elem) SetValue(content string) *Elem {
	e.Content = content
	return e
}

func (e *Elem) SetColors(fg, bg *string) *Elem {
	e.FgColor = fg
	e.BgColor = bg
	return e
}

func (e *Elem) SetWidth(width int) *Elem {
	e.Width = width
	return e
}

//	func (e *Elem) SetPartition(p Partition) {
//		e.Partition = p
//	}
func (e Elem) Render(h int) string {
	style := lipgloss.NewStyle().
		Width(e.Width).
		Height(h).
		Align(lipgloss.Center)

	if e.FgColor != nil {
		style = style.Foreground(lipgloss.Color(*e.FgColor))
	}

	if e.BgColor != nil {
		style = style.Background(lipgloss.Color(*e.BgColor))
	}

	return style.Render(e.Content)
}

type Model struct {
	LeftElems  []*Elem
	RightElems []*Elem
	ElemsMap   map[string]*Elem //we use maps to keep a quick and convienent access to the elems;

	Height int
	Width  int
}

// New creates a new statusbar model
func New(ops ...BarOptions) Model {
	m := Model{
		LeftElems:  []*Elem{},
		RightElems: []*Elem{},
		ElemsMap:   make(map[string]*Elem),
		Height:     1,
		Width:      100,
	}
	for _, i := range ops {
		i(&m)
	}
	return m
}

// WithColumns sets the table columns (headers).
func WithWidth(w int) BarOptions {
	return func(m *Model) {
		m.Width = w
	}
}

func WithHeight(h int) BarOptions {
	return func(m *Model) {
		m.Height = h
	}
}

func WithLeftLen(n int) BarOptions {
	return func(m *Model) {
		// Create a new slice with the specified length
		m.LeftElems = make([]*Elem, n)

		// Initialize each element with default values
		for i := range n {
			m.LeftElems[i] = &Elem{
				tag:     "",
				Content: "",
				Width:   10, // Default width
				BgColor: colorPtr("236"),
				FgColor: colorPtr("252"),
			}
		}
	}
}

func WithRightLen(n int) BarOptions {
	return func(m *Model) {
		m.RightElems = make([]*Elem, n)
		for i := range n {
			m.RightElems[i] = &Elem{
				tag:     "",
				Content: "",
				Width:   10, // Default width
				BgColor: colorPtr("236"),
				FgColor: colorPtr("252"),
			}
		}
	}
}

// Set the tag for one elem, and accessable through the map;
func (m *Model) SetTag(e *Elem, tag string) {
	if e != nil {
		e.tag = tag
		m.ElemsMap[tag] = e
	}
}

func (m *Model) GetTag(tag string) *Elem {
	if m.ElemsMap == nil {
		return nil
	}

	// Look up the element by tag
	elem, exists := m.ElemsMap[tag]
	if !exists {
		return nil
	}
	return elem
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}
	return m, nil
}

func (m Model) View() tea.View {
	v := tea.NewView(m.Render())
	return v
}

// Add an element to the left partition of the bar.
func (m *Model) AddLeft(w int, c string) *Elem {
	newElem := &Elem{
		Content: c,
		Width:   w,
		BgColor: colorPtr("236"),
		FgColor: colorPtr("252"),
	}
	m.LeftElems = append(m.LeftElems, newElem)
	return newElem
}

// Remove by id from left partition.
func (m *Model) RemoveLeft(i int) *Model {
	if i >= 0 && i < len(m.LeftElems) {
		elem := m.LeftElems[i]
		if elem != nil && elem.tag != "" {
			delete(m.ElemsMap, elem.tag)
		}
		m.LeftElems = append(m.LeftElems[:i], m.LeftElems[i+1:]...)
	}
	return m
}

func (m *Model) RemoveRight(i int) *Model {
	if i >= 0 && i < len(m.RightElems) {
		elem := m.RightElems[i]
		if elem != nil && elem.tag != "" {
			delete(m.ElemsMap, elem.tag)
		}
		m.RightElems = append(m.RightElems[:i], m.RightElems[i+1:]...)
	}
	return m
}

func (m *Model) AddRight(w int, c string) *Elem {
	newElem := &Elem{
		Content: c,
		Width:   w,
		BgColor: colorPtr("236"),
		FgColor: colorPtr("252"),
	}
	m.RightElems = append(m.RightElems, newElem)
	return newElem
}

// Get element by id from left partition.
func (m *Model) GetLeft(index int) *Elem {
	if index >= 0 && index < len(m.LeftElems) {
		return m.LeftElems[index]
	}
	return nil
}

func (m *Model) GetRight(index int) *Elem {
	if index >= 0 && index < len(m.RightElems) {
		return m.RightElems[index]
	}
	return nil
}

func (m *Model) SetWidth(w int) *Model {
	m.Width = w
	return m
}

func (m *Model) SetHeight(h int) *Model {
	m.Height = h
	return m
}

// Render returns the statusbar as a string
func (m Model) Render() string {
	// Render left elements
	leftContent := ""
	for _, elem := range m.LeftElems {
		leftContent += elem.Render(m.Height)
	}

	// Render right elements
	rightContent := ""
	for _, elem := range m.RightElems {
		rightContent += elem.Render(m.Height)
	}

	// Calculate space between left and right elements
	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	middleWidth := max(0, m.Width-leftWidth-rightWidth)

	// Create empty middle space with appropriate width
	middleSpace := lipgloss.NewStyle().
		Width(middleWidth).
		Height(m.Height).
		Render("")

	// Join left content, middle space, and right content
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftContent,
		middleSpace,
		rightContent,
	)
}
