package main

import (
	//"strings"
	"math/rand"
)

type Theme struct {
	Name string
	Color string
	Custom bool
	Favorite bool
}

var CurrentTheme Theme


func ShuffleThemes() {
	for i := range Config.Themes {
		j := rand.Intn(i + 1)
		Config.Themes[i], Config.Themes[j] = Config.Themes[j], Config.Themes[i]

	}
}
