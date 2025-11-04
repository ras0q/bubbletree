package main

import (
	"iter"

	tea "github.com/charmbracelet/bubbletea"
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
	tree.OnUpdate = func(_ []bubbletree.RenderedLine[int], focusedID int, msg tea.Msg) tea.Cmd {
		var cmd tea.Cmd

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "h":
				searchResult, ok := mockTree.search(focusedID)
				if !ok {
					break
				}

				item := searchResult.item
				parent := searchResult.parent
				newFocusedID := focusedID

				if item.isOpened {
					item.isOpened = false
				} else {
					parent.isOpened = false
					newFocusedID = parent.ID()
				}

				cmd = tea.Batch(
					tree.SetTree(mockTree),
					tree.SetFocusedID(newFocusedID),
				)

			case "l":
				searchResult, ok := mockTree.search(focusedID)
				if !ok {
					break
				}

				item := searchResult.item
				if item.isLeaf() {
					break
				}

				item.isOpened = true

				cmd = tea.Batch(
					tree.SetTree(mockTree),
					tree.SetFocusedID(item.ID()),
				)
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
					id:      3,
					content: "Charlie",
					children: []*itemTree{
						{
							id:       4,
							content:  "Diana",
							children: []*itemTree{},
						},
					},
				},
			},
		},
		{
			id:       5,
			content:  "Eve",
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

var _ bubbletree.Node[int] = itemTree{}

func (t itemTree) ID() int {
	return t.id
}

func (t itemTree) Content() string {
	prefix := ""
	if !t.isLeaf() {
		if t.isOpened {
			prefix = "▼"
		} else {
			prefix = "▶"
		}
	}

	return prefix + t.content
}

func (t itemTree) Children() iter.Seq2[bubbletree.Node[int], bool] {
	return func(yield func(bubbletree.Node[int], bool) bool) {
		if !t.isOpened {
			return
		}

		for i, child := range t.children {
			hasNext := i < len(t.children)-1
			yield(child, hasNext)
		}
	}
}

func (t itemTree) isLeaf() bool {
	return len(t.children) == 0
}

type treeSearchResult struct {
	item   *itemTree
	parent *itemTree
}

func (t *itemTree) search(id int) (*treeSearchResult, bool) {
	if t.id == id {
		return &treeSearchResult{
			item:   t,
			parent: nil,
		}, true
	}

	for _, child := range t.children {
		if child.id == id {
			return &treeSearchResult{
				item:   child,
				parent: t,
			}, true
		}

		result, ok := child.search(id)
		if ok {
			return result, true
		}
	}

	return nil, false
}
