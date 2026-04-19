package main

import (
	"os"
	"fmt"
	"strings"
)

const (
	VimRuntime string = "/usr/share/vim/vim91/colors/"
)

var (
	CachedThemes []string
	VimConfig string
)


func GetThemesLocal() {
	VimConfig = fmt.Sprintf("%s/.vimrc", os.Getenv("HOME"))

	c, err := os.ReadDir(VimRuntime)
	if err != nil {
		fmt.Println("ERR:", err)
	}

	for _, entry := range c {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vim") {
			CachedThemes = append(CachedThemes, strings.TrimSuffix(entry.Name(), ".vim"))
		}
	}
}


func ChangeTheme() {
	contents, err := os.ReadFile(VimConfig)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	
	lines := strings.Split(string(contents), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "colorscheme ") {
			lines[i] = fmt.Sprintf("colorscheme %s", RolledTheme)
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(VimConfig, []byte(output), 0644)
	if err != nil {
		fmt.Println("ERR:", err)
	}
}

