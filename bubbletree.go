package bubbletree

import (
	"iter"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// Bubbletea Messages
type (
	reconstructMsg[T comparable]   Tree[T]
	reconstructedMsg[T comparable] []TreeLine[T]
	setFocusedIDMsg[T comparable]  struct{ value T }
)

type Model[T comparable] struct {
	UpdateFunc func(lines []TreeLine[T], focusedID T, msg tea.Msg) tea.Cmd

	currentLines []TreeLine[T]
	focusedID    T
}

type Tree[T comparable] interface {
	ID() T
	Content() string
	Children() iter.Seq2[Tree[T], bool]
}

type TreeLine[T comparable] struct {
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
			lines := constructTree(msg)
			return reconstructedMsg[T](lines)
		})

	case reconstructedMsg[T]:
		m.currentLines = msg
		var zero T
		if m.focusedID == zero && len(m.currentLines) > 0 {
			m.focusedID = m.currentLines[0].ID
		}

	case setFocusedIDMsg[T]:
		m.focusedID = msg.value

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			var ok bool
			for i, line := range m.currentLines {
				if ok {
					break
				}

				if line.ID == m.focusedID {
					cursor := (i + 1) % len(m.currentLines)
					m.focusedID = m.currentLines[cursor].ID
					ok = true
				}
			}

		case "k":
			var ok bool
			for i, line := range m.currentLines {
				if ok {
					break
				}

				if line.ID == m.focusedID {
					cursor := (i - 1 + len(m.currentLines)) % len(m.currentLines)
					m.focusedID = m.currentLines[cursor].ID

					break
				}
			}
		}
	}

	if m.currentLines != nil && m.UpdateFunc != nil {
		cmd := m.UpdateFunc(m.currentLines, m.focusedID, msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model[T]) View() string {
	b := strings.Builder{}
	for _, line := range m.currentLines {
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

func (m Model[T]) SetTree(tree Tree[T]) tea.Cmd {
	return func() tea.Msg {
		return reconstructMsg[T](tree)
	}
}

func (m Model[T]) SetFocusedID(focusedID T) tea.Cmd {
	return func() tea.Msg {
		return setFocusedIDMsg[T]{
			value: focusedID,
		}
	}
}

func constructTree[T comparable](tree Tree[T]) []TreeLine[T] {
	lines := make([]TreeLine[T], 0) // TODO: set appropriate cap
	lines = append(lines, TreeLine[T]{
		ID:  tree.ID(),
		Raw: tree.Content(),
	})
	childLines := constructChildren(tree, "")
	lines = append(lines, childLines...)

	return lines
}

func constructChildren[T comparable](tree Tree[T], prefix string) []TreeLine[T] {
	lines := make([]TreeLine[T], 0)
	b := strings.Builder{}
	for child, hasNext := range tree.Children() {
		connector := "├── "
		nextPrefix := "│   "
		if !hasNext {
			connector = "└── "
			nextPrefix = "    "
		}

		b.WriteString(prefix)
		b.WriteString(connector)
		b.WriteString(child.Content())

		lines = append(lines, TreeLine[T]{
			ID:  child.ID(),
			Raw: b.String(),
		})
		b.Reset()

		childLines := constructChildren(child, prefix+nextPrefix)
		lines = append(lines, childLines...)
	}

	return lines
}
