package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// State represents the current state of the CLI application
type State int

const (
	// StateHome represents the home screen
	StateHome State = iota
	// StateUsername represents the screen asking for username
	StateUsername
	// StateCollection represents the screen asking for collection name
	StateCollection
	// StateOperations represents the screen for CRUD operations
	StateOperations
	// StateKeyValue represents the screen for entering key and value
	StateKeyValue
	// StateMessage represents the screen displaying messages
	StateValue

	StateKeyOperations

	StateMessage
)

// Model represents the data model for the CLI application
type Model struct {
	State           State
	list            list.Model
	keyList         list.Model
	Username        string
	CollectionName  string
	Message         string
	keyInput        textinput.Model
	valueInput      textinput.Model
	userNameInput   textinput.Model
	collectionInput textinput.Model
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var messageStyle lipgloss.Style

type item struct {
	title, desc string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Username"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	ti2 := textinput.New()
	ti2.Placeholder = "Enter Collection"
	ti2.CharLimit = 156
	ti2.Width = 20

	ti3 := textinput.New()
	ti3.Placeholder = "Enter Key"
	ti3.CharLimit = 156
	ti3.Width = 20

	ti4 := textinput.New()
	ti4.Placeholder = "Enter Value"
	ti4.CharLimit = 156
	ti4.Width = 20

	items := []list.Item{

		item{title: "Read", desc: "To Use the collection"},
		item{title: "Delete", desc: "To Delete the collection"},
		item{title: "Create", desc: "To Create a new collection"},
	}

	items2 := []list.Item{
		item{title: "Read", desc: "To Read the key"},
		item{title: "Delete", desc: "To Delete the key"},
		item{title: "Create", desc: "To Create a new Key-Value pair"},
		item{title: "Update", desc: "To Update the key"},
	}
	messageStyle = lipgloss.NewStyle().Bold(true)
	m := Model{
		list:            list.New(items, list.NewDefaultDelegate(), 0, 0),
		userNameInput:   ti,
		collectionInput: ti2,
		keyInput:        ti3,
		valueInput:      ti4,
		keyList:         list.New(items2, list.NewDefaultDelegate(), 0, 0),
		State:           StateHome,
		Message:         ""}
	m.list.Title = "Welcome to Arachne-DB, enter collection to continue.."
	m.keyList.Title = "Perform one of the four operations"
	return m
}

// Init implements tea.Model.
func (Model) Init() tea.Cmd {
	return nil
}

// Update function updates the model based on user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd = []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			switch m.State {
			case StateHome:
				m.State = StateCollection
				m.userNameInput.Blur()
				m.collectionInput.Focus()
				// Prepare the JSON payload
				payload := fmt.Sprintf(`{"username": "%s"}`, m.userNameInput.Value())
				// Send a POST request
				resp, err := http.Post("http://localhost:8080/username", "application/json", strings.NewReader(payload))
				if err != nil {
					m.State = StateHome
					m.collectionInput.Blur()
					m.userNameInput.Focus()
					m.Message = "Arachne-DB is down ˙◠˙"
					messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
				} else {
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						m.Message = fmt.Sprintf("Error reading response: %v", err)
						messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
					} else {
						m.Message = string(body)
						messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
					}
				}

			case StateOperations:
				m.State = StateKeyValue
				m.collectionInput.Blur()
				m.keyInput.Focus()
				switch State(m.list.Cursor()) {
				case 0:
					url := fmt.Sprintf("http://localhost:8080/user/get/collection/%s/%s",
					m.userNameInput.Value(), m.collectionInput.Value())

					// Send a GET request
					resp, err := http.Get(url)
					if err != nil {
						m.State = StateHome
						m.collectionInput.Blur()
						m.userNameInput.Focus()
						m.Message = "Arachne-DB is down ˙◠˙"
						messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
					} else {
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							m.Message = fmt.Sprintf("Error reading response: %v", err)
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
						} else {
							m.Message = string(body)
							if m.Message=="\"Collection not found!\""{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
								m.State = StateCollection
								m.keyInput.Blur()
								m.collectionInput.Focus()
							}else{
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
							}
						}

					}
				case 1:
					payload := fmt.Sprintf(`{"username": "%s", "collection_name": "%s"}`, m.userNameInput.Value(), m.collectionInput.Value())

					// Send a POST request
					resp, err := http.Post("http://localhost:8080/user/delete/collection","application/json",strings.NewReader(payload))
					if err != nil {
						m.State = StateHome
						m.collectionInput.Blur()
						m.userNameInput.Focus()
						m.keyInput.Blur()

						m.Message = "Arachne-DB is down ˙◠˙"
						messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
					} else {
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							m.Message = fmt.Sprintf("Error reading response: %v", err)
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
						} else {
							m.Message = string(body)
							if m.Message=="\"Collection not found!\""{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
								m.State = StateCollection
								m.keyInput.Blur()
								m.collectionInput.Focus()
							}else{
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
							m.State = StateCollection
							m.keyInput.Blur()
							m.collectionInput.Focus()
							}
						}

					}
				case 2:
					payload := fmt.Sprintf(`{"username": "%s", "collection_name": "%s"}`, m.userNameInput.Value(), m.collectionInput.Value())

					resp, err := http.Post("http://localhost:8080/user/create/collection","application/json",strings.NewReader(payload))
					if err != nil {
						m.State = StateHome
						m.collectionInput.Blur()
						m.userNameInput.Focus()
						m.keyInput.Blur()

						m.Message = "Arachne-DB is down ˙◠˙"
						messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
					} else {
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							m.Message = fmt.Sprintf("Error reading response: %v", err)
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
							m.State = StateCollection
							m.keyInput.Blur()
							m.collectionInput.Focus()
						} else {
							m.Message = string(body)
							if m.Message=="\"Collection not found!\""{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
								m.State = StateCollection
								m.keyInput.Blur()
								m.collectionInput.Focus()
							}else{
								
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
							}
						}

					}
				}
				
			case StateKeyOperations:
				switch State(m.keyList.Cursor()) {
					case 0:
						url := fmt.Sprintf("http://localhost:8080/user/get/key/%s/%s/%s",
						m.userNameInput.Value(), m.collectionInput.Value(),m.keyInput.Value())

						// Send a GET request
						resp, err := http.Get(url)

						if err != nil {
							m.State = StateHome
							m.collectionInput.Blur()
							m.userNameInput.Focus()
							m.Message = "Arachne-DB is down ˙◠˙"
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
						} else {
							defer resp.Body.Close()
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								m.Message = fmt.Sprintf("Error reading response: %v", err)
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
							} else {
								m.Message = string(body)
								if m.Message=="{\"error\":\"Key not found\"}"{
									messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
									m.State = StateKeyValue
									m.keyInput.Focus()
									m.collectionInput.Blur()
									m.valueInput.Blur()
								}else{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
								m.State=StateKeyValue
								m.keyInput.Focus()
							}
							}

						}
					case 1:
						payload := fmt.Sprintf(`{"username": "%s", "collection": "%s", "key": "%s"}`, m.userNameInput.Value(), m.collectionInput.Value(),m.keyInput.Value())

						resp, err := http.Post("http://localhost:8080/user/delete/key","application/json",strings.NewReader(payload))

						if err != nil {
							m.State = StateHome
							m.collectionInput.Blur()
							m.userNameInput.Focus()
							m.Message = "Arachne-DB is down ˙◠˙"
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
						} else {
							defer resp.Body.Close()
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								m.Message = fmt.Sprintf("Error reading response: %v", err)
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
							} else {
								m.Message = string(body)
								if m.Message=="{\"error\":\"Key not found\"}"{
									messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
									m.State = StateKeyValue
									m.keyInput.Focus()
									m.collectionInput.Blur()
									m.valueInput.Blur()
								}else{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
								m.keyInput.Focus()
								m.State=StateKeyValue
								}
							}

						}
					case 2:
						payload := fmt.Sprintf(`{"username": "%s", "collection": "%s", "key": "%s", "value": "%s"}`, m.userNameInput.Value(), m.collectionInput.Value(),m.keyInput.Value(), m.valueInput.Value())

						// Send a GET request
						resp, err := http.Post("http://localhost:8080/user/create/key","application/json",strings.NewReader(payload))


						if err != nil {
							m.State = StateHome
							m.collectionInput.Blur()
							m.userNameInput.Focus()
							m.Message = "Arachne-DB is down ˙◠˙"
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
						} else {
							defer resp.Body.Close()
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								m.Message = fmt.Sprintf("Error reading response: %v", err)
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
							} else {
								m.Message = string(body)
								if m.Message=="{\"error\":\"Key not found\"}"{
									messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
									m.State = StateKeyValue
									m.keyInput.Focus()
									m.collectionInput.Blur()
									m.valueInput.Blur()
								}else{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
								m.keyInput.Focus()
								m.State=StateKeyValue
								}
							}

						}
					case 3:
						payload := fmt.Sprintf(`{"username": "%s", "collection": "%s", "key": "%s", "value": "%s"}`, m.userNameInput.Value(), m.collectionInput.Value(),m.keyInput.Value(), m.valueInput.Value())

						// Send a POST request
						resp, err := http.Post("http://localhost:8080/user/create/key","application/json",strings.NewReader(payload))

						if err != nil {
							m.State = StateHome
							m.collectionInput.Blur()
							m.userNameInput.Focus()
							m.Message = "Arachne-DB is down ˙◠˙"
							messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
						} else {
							defer resp.Body.Close()
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								m.Message = fmt.Sprintf("Error reading response: %v", err)
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
							} else {
								m.Message = string(body)
								if m.Message=="{\"error\":\"Key not found\"}"{
									messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
									m.State = StateKeyValue
									m.keyInput.Focus()
									m.collectionInput.Blur()
									m.valueInput.Blur()
								}else{
								messageStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00AA00")) // Green color for success messages
								m.keyInput.Focus()
								m.State=StateKeyValue
								}
							}
						}
					}
			case StateCollection:
				m.State = StateOperations

			case StateKeyValue:
				m.State = StateValue
				m.keyInput.Blur()
				m.valueInput.Focus()
			case StateValue:
				m.State=StateKeyOperations
				m.valueInput.Blur()
				return m,tea.Batch(cmd...)
			}
			return m, tea.Batch(cmd...)
		}
		if msg.String() == "tab" {
			switch m.State {
			case StateCollection:
				m.State = StateOperations
				return m, tea.Batch(cmd...)
			}

		}
		if msg.String() == "shift+tab" {
			switch m.State {
			case StateOperations:
				m.State = StateCollection
				m.collectionInput.Focus()
			case StateValue:
				m.State = StateKeyValue
				m.valueInput.Blur()
				m.keyInput.Focus()
			}
			return m, tea.Batch(cmd...)
		}
		if msg.String() == "esc" {
			switch m.State {

			case StateOperations:
				m.State = StateHome
				m.collectionInput.Blur()
				m.userNameInput.Focus()
			case StateCollection:
				m.State = StateHome
				m.collectionInput.Blur()
				m.userNameInput.Focus()

			case StateKeyValue:
				m.State = StateCollection
				m.keyInput.Blur()
				m.valueInput.Blur()
				m.collectionInput.Focus()

			case StateValue:
				m.State = StateCollection
				m.keyInput.Blur()
				m.valueInput.Blur()
				m.collectionInput.Focus()
			case StateKeyOperations:
				m.State = StateKeyValue
				m.valueInput.Blur()
				m.keyInput.Focus()
			}
			m.Message = ""
			return m, tea.Batch(cmd...)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.keyList.SetSize(msg.Width-h, msg.Height-v)

	}
	switch m.State {
	case StateHome:
		var cmd1 tea.Cmd
		m.userNameInput, cmd1 = m.userNameInput.Update(msg)
		cmd = append(cmd, cmd1)
	case StateCollection:
		var cmd1 tea.Cmd

		m.collectionInput, cmd1 = m.collectionInput.Update(msg)

		cmd = append(cmd, cmd1)
	case StateOperations:
		var cmd2 tea.Cmd
		m.list, cmd2 = m.list.Update(msg)
		cmd = append(cmd, cmd2)
	case StateKeyValue:
		var cmd1 tea.Cmd
		m.keyInput, cmd1 = m.keyInput.Update(msg)
		cmd = append(cmd, cmd1)
	case StateValue:
		var cmd1 tea.Cmd
		m.valueInput, cmd1 = m.valueInput.Update(msg)
		cmd = append(cmd, cmd1)
	case StateKeyOperations:
		var cmd2 tea.Cmd
		m.keyList, cmd2 = m.keyList.Update(msg)
		cmd = append(cmd, cmd2)
	}

	return m, tea.Batch(cmd...)
}

