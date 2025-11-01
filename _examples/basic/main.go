package main

import (
	"iter"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ras0q/bubbletree"
)

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

type rootModel struct {
	tree bubbletree.Model[int]
}

func New() rootModel {
	tree := bubbletree.New[int]()
	tree.UpdateFunc = func(line bubbletree.TreeLine[int], msg tea.Msg) tea.Cmd {
		var cmd tea.Cmd

		switch msg := msg.(type) {
		case tea.KeyMsg:
			focusedItem := mockTree.search(line.ID)

			switch msg.String() {
			case "h":
				focusedItem.isOpened = false
				cmd = tree.SetTree(mockTree)
			case "l":
				focusedItem.isOpened = true
				cmd = tree.SetTree(mockTree)
			}
		}

		return cmd
	}

	return rootModel{
		tree: tree,
	}
}

var _ tea.Model = rootModel{}
var _ tea.ViewModel = rootModel{}

var mockTree = itemTree{
	id:      1,
	content: "Alice",
	children: []*itemTree{
		{
			id:      2,
			content: "Bob",
			children: []*itemTree{
				{
					id:       3,
					content:  "Charlie",
					children: []*itemTree{},
				},
			},
		},
		{
			id:       4,
			content:  "Diana",
			children: []*itemTree{},
		},
	},
}

// Init implements tea.Model.
func (m rootModel) Init() tea.Cmd {
	return tea.Batch(
		m.tree.Init(),
		m.tree.SetTree(mockTree),
	)
}

// Update implements tea.Model.
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)

	return m, cmd
}

// View implements tea.ViewModel.
func (m rootModel) View() string {
	return m.tree.View()
}

type itemTree struct {
	id       int
	content  string
	children []*itemTree
	isOpened bool
}

var _ bubbletree.Tree[int] = itemTree{}

func (t itemTree) ID() int {
	return t.id
}

func (t itemTree) Content() string {
	prefix := ""
	if len(t.children) > 0 {
		if t.isOpened {
			prefix = "▼"
		} else {
			prefix = "▶"
		}
	}

	return prefix + t.content
}

func (t itemTree) Children() iter.Seq2[bubbletree.Tree[int], bool] {
	return func(yield func(bubbletree.Tree[int], bool) bool) {
		if !t.isOpened {
			return
		}

		for i, child := range t.children {
			hasNext := i < len(t.children)-1
			yield(child, hasNext)
		}
	}
}

func (t *itemTree) search(id int) *itemTree {
	if t.id == id {
		return t
	}

	for _, child := range t.children {
		result := child.search(id)
		if result != nil {
			return result
		}
	}

	return nil
}
