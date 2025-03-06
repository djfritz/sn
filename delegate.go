package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		//i := m.SelectedItem().(*item)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.increase):
				H++
			case key.Matches(msg, keys.decrease):
				if H > 2 {
					H--
				}
			}
		}

		return nil
	}

	help := []key.Binding{keys.increase, keys.decrease}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	increase key.Binding
	decrease key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.increase,
		d.decrease,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.increase,
			d.decrease,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		increase: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "increase"),
		),
		decrease: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "decrease"),
		),
	}
}
