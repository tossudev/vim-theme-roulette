package main

import (
	"os"
	"fmt"
	"strings"
)

var (
	VimRuntime string = "/usr/share/vim/"
	CachedThemes []string
	VimConfig string
)


func GetThemesLocal() {
	VimConfig = fmt.Sprintf("%s/.vimrc", os.Getenv("HOME"))

	runtimePath := getRuntimePath()
	if runtimePath == "" {
		fmt.Println("ERR: Couldn't find VIMRUNTIME!")
		return
	}

	c, err := os.ReadDir(runtimePath)
	if err != nil {
		fmt.Println("ERR reading dir:", runtimePath, err)
	}

	for _, entry := range c {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vim") {
			AddTheme(entry.Name(), true)

			//CachedThemes = append(CachedThemes, strings.TrimSuffix(entry.Name(), ".vim"))
		}
	}

	//fmt.Println(Themes)
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
		fmt.Println("ERR:", err)
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
