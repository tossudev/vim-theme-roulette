package main

import (
	"strings"
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
