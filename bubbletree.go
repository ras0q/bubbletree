package bubbletree

import (
	"iter"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Model[T comparable] struct {
	ItemTree ItemTree[T]

	currentLines []TreeLine[T]
	focusedLine  int
}

func New[T comparable](tree ItemTree[T]) Model[T] {
	return Model[T]{
		ItemTree:    tree,
		focusedLine: 0,
	}
}

func (m Model[T]) Init() tea.Cmd {
	return func() tea.Msg {
		lines := constructTree(m.ItemTree)
		return TreeLinesConstructedMsg[T](lines)
	}
}

func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	switch msg := msg.(type) {
	case TreeLinesConstructedMsg[T]:
		m.currentLines = msg

	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			l := len(m.currentLines)
			m.focusedLine = (l + m.focusedLine + 1) % l
		case "k":
			l := len(m.currentLines)
			m.focusedLine = (l + m.focusedLine - 1) % l
		}
	}

	return m, nil
}

func (m Model[T]) View() string {
	b := strings.Builder{}
	for i, line := range m.currentLines {
		style := lipgloss.NewStyle().Width(50) // TODO: set appropriate width
		if i == m.focusedLine {
			style = style.Background(lipgloss.Color("205")).Bold(true)
		}
		b.WriteString(style.Render(line.Raw))
		b.WriteString("\n")
	}

	return b.String()
}

type ItemTree[T comparable] interface {
	ID() T
	Content() string
	Children() iter.Seq2[ItemTree[T], bool]
}

type TreeLine[T comparable] struct {
	ID  T
	Raw string
}

type TreeLinesConstructedMsg[T comparable] []TreeLine[T]

func constructTree[T comparable](tree ItemTree[T]) []TreeLine[T] {
	lines := make([]TreeLine[T], 0) // TODO: set appropriate cap
	lines = append(lines, TreeLine[T]{
		ID:  tree.ID(),
		Raw: tree.Content(),
	})
	childLines := constructChildren(tree.Children(), "")
	lines = append(lines, childLines...)

	return lines
}

func constructChildren[T comparable](children iter.Seq2[ItemTree[T], bool], prefix string) []TreeLine[T] {
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
