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
	tree := itemTree{
		id:      1,
		content: "Alice",
		children: []itemTree{
			{
				id:      2,
				content: "Bob",
				children: []itemTree{
					{
						id:       3,
						content:  "Charlie",
						children: []itemTree{},
					},
				},
			},
			{
				id:       4,
				content:  "Diana",
				children: []itemTree{},
			},
		},
	}

	return rootModel{
		tree: bubbletree.New(tree),
	}
}

var _ tea.Model = rootModel{}
var _ tea.ViewModel = rootModel{}

// Init implements tea.Model.
func (m rootModel) Init() tea.Cmd {
	return m.tree.Init()
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
	children []itemTree
}

var _ bubbletree.ItemTree[int] = itemTree{}

func (t itemTree) ID() int {
	return t.id
}

func (t itemTree) Content() string {
	return t.content
}

func (t itemTree) Children() iter.Seq2[bubbletree.ItemTree[int], bool] {
	return func(yield func(bubbletree.ItemTree[int], bool) bool) {
		for i, child := range t.children {
			hasNext := i < len(t.children)-1
			yield(child, hasNext)
		}
	}
}
