package main

import (
    "fmt"
	"os"
	"time"

    tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var exit bool = false
var displaySize int = 50
var speed int = 50
var stopSpin = false

const (
	ColorReset string = "\033[0m"
	ColorGreen string = "\u001B[32m"
	ColorYellow string = "\u001B[33m"
	ColorWhite string = "\u001B[37m"
)

type model struct {
	fullText string
	displayText string
	index int
	themes []Theme
	themeIndices []int

	width int
	height int
}


func main() {
	FetchConfig()
	ShuffleThemes()

    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("ERR in main(): %v", err)
        os.Exit(1)
    }
}


func initialModel() model {
	text := ""
	themeIndices := []int{}
	totalLength := 0

	for _, theme := range(Config.Themes) {
		themeLength := len(theme.Name) + 4
		themeIndices = append(themeIndices, themeLength + totalLength)
		totalLength += themeLength
		text += fmt.Sprintf("| %s |", theme.Name)
	}

	return model {
		fullText:		text,
		index:			0,
		themes:			Config.Themes,
		themeIndices:	themeIndices,
	}
}

func (m model) Init() tea.Cmd {
    return tick
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyPressMsg:
		switch msg.String() {
		
		case "ctrl+c", "q":
			exit = true
			return m, tea.Quit

		case "space", "enter":
			stopSpin = true
		}

	case tickMsg:
		m.UpdateRoulette()
		return m, tick

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}


func (m model) View() tea.View {
	width := m.width

	if exit {
		if CurrentTheme.Name != "" {
			ChangeTheme()
			return tea.NewView(fmt.Sprintf("Changed Vim theme to: %s\n", CurrentTheme.Name))
		}

		return tea.NewView(fmt.Sprintf("Program exited by user"))
	}

	s := ColorGreen + "\n★ ☆ ★  Welcome to the Vim Theme Roulette!  ☆ ★ ☆ ★\n\n" + ColorReset

	s += "░▒▓ "
	s += m.displayText
	s += " ▓▒░"
	

	s += "\n"
	for range(m.width) {
		s += " "
	}
	s += "^\n"

	if speed <= 0 {
		s += fmt.Sprintf("You rolled %s!", CurrentTheme.Name)
	} else {
		s += "Press space or enter to stop spinning."
	}

	s += "\nPress q to exit."
	s += "\n"

	centered := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("5")).
		Render(s)

	return tea.NewView(centered)
}


func (m *model) UpdateRoulette() {
	if speed <= 0 {
		if CurrentTheme.Name == "" {
			m.GetTheme()
		}

		return
	}

	m.index += (speed / 10)
	m.displayText = ""
	wrap := 0

	if stopSpin {
		speed -= 1
	}

	if m.index > len(m.fullText) {
		m.index = 0
	}

	for i := range(displaySize) {
		index := m.index + i

		if wrap > 0 {
			index = i - wrap
		}

		if index >= len(m.fullText) - 1 {
			wrap = i
		}

		m.displayText += string(rune(m.fullText[index]))
	}
}

func (m *model) GetTheme() {
	for i, v := range(m.themeIndices) {
		if v > m.index + displaySize/2 {
			CurrentTheme = m.themes[i]
			return
		}
	}
	CurrentTheme = m.themes[0]
}


type tickMsg time.Time

func tick() tea.Msg {                                                                                                                                                                         
    time.Sleep(time.Duration(16) * time.Millisecond)
    return tickMsg(time.Now())
} 

