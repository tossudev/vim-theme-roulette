package main

import (
	"strings"
	"math/rand"
)

type Theme struct {
	name string
	color string
	builtin bool
}

var Themes []Theme
var CurrentTheme Theme


func AddTheme(filename string, builtin bool) {
	theme := Theme{
		name:		strings.TrimSuffix(filename, ".vim"),
		color:		"\033[32m",
		builtin:	builtin,
	}

	Themes = append(Themes, theme)
}


func ShuffleThemes() {
	for i := range Themes {
		j := rand.Intn(i + 1)
		Themes[i], Themes[j] = Themes[j], Themes[i]

	}
}
