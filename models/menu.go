package models

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)


var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).AlignHorizontal(lipgloss.Center).Padding(0, 0, 2, 1)
var activeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))

var art = `┏┓┏┳┓┏┓┏┓┳┳┓  ┳┳┓┏┓┳┓┳┏┓┏┓┏┓┏┳┓
┗┓ ┃ ┣ ┣┫┃┃┃  ┃┃┃┣┫┃┃┃┣ ┣ ┗┓ ┃ 
┗┛ ┻ ┗┛┛┗┛ ┗  ┛ ┗┛┗┛┗┻┻ ┗┛┗┛ ┻ `


type menu struct {
	choices []string // items on the to-do list
	cursor  int      // which to-do list item our cursor is pointing at
}

func InitialMenu() menu {
	return menu{
		choices: []string{"Download Game", "Settings", "Quit"},
	}
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the space bar toggle the selected state
		// for the item that the cursor is pointing at.
		case "enter", "space":
			switch m.cursor {
			case 0:
				return initialDownload(), nil
			case 1:
				return initialSettings(), nil
			case 2:
				return m, tea.Quit
			}
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m menu) View() tea.View {
	// The header
	s := titleStyle.Render(art) + "\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " "  // no cursor
		checked := " " // not selected
		if m.cursor == i {
			cursor = activeStyle.Render(">")
			checked = activeStyle.Render("x")
			choice = activeStyle.Render(choice)
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		
	}
	s += Footer() + "\n"
	// Send the UI for rendering
	return tea.NewView(s)
}

func (m menu) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func Footer() string {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).PaddingTop(2).
        Render(" ______________________________________\n| ESC • Back     |     ENTER • Confirm |")
}