package bubbletree

import (
	"iter"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/v2"
)

// Bubbletea Messages
type (
	reconstructMsg[T comparable]   Node[T]
	reconstructedMsg[T comparable] []RenderedLine[T]
	setFocusMsg[T comparable]      struct{ id T }
)

// Bubbletea Model
type Model[T comparable] struct {
	OnUpdate func(lines []RenderedLine[T], focusedID T, msg tea.Msg) tea.Cmd

	renderedLines []RenderedLine[T]
	focusedID     T
}

type Node[T comparable] interface {
	ID() T
	Content() string
	Children() iter.Seq2[Node[T], bool]
}

type RenderedLine[T comparable] struct {
	ID  T
	Raw string
}

func New[T comparable]() Model[T] {
	return Model[T]{}
}

// MARK: Elm architecture implementation

func (m Model[T]) Init() tea.Cmd {
	return nil
}

func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case reconstructMsg[T]:
		cmds = append(cmds, func() tea.Msg {
			lines := renderTree(msg)
			return reconstructedMsg[T](lines)
		})

	case reconstructedMsg[T]:
		m.renderedLines = msg
		var zero T
		if m.focusedID == zero && len(m.renderedLines) > 0 {
			m.focusedID = m.renderedLines[0].ID
		}

	case setFocusMsg[T]:
		m.focusedID = msg.id

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			var ok bool
			for i, line := range m.renderedLines {
				if ok {
					break
				}

				if line.ID == m.focusedID {
					cursor := (i + 1) % len(m.renderedLines)
					m.focusedID = m.renderedLines[cursor].ID
					ok = true
				}
			}

		case "k":
			var ok bool
			for i, line := range m.renderedLines {
				if ok {
					break
				}

				if line.ID == m.focusedID {
					cursor := (i - 1 + len(m.renderedLines)) % len(m.renderedLines)
					m.focusedID = m.renderedLines[cursor].ID

					break
				}
			}
		}
	}

	if m.renderedLines != nil && m.OnUpdate != nil {
		cmd := m.OnUpdate(m.renderedLines, m.focusedID, msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model[T]) View() string {
	b := strings.Builder{}
	for _, line := range m.renderedLines {
		style := lipgloss.NewStyle().Width(50) // TODO: set appropriate width
		if line.ID == m.focusedID {
			style = style.Background(lipgloss.Color("205")).Bold(true)
		}
		b.WriteString(style.Render(line.Raw))
		b.WriteString("\n")
	}

	return b.String()
}

// MARK: Helper methods

func (m Model[T]) SetTree(node Node[T]) tea.Cmd {
	return func() tea.Msg {
		return reconstructMsg[T](node)
	}
}

func (m Model[T]) SetFocusedID(focusedID T) tea.Cmd {
	return func() tea.Msg {
		return setFocusMsg[T]{
			id: focusedID,
		}
	}
}

func renderTree[T comparable](node Node[T]) []RenderedLine[T] {
	lines := make([]RenderedLine[T], 0) // TODO: set appropriate cap
	lines = append(lines, RenderedLine[T]{
		ID:  node.ID(),
		Raw: node.Content(),
	})
	childLines := renderChildren(node, "")
	lines = append(lines, childLines...)

	return lines
}

func renderChildren[T comparable](node Node[T], prefix string) []RenderedLine[T] {
	lines := make([]RenderedLine[T], 0)
	b := strings.Builder{}
	for child, hasNext := range node.Children() {
		connector := "├── "
		nextPrefix := "│   "
		if !hasNext {
			connector = "└── "
			nextPrefix = "    "
		}

		b.WriteString(prefix)
		b.WriteString(connector)
		b.WriteString(child.Content())

		lines = append(lines, RenderedLine[T]{
			ID:  child.ID(),
			Raw: b.String(),
		})
		b.Reset()

		childLines := renderChildren(child, prefix+nextPrefix)
		lines = append(lines, childLines...)
	}

	return lines
}
