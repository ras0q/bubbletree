package bubbletree

import (
	"fmt"
	"iter"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Model[T comparable] struct {
	ItemTree ItemTree[T]
}

func New[T comparable](tree ItemTree[T]) Model[T] {
	return Model[T]{
		ItemTree: tree,
	}
}

func (m Model[T]) Init() tea.Cmd {
	return nil
}

func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
	return m, nil
}

func (m Model[T]) View() string {
	return renderTree(m.ItemTree)
}

type ItemTree[T comparable] interface {
	ID() T
	Children() iter.Seq2[ItemTree[T], bool]
}

func renderTree[T comparable](tree ItemTree[T]) string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "%+v\n", tree.ID())
	renderChildren(&b, tree.Children(), "")

	return b.String()
}

func renderChildren[T comparable](b *strings.Builder, children iter.Seq2[ItemTree[T], bool], prefix string) {
	for child, hasNext := range children {
		connector := "├── "
		nextPrefix := "│   "
		if !hasNext {
			connector = "└── "
			nextPrefix = "    "
		}

		b.WriteString(prefix)
		b.WriteString(connector)
		fmt.Fprintf(b, "%+v\n", child.ID())

		renderChildren(b, child.Children(), nextPrefix)
	}
}
