package main

import (
    "fmt"
	"os"
	"time"
	"slices"

    tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type View int

const (
	ANSIColorReset string = "\033[0m"
	ANSIColorGreen string = "\u001B[32m"
	ANSIColorYellow string = "\u001B[33m"
	ANSIColorWhite string = "\u001B[37m"
	ANSIColorPastel string = "\u001B[210m"

	ViewStart View = iota
	ViewRoulette
	ViewExit
)

type model struct {
	Header 			string
	HeaderHue		float64
	BlockLeft 		string
	BlockRight 		string

	FullText 		string
	DisplayText 	string
	Index 			int
	Themes 			[]Theme
	ThemeIndices	[]int
	MenuChoices 	[]string
	Cursor 			int
	Selected		int

	Width 			int
	Height 			int
	CurrentView		View
}

var (
	displaySize int = 50
	speed int = 50
	stopSpin = false
	headerColorSpeed float64 = 0.002
	rouletteColors []string = []string{
		"\033[38;5;10m",
		"\033[38;5;11m",
		"\033[38;5;12m",
		"\033[38;5;13m",
		"\033[38;5;14m",
		"\033[38;5;15m",
	}
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
		Header:			"\n☆ ★ ☆ ★  Vim Theme Roulette  ☆ ★ ☆ ★\n\n",
		BlockLeft:		"░▒▓ ",
		BlockRight:		" ▓▒░",

		MenuChoices:	[]string{"All", "Favorites", "Start!", "Quit"},
		CurrentView:	ViewStart,
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
			m.CurrentView = ViewExit
			return m, tea.Quit

		case "space", "enter":
			if m.CurrentView == ViewStart {
				if m.Cursor == 2 {
					m.AddThemes()
					m.CurrentView = ViewRoulette

				} else if m.Cursor == 3 {
					m.CurrentView = ViewExit
					return m, tea.Quit

				} else {
					m.Selected = m.Cursor
				}

			} else {
				stopSpin = true
			}

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		
		case "down", "j":
            if m.Cursor < len(m.MenuChoices)-1 {
                m.Cursor++
            }
		}

	case tickMsg:
		m.HeaderHue += headerColorSpeed
		if m.HeaderHue > 1.0 {
			m.HeaderHue = 0.0
		}

		if m.CurrentView == ViewRoulette {
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
	switch m.CurrentView {
		case ViewStart:
			return startView(m)

		case ViewExit:
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

	s := Rainbow(m.HeaderHue) + m.Header + ANSIColorReset
	s += "\nSelect themes to include:\n"

	for i := range m.MenuChoices {
		prefix := ""
		suffix := ""

		if i == m.Selected {
			prefix = ANSIColorGreen
			suffix = " ✅"
		}
		
		if i == 2 {
			s += "\n"
		}
		
		if m.Cursor == i {
			prefix = ANSIColorYellow
		}

		s += fmt.Sprintf("%s%s%s%s\n", prefix, m.MenuChoices[i], suffix, ANSIColorReset)
	}

	s += "\n"

	centered := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(s)

	return tea.NewView(centered)
}


func rouletteView(m model) tea.View {
	width := m.Width

	s := Rainbow(m.HeaderHue) + m.Header + ANSIColorReset
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
		Render(s)

	return tea.NewView(centered)
}


func (m *model) AddThemes() {
	text := ""
	themeIndices := []int{}
	totalLength := 0

	for _, theme := range(Config.Themes) {
		if m.Selected == 1 && !slices.Contains(Config.Favorites, theme.Name) {
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

	/*
	colors := []string{
		"\033[38;5;10m",
		"\033[38;5;11m",
		"\033[38;5;12m",
		"\033[38;5;13m",
		"\033[38;5;14m",
		"\033[38;5;15m",
	}
	*/
	colorIndex := 0

	for i := range(displaySize) {
		for j, v := range(m.ThemeIndices) {
			if v > m.Index + i {
				colorIndex = j
				break
			}
		}
		color := rouletteColors[colorIndex % len(rouletteColors)]

		index := m.Index + i

		if wrap > 0 {
			index = i - wrap
		}

		if index >= len(m.FullText) - 1 {
			wrap = i
			continue
		}

		m.DisplayText += color + string(rune(m.FullText[index])) + ANSIColorReset
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

