// Alot of this is inspired by bubble gums fancy list example
// https://github.com/charmbracelet/bubbletea/tree/master/examples/list-fancy

package funcs

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
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
	var url string
	if i, ok := m.list.SelectedItem().(item); ok {
		url = i.Url()
	} else {
		fmt.Println("Failed to select item url")
		os.Exit(1)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "c" {
			m.list.NewStatusMessage(fmt.Sprintf("Repo URL: %v", url))
		}
		if msg.String() == "enter" {
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

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	help := []key.Binding{keys.choose, keys.copy}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type listKeyMap struct {
	choose key.Binding
	copy   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open"),
		),
		copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "view url"),
		),
	}
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
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	delegate := newItemDelegate(delegateKeys)
	llist := list.New(items, delegate, 0, 0)
	llist.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.choose,
			listKeys.copy,
		}
	}
	m := model{list: llist}
	m.list.Title = fmt.Sprintf("%v's Repos", name)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type delegateKeyMap struct {
	choose key.Binding
	copy   key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.copy,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.copy,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "view url"),
		),
	}
}
