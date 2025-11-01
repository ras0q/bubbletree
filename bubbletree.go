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
)

type Model[T comparable] struct {
	UpdateFunc func(line TreeLine[T], msg tea.Msg) tea.Cmd

	currentLines   []TreeLine[T]
	focusedLineNum int
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

	tea.Println(msg)
	switch msg := msg.(type) {
	case reconstructMsg[T]:
		cmds = append(cmds, func() tea.Msg {
			lines := constructTree(msg)
			return reconstructedMsg[T](lines)
		})

	case reconstructedMsg[T]:
		m.currentLines = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			l := len(m.currentLines)
			m.focusedLineNum = (l + m.focusedLineNum + 1) % l
		case "k":
			l := len(m.currentLines)
			m.focusedLineNum = (l + m.focusedLineNum - 1) % l
		}
	}

	if m.currentLines != nil && m.UpdateFunc != nil {
		cmd := m.UpdateFunc(m.currentLines[m.focusedLineNum], msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model[T]) View() string {
	b := strings.Builder{}
	for i, line := range m.currentLines {
		style := lipgloss.NewStyle().Width(50) // TODO: set appropriate width
		if i == m.focusedLineNum {
			style = style.Background(lipgloss.Color("205")).Bold(true)
		}
		b.WriteString(style.Render(line.Raw))
		b.WriteString("\n")
	}

	return b.String()
}

// MARK: Helper methods

func (m *Model[T]) SetTree(tree Tree[T]) tea.Cmd {
	return func() tea.Msg {
		return reconstructMsg[T](tree)
	}
}

func constructTree[T comparable](tree Tree[T]) []TreeLine[T] {
	lines := make([]TreeLine[T], 0) // TODO: set appropriate cap
	lines = append(lines, TreeLine[T]{
		ID:  tree.ID(),
		Raw: tree.Content(),
	})
	childLines := constructChildren(tree.Children(), "")
	lines = append(lines, childLines...)

	return lines
}

func constructChildren[T comparable](children iter.Seq2[Tree[T], bool], prefix string) []TreeLine[T] {
	lines := make([]TreeLine[T], 0)
	b := strings.Builder{}
	for child, hasNext := range children {
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

		childLines := constructChildren(child.Children(), nextPrefix)
		lines = append(lines, childLines...)
	}

	return lines
}
