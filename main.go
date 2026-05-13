package main

import (
    "fmt"
	"os"
	"time"

    tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	ColorReset string = "\033[0m"
	ColorGreen string = "\u001B[32m"
	ColorYellow string = "\u001B[33m"
	ColorWhite string = "\u001B[37m"
)

type model struct {
	Header 			string
	BlockLeft 		string
	BlockRight 		string

	FullText 		string
	DisplayText 	string
	Index 			int
	Themes 			[]Theme
	ThemeIndices	[]int
	ThemeChoices 	[]string
	ThemesSelected 	map[int]bool
	Cursor 			int

	Width 			int
	Height 			int
}

var (
	start bool = true
	exit bool = false
	displaySize int = 50
	speed int = 50
	stopSpin = false
)


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
	m := model {
		Header:			ColorGreen + "\n☆ ★ ☆ ★  Vim Theme Roulette  ☆ ★ ☆ ★\n\n" + ColorReset,
		BlockLeft:		"░▒▓ ",
		BlockRight:		" ▓▒░",

		ThemeChoices:	[]string{"All", "Favorites", "Start!"},
		ThemesSelected:	make(map[int]bool),
	}

	
	return m
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
			if start {
				if m.Cursor == 2 {
					m.AddThemes()
					start = false
				} else {
					if _, ok := m.ThemesSelected[m.Cursor]; ok {
						delete(m.ThemesSelected, m.Cursor);
					} else {
						m.ThemesSelected[m.Cursor] = true
					}
				}

			} else {
				stopSpin = true
			}

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		
		case "down", "j":
            if m.Cursor < len(m.ThemeChoices)-1 {
                m.Cursor++
            }
		}

	case tickMsg:
		if !start {
			m.UpdateRoulette()
		}
		return m, tick

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}


func (m model) View() tea.View {
	if start {
		return startView(m)
	}

	if exit {
		if CurrentTheme.Name != "" {
			ChangeTheme()
			return tea.NewView(fmt.Sprintf("Changed Vim theme to: %s\n", CurrentTheme.Name))
		}

		return tea.NewView(fmt.Sprintf("Program exited by user"))
	}

	return rouletteView(m)
}


func startView(m model) tea.View {
	width := m.Width

	s := m.Header

	for i := range m.ThemeChoices {
		checked := " "
		if _, ok := m.ThemesSelected[i]; ok {
			checked = "X"
		}
		
		if m.Cursor == i {
			s += "> "
		}

		if i == 2 {
			s += fmt.Sprintf("Start!")
		} else {
			s += fmt.Sprintf("[%s] %s\n", checked,  m.ThemeChoices[i])
		}
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


func rouletteView(m model) tea.View {
	width := m.Width

	s := m.Header
	s += m.BlockLeft + m.DisplayText + m.BlockRight
	s += "\n"
	
	for range(m.Width) {
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


func (m *model) AddThemes() {
	text := ""
	themeIndices := []int{}
	totalLength := 0

	for _, theme := range(Config.Themes) {
		if m.ThemesSelected[1] && !theme.Favorite {
			continue
		}

		themeLength := len(theme.Name) + 4
		themeIndices = append(themeIndices, themeLength + totalLength)
		totalLength += themeLength
		text += fmt.Sprintf("| %s |", theme.Name)

		m.Themes = append(m.Themes, theme)
	}
	
	m.FullText = text
	m.Index = 0
	m.ThemeIndices = themeIndices
}


func (m *model) UpdateRoulette() {
	if speed <= 0.0 {
		if CurrentTheme.Name == "" {
			m.GetTheme()
		}

		return
	}

	m.Index += speed / 10
	m.DisplayText = ""
	wrap := 0

	if stopSpin {
		speed -= 1
	}

	if m.Index > len(m.FullText) {
		m.Index = 0
	}

	for i := range(displaySize) {
		index := m.Index + i

		if wrap > 0 {
			index = i - wrap
		}

		if index >= len(m.FullText) - 1 {
			wrap = i
		}

		m.DisplayText += string(rune(m.FullText[index]))
	}
}

func (m *model) GetTheme() {
	for i, v := range(m.ThemeIndices) {
		if v > m.Index + displaySize/2 {
			CurrentTheme = m.Themes[i]
			return
		}
	}
	CurrentTheme = m.Themes[0]
}


type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Duration(16) * time.Millisecond)
	return tickMsg(time.Now())
}

