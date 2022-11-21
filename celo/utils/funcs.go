package funcs

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Repo struct {
	Description string `json:"description"`
	Full_name   string `json:"full_nane"`
	Html_url    string `json:"html_url"`
	Name        string `json:"name"`
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type item struct {
	title, desc, url string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) Url() string         { return i.url }
func (i item) FilterValue() string { return i.title }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			var url string
			if i, ok := m.list.SelectedItem().(item); ok {
				url = i.Url()
			} else {
				fmt.Println("Failed to select item url")
				os.Exit(1)
			}
			browser.OpenURL(url)
		}
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

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func RepoList(repos []Repo, name string) {
	items := []list.Item{}
	for x := 0; x < len(repos); x++ {
		r := repos[x]
		items = append(items, item{
			title: r.Name,
			desc:  r.Description,
			url:   r.Html_url,
		})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = fmt.Sprintf("%v's Repos", name)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
