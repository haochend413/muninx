// Package statusbar provides a signal-driven status bar for terminal UIs.
//
// Usage:
//
//	bar := statusbar.New(termWidth, 1)
//	bar.Register("mode",  statusbar.Left,   statusbar.ElemConfig{Width: 8,  Fg: "0",   Bg: "34",  Align: statusbar.AlignCenter})
//	bar.Register("msg",   statusbar.Center, statusbar.ElemConfig{Width: 30, Fg: "252", Bg: "236", Align: statusbar.AlignCenter})
//	bar.Register("count", statusbar.Right,  statusbar.ElemConfig{Width: 20, Fg: "252", Bg: "236"})
//
//	bar.Set("mode", "INSERT")
//
//	// Flash a temporary notice — returns a cmd that auto-clears after 2 s:
//	cmd := bar.Signal("msg", statusbar.Signal{
//	    Content: "Saved!", Fg: "0", Bg: "40", Bold: true, Duration: 2 * time.Second,
//	})
//	return m, cmd
//
// Embed Update in the parent's Update so expiry messages are processed:
//
//	*m.bar, cmd = m.bar.Update(msg)
package statusbar

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/lipgloss/v2"
)

// DefaultHeight is the standard single-line height views should reserve for
// a status bar, so every view that embeds one stays visually consistent.
const DefaultHeight = 1

// Partition determines which section of the bar an element belongs to.
type Partition int

const (
	Left Partition = iota
	Center
	Right
)

// Align controls text alignment within an element cell.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// ElemConfig is the static configuration for a named slot.
type ElemConfig struct {
	Width int
	Fg    string // lipgloss foreground color string; "" = inherit terminal default
	Bg    string // lipgloss background color string
	Align Align
}

// Signal is a temporary visual override for one element.
// The element reverts to its base state after Duration.
// Duration == 0 means the signal persists until Clear is called.
type Signal struct {
	Content  string
	Fg       string // overrides ElemConfig.Fg when non-empty
	Bg       string // overrides ElemConfig.Bg when non-empty
	Bold     bool
	Duration time.Duration
}

// clearMsg is sent internally when a timed signal expires.
type clearMsg struct {
	bar *byte  // points to the bar that created this message; prevents cross-bar clearing
	tag string
	gen int // generation at Signal() time; stale if a newer signal was pushed
}

type elem struct {
	cfg       ElemConfig
	content   string
	signal    *Signal
	signalGen int
}

// Model is the statusbar.
type Model struct {
	id     *byte // unique per-instance; copied by value, so all copies of the same bar share the same id
	elems  map[string]*elem
	left   []string
	center []string
	right  []string
	width  int
	height int
}

// New creates an empty statusbar with the given width and height.
func New(width, height int) Model {
	return Model{
		id:     new(byte),
		elems:  make(map[string]*elem),
		width:  width,
		height: height,
	}
}

// SetWidth resizes the bar. Call this from a WindowSizeMsg handler.
func (m *Model) SetWidth(w int) { m.width = w }

// SetHeight resizes the bar height.
func (m *Model) SetHeight(h int) { m.height = h }

// SetElemWidth updates only the width of a named element, leaving its
// colors and alignment untouched. Useful for an element that should track
// the bar's full width across resizes.
func (m *Model) SetElemWidth(tag string, width int) {
	if e, ok := m.elems[tag]; ok {
		e.cfg.Width = width
	}
}

// Register adds a named slot to the given partition in declaration order.
// Calling Register with an already-registered tag is a no-op.
func (m *Model) Register(tag string, p Partition, cfg ElemConfig) {
	if _, exists := m.elems[tag]; exists {
		return
	}
	m.elems[tag] = &elem{cfg: cfg}
	switch p {
	case Left:
		m.left = append(m.left, tag)
	case Center:
		m.center = append(m.center, tag)
	default:
		m.right = append(m.right, tag)
	}
}

// Set updates the base (non-signal) content of a named element.
func (m *Model) Set(tag, content string) {
	if e, ok := m.elems[tag]; ok {
		e.content = content
	}
}

// SetColors updates the base foreground and background of a named element.
func (m *Model) SetColors(tag, fg, bg string) {
	if e, ok := m.elems[tag]; ok {
		e.cfg.Fg = fg
		e.cfg.Bg = bg
	}
}

// Configure replaces the full ElemConfig of a named element.
func (m *Model) Configure(tag string, cfg ElemConfig) {
	if e, ok := m.elems[tag]; ok {
		e.cfg = cfg
	}
}

// Signal pushes a temporary override onto a named element.
// Returns a Cmd that fires a clearMsg after s.Duration; returns nil when
// Duration == 0 (the signal stays until Clear is called explicitly).
// If Signal is called again before the previous one expires, the old timer
// becomes a no-op (generation counter mismatch).
func (m *Model) Signal(tag string, s Signal) tea.Cmd {
	e, ok := m.elems[tag]
	if !ok {
		return nil
	}
	e.signalGen++
	gen := e.signalGen
	e.signal = &s
	if s.Duration <= 0 {
		return nil
	}
	barID := m.id
	return tea.Tick(s.Duration, func(time.Time) tea.Msg {
		return clearMsg{bar: barID, tag: tag, gen: gen}
	})
}

// Clear removes any active signal from a named element immediately.
func (m *Model) Clear(tag string) {
	if e, ok := m.elems[tag]; ok {
		e.signal = nil
	}
}

// Init satisfies tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update handles signal-expiry and window-resize messages.
// Embed in the parent model's Update loop:
//
//	*m.bar, cmd = m.bar.Update(msg)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case clearMsg:
		if msg.bar == m.id {
			if e, ok := m.elems[msg.tag]; ok && e.signalGen == msg.gen {
				e.signal = nil
			}
		}
	}
	return m, nil
}

// Render returns the bar as a styled string ready to append below a view.
func (m Model) Render() string {
	leftStr := m.renderPartition(m.left)
	centerStr := m.renderPartition(m.center)
	rightStr := m.renderPartition(m.right)

	lw := lipgloss.Width(leftStr)
	cw := lipgloss.Width(centerStr)
	rw := lipgloss.Width(rightStr)

	fill := max(0, m.width-lw-cw-rw)
	lfill := fill / 2
	rfill := fill - lfill

	gap := lipgloss.NewStyle().Height(m.height).Background(lipgloss.Color("236"))
	return lipgloss.JoinHorizontal(lipgloss.Left,
		leftStr,
		gap.Width(lfill).Render(""),
		centerStr,
		gap.Width(rfill).Render(""),
		rightStr,
	)
}

func (m Model) renderPartition(tags []string) string {
	var sb strings.Builder
	for _, tag := range tags {
		sb.WriteString(m.renderElem(tag))
	}
	return sb.String()
}

func (m Model) renderElem(tag string) string {
	e, ok := m.elems[tag]
	if !ok {
		return ""
	}
	cfg := e.cfg
	content := e.content
	fg, bg := cfg.Fg, cfg.Bg
	bold := false

	if e.signal != nil {
		content = e.signal.Content
		if e.signal.Fg != "" {
			fg = e.signal.Fg
		}
		if e.signal.Bg != "" {
			bg = e.signal.Bg
		}
		bold = e.signal.Bold
	}

	align := lipgloss.Left
	switch cfg.Align {
	case AlignCenter:
		align = lipgloss.Center
	case AlignRight:
		align = lipgloss.Right
	}

	style := lipgloss.NewStyle().Width(cfg.Width).Height(m.height).Align(align)
	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if bold {
		style = style.Bold(true)
	}
	return style.Render(content)
}
