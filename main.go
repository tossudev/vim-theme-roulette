package main

import (
    "fmt"
	"os"

    tea "charm.land/bubbletea/v2"
)

type model struct {
	themes	[]string
}


func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}


func initialModel() model {
	return model{
		themes:  []string{"gruvbox", "elflord", "industry", "morning"},
	}
}

func (m model) Init() tea.Cmd {
    return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyPressMsg:
		switch msg.String() {
		
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	
	return m, nil
}


func (m model) View() tea.View {
	s := "Vim Theme Roulette >:D\n\n"

	for _, theme := range(m.themes) {
		s += fmt.Sprintf("  %s  ", theme)
	}

	s += "\n\nPress q to quit."

	return tea.NewView(s)
}



