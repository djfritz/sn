// Copyright 2025 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package main

import (
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/djfritz/sn/stack"
)

var (
	S        *stack.Stack
	H        = 10
	delegate list.DefaultDelegate
	items    []list.Item
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
}

func (i *item) Title() string { return i.title }
func (i *item) Description() string {
	return i.description
}
func (i *item) FilterValue() string { return i.title + " " + i.description }

type model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

func newModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	for _, v := range S.R {
		items = append(items, &item{
			title:       v.Header(),
			description: v.CSString(),
		})
	}

	// Setup list
	delegate = newItemDelegate(delegateKeys)

	stackList := list.New(items, delegate, 0, 0)
	stackList.Title = os.Args[1]
	stackList.Styles.Title = titleStyle
	stackList.SetFilteringEnabled(true)
	stackList.Filter = filter

	return model{
		list:         stackList,
		delegateKeys: delegateKeys,
	}
}

func filter(f string, d []string) []list.Rank {
	var ret []list.Rank
	for i, v := range d {
		if x := strings.Index(v, f); x != -1 {
			mi := []int{x}
			for i := 0; i < len(f); i++ {
				mi = append(mi, x+i)
			}
			ret = append(ret, list.Rank{
				Index:          i,
				MatchedIndexes: mi,
			})
		}
	}
	return ret
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.MaxHeight(H)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.MaxHeight(H)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.MaxHeight(H)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.MaxHeight(H)
	delegate.Styles.DimmedTitle = delegate.Styles.DimmedTitle.MaxHeight(H)
	delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.MaxHeight(H)
	delegate.SetHeight(H)

	m.list.SetDelegate(delegate)

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invalid arguments: expected filename")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("opening file:", err)
	}

	S, err = stack.Decode(f)
	if err != nil {
		log.Fatal("decoding stacktrace:", err)
	}

	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
