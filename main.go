package main

import (
    "fmt"
	"os"
	"time"

    tea "charm.land/bubbletea/v2"
)

var exit bool = false
var displaySize int = 25
var speed int = 5
var stopSpin = false


type model struct {
	fullText string
	displayText string
	index int
}


func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}


func initialModel() model {
	text := ""
	themes := []string{"gruvbox", "elflord", "industry", "morning"}

	for _, theme := range(themes) {
		text += fmt.Sprintf("| %s |", theme)
	}

	return model {
		fullText:	text,
		index:		0,
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
		return tea.NewView("Changed Vim theme to: [theme]\n")
	}

	s := "Vim Theme Roulette >:D\n\n"

	s += m.displayText
	s += "\n"
	for range(displaySize/2) {
		s += " "
	}
	s += "^"

	if speed <= 0 {
		s += "\nPress q to accept your fate!"
	} else {
		s += "\nPress space or enter to stop spinning."
	}

	return tea.NewView(s)
}


func (m *model) UpdateRoulette() {
	if speed <= 0 {
		return
	}

	m.index += speed
	m.displayText = ""
	wrap := 0

	if stopSpin {
		speed -= 1
	}

	if m.index >= len(m.fullText) - 1 {
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

type tickMsg time.Time                                                                                                                                                                        
                                                                                                                                                                                              
func tick() tea.Msg {                                                                                                                                                                         
    time.Sleep(time.Duration(50) * time.Millisecond)
    return tickMsg(time.Now())
} 