// View function renders the view based on the current state
func (m Model) View() string {
	var view string

	switch m.State {
	case StateHome:
		view = m.homeView()
	case StateValue:
		view = m.keyValueView()
	case StateCollection:
		view = m.collectionView()
	case StateOperations:
		view = m.collectionView()
	case StateKeyValue:
		view = m.keyValueView()
	case StateMessage:
		view = messageView(m.Message)
	case StateKeyOperations:
		view = m.keyValueView()
	}


	return view
}

func (m Model) homeView() string {
	banner := "\n\n" + "                                        █████                                 █████ █████    " + "     /      \\      \n" +
		"                                       ░░███                                 ░░███ ░░███     " + "  \\  \\  ,,  /  /   \n" +
		"  ██████   ████████   ██████    ██████  ░███████   ████████    ██████      ███████  ░███████ " + "    '-.\\()/'.-'     \n" +
		" ░░░░░███ ░░███░░███ ░░░░░███  ███░░███ ░███░░███ ░░███░░███  ███░░███    ███░░███  ░███░░███" + "  .--_'(  )'_--.   \n" +
		"  ███████  ░███ ░░░   ███████ ░███ ░░░  ░███ ░███  ░███ ░███ ░███████    ░███ ░███  ░███ ░███" + " / /` /`\"\"`\\ `\\ \\  \n" +
		" ███░░███  ░███      ███░░███ ░███  ███ ░███ ░███  ░███ ░███ ░███░░░     ░███ ░███  ░███ ░███" + "  |  |  ><  |  |   \n" +
		"░░████████ █████    ░░████████░░██████  ████ █████ ████ █████░░██████    ░░████████ ████████ " + "  \\  \\      /  / \n" +
		" ░░░░░░░░ ░░░░░      ░░░░░░░░  ░░░░░░  ░░░░ ░░░░░ ░░░░ ░░░░░  ░░░░░░      ░░░░░░░░ ░░░░░░░░  \n"

	bannerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8"))
	return lipgloss.JoinVertical(lipgloss.Bottom, lipgloss.JoinVertical(lipgloss.Left, bannerStyle.Render(banner), (m.userNameInput.View())), messageStyle.Render(m.Message))

}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m Model) collectionView() string {
	// Render the screen asking for collection name
	return lipgloss.JoinVertical(lipgloss.Bottom,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			docStyle.Render(m.collectionInput.View()),
			docStyle.Render(m.list.View()), // Render the list component
		),
		messageStyle.Render(m.Message),
	)
}

// func operationsView() string {
// 	// Render the screen for CRUD operations
// 	return ""
// }

func (m Model) keyValueView() string {
	// Render the screen for entering key and value
	return lipgloss.JoinVertical(lipgloss.Bottom, lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Center, docStyle.Render(m.keyInput.View()), docStyle.Render(m.valueInput.View())),
		docStyle.Render(m.keyList.View()), // Render the list component
	),
		messageStyle.Render(m.Message))
}

func messageView(msg string) string {
	// Render the message and return values of API calls
	return ""
}
