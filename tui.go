package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os/exec"
	"runtime"
)

// Styles
var docStyle = lipgloss.NewStyle().Margin(1, 2)
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// item is a list item for the Hacker News list
type item struct {
	title, url, storyText string
}

// Title returns the title of the item
func (i item) Title() string { return i.title }

// Description returns the url of the item
func (i item) Description() string { return i.url }

// FilterValue returns the title of the item
func (i item) FilterValue() string { return i.title }

// listKeyMap maps keys to functions
type listKeyMap struct {
	openUrl key.Binding
}

// newListKeyMap returns a new listKeyMap
func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		openUrl: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Open link"),
		),
	}
}

// model is the Bubble Tea model
type model struct {
	list        list.Model
	read        viewport.Model
	keys        *listKeyMap
	readingMode bool
}

// Init initializes the Bubble Tea model
func (m model) Init() tea.Cmd {
	return nil
}

// Update updates the Bubble Tea model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle key presses
	case tea.KeyMsg:
		switch msg.String() {

		// Quit on ctrl+c
		case "ctrl+c":
			return m, tea.Quit

		// Open news item
		case "enter":

			// If in reading mode, go back to list
			if m.readingMode {
				m.readingMode = false
				return m, nil
			}

			urlToOpen := m.list.Items()[m.list.Index()].(item).url

			// If no url, show story text
			if urlToOpen == "" {
				content := m.list.Items()[m.list.Index()].(item).storyText
				renderer, err := glamour.NewTermRenderer(
					glamour.WithAutoStyle(),
					glamour.WithWordWrap(m.list.Width()),
				)
				if err != nil {
					fmt.Printf("Error creating render: %s\n", err)
				}

				str, err := renderer.Render(content)
				if err != nil {
					fmt.Printf("Error rendering content: %s\n", err)
				}

				m.read.SetContent(str)
				m.readingMode = true

				return m, nil
			} else {
				// Open url
				openBrowser(m.list.Items()[m.list.Index()].(item).url)
			}

			return m, nil
		}

	// Handle window resizes
	case tea.WindowSizeMsg:
		// Update list dimensions
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		// Update read dimensions
		m.read.Width = msg.Width - h
		m.read.Height = msg.Height - v - 2

		return m, nil
	}

	var cmd tea.Cmd
	if !m.readingMode {
		// Update list
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	} else {
		// Update read
		m.read, cmd = m.read.Update(msg)
		return m, cmd
	}
}

// View returns the Bubble Tea view
func (m model) View() string {
	// If in reading mode, show read view
	if m.readingMode {
		return m.read.View() + m.readHelpView()
	} else {
		return docStyle.Render(m.list.View())
	}

}

// readHelpView returns the help view for the read view
func (m model) readHelpView() string {
	return helpStyle("\n  ↑/↓: Navigate • enter: Go back\n")
}

// initialModel returns the initial Bubble Tea model
func initialModel() model {

	// Get the data from the API
	a := getApi()
	urls := getHnUrl()
	urls.topStories()
	data, err := a.request(urls.string())

	if err != nil {
		fmt.Printf("Error getting data: %s\n", err)
		return model{}
	}

	// Create list items with data
	var items []list.Item
	for _, hit := range data.Hits {
		formattedTitle := fmt.Sprintf("▴%d | %s", hit.Points, hit.Title)
		items = append(items, item{title: formattedTitle, url: hit.URL, storyText: hit.StoryText})
	}

	// Create list view
	listKeys := newListKeyMap()
	newsList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	newsList.Title = "Hacker News"
	keysExtra := func() []key.Binding {
		return []key.Binding{
			listKeys.openUrl,
		}
	}
	newsList.AdditionalShortHelpKeys = keysExtra
	newsList.AdditionalFullHelpKeys = keysExtra

	// Create read view
	vp := viewport.New(20, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return model{
		list: newsList,
		keys: listKeys,
		read: vp,
	}
}

// openBrowser opens the specified URL in the default browser of the user.
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
