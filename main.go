package main

import (
    "fmt"
	"os"
	"time"

    tea "charm.land/bubbletea/v2"
)

var exit bool = false
var displaySize int = 50
var speed int = 5
var stopSpin = false
var ColorReset string = "\033[0m"

type model struct {
	fullText string
	displayText string
	index int
	themes []Theme
	themeIndices []int
}


func main() {
	GetThemesLocal(true)
	GetThemesLocal(false)
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

	for _, theme := range(Themes) {
		themeLength := len(theme.name) + 4
		themeIndices = append(themeIndices, themeLength + totalLength)
		totalLength += themeLength
		text += fmt.Sprintf("| %s |", theme.name)
	}

	return model {
		fullText:		text,
		index:			0,
		themes:			Themes,
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
	}

	return m, nil
}


func (m model) View() tea.View {
	if exit {
		ChangeTheme()
		return tea.NewView(fmt.Sprintf("Changed Vim theme to: %s\n", CurrentTheme.name))
	}

	s := "Vim Theme Roulette >:D\n\n"

	s += m.displayText
	s += "\n"
	for range(displaySize/2) {
		s += " "
	}
	s += "^\n"

	if speed <= 0 {
		s += fmt.Sprintf("You rolled %s!", CurrentTheme.name)
		s += "\nPress q to accept your fate!"
	} else {
		s += "\nPress space or enter to stop spinning."
	}

	return tea.NewView(s)
}


func (m *model) UpdateRoulette() {
	if speed <= 0 {
		if CurrentTheme.name == "" {
			m.GetTheme()
		}

		return
	}

	m.index += speed
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
    time.Sleep(time.Duration(50) * time.Millisecond)
    return tickMsg(time.Now())
} 

