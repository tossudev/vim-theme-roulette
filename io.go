package main

import (
	"os"
	"fmt"
	"strings"
)

var (
	VimRuntime string = "/usr/share/vim/"
	VimConfig string
)


func GetThemesLocal(builtin bool) {
	VimConfig = fmt.Sprintf("%s/.vimrc", os.Getenv("HOME"))

	var path string

	if builtin {
		path = getRuntimePath()
	} else {
		path = fmt.Sprintf("%s/.vim/colors", os.Getenv("HOME"))
	}

	if path == "" {
		fmt.Println("ERR: Couldn't find themes path!")
		return
	}

	c, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("ERR reading dir:", path, err)
	}

	for _, entry := range c {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vim") {
			AddTheme(entry.Name(), true)
		}
	}
}


func ChangeTheme() {
	contents, err := os.ReadFile(VimConfig)
	if err != nil {
		fmt.Println("ERR reading file:", VimConfig, err)
	}
	
	lines := strings.Split(string(contents), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "colorscheme ") {
			lines[i] = fmt.Sprintf("colorscheme %s", CurrentTheme.name)
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(VimConfig, []byte(output), 0644)
	if err != nil {
		fmt.Println("ERR writing to file:", VimConfig, err)
	}
}


func getRuntimePath() string {
	runtimePath := VimRuntime
	dirContents, err := os.ReadDir(VimRuntime)
	if err != nil {
		fmt.Println("ERR reading dir:", VimRuntime, err)
	}

	for _, entry := range dirContents {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "vim") {
			runtimePath += fmt.Sprintf("%s/colors/", entry.Name())
			return runtimePath
		}
	}

	return ""
}
