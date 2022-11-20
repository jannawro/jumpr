package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	nickname, name, region, profile string
}

func (i item) Nickname() string {
	return i.nickname
}

func (i item) Name() string {
	return i.name
}

func (i item) Region() string {
	return i.region
}

func (i item) Profile() string {
	return i.profile
}

func (i item) FilterValue() string { return i.nickname }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func TeaPrompt() {
	fmt.Println(cfg[0].Nickname())

	//m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	//m.list.Title = "Clusters to choose from"
	//
	//p := tea.NewProgram(m, tea.WithAltScreen())
	//
	//if _, err := p.Run(); err != nil {
	//	fmt.Println("error running program: ", err)
	//	os.Exit(1)
	//}
}
