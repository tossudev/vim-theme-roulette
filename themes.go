package main

import (
	//"strings"
	"math/rand"
)

type Theme struct {
	Name string
	Color string
	Custom bool
}

var CurrentTheme Theme

/*
func AddTheme(filename string, builtin bool) {
	theme := Theme{
		name:		strings.TrimSuffix(filename, ".vim"),
		color:		"\033[32m",
		builtin:	builtin,
	}

	Themes = append(Themes, theme)
}
*/

func ShuffleThemes() {
	for i := range Config.Themes {
		j := rand.Intn(i + 1)
		Config.Themes[i], Config.Themes[j] = Config.Themes[j], Config.Themes[i]

	}
}
