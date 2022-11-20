package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jannawro/jumpr/internal/clusterLogin"
	"log"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	choice   clusterLogin.Cluster
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	cursor int
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			choice = cfg[m.cursor]
			return m, tea.Quit
		case "down", "j":
			m.cursor++
			if m.cursor >= len(cfg) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(cfg) - 1
			}
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
	var items []list.Item
	for _, cluster := range cfg {
		items = append(items, item{
			title: cluster.Nickname,
			desc:  "| " + cluster.Name + " | " + cluster.Profile + " | " + cluster.Region + " |",
		})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Clusters to choose from"

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		log.Fatal("Failed to run prompt: ", err)
	}

	fmt.Println(choice.Nickname)
}
