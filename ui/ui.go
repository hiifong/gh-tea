package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hiifong/gh-tea/ui/component/tab"
	tkey "github.com/hiifong/gh-tea/ui/key"
)

type Model struct {
	keys *tkey.Map
	Tab  tea.Model
}

func New() Model {
	return Model{
		keys: tkey.Default,
		Tab:  tab.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd     tea.Cmd
		tabsCmd tea.Cmd

		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tkey.Default.Quit):
			return m, tea.Quit
		}
		m.Tab, tabsCmd = m.Tab.Update(msg)
		return m, tabsCmd
	}

	cmds = append(
		cmds,
		cmd,
		tabsCmd,
	)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.Tab.View()
}

func Run() error {
	p := tea.NewProgram(New())
	_, err := p.Run()
	return err
}
