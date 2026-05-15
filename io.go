package main

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
	"github.com/pelletier/go-toml/v2"
)

var (
	VimConfig string
	ConfigPath string = "config.toml"
	Config ConfigFile
)

type ConfigFile struct {
	Vimruntime string
	Favorites []string
	Themes []Theme
}


func FetchConfig() {
	VimConfig = fmt.Sprintf("%s/.vimrc", os.Getenv("HOME"))
	
	_, err := os.Stat(ConfigPath)
	if err != nil {
		createConfig()
		return
	}

	tomlData, err2 := os.ReadFile(ConfigPath)
	if err2 != nil {
		fmt.Println("ERR reading file:", ConfigPath, err)
	}
	
	cfgFile := ConfigFile{}
	err3 := toml.Unmarshal(tomlData, &cfgFile)
	if err3 != nil {
		fmt.Println("TOML Unmarshal ERR:", err3)
	}
	
	Config = cfgFile
}


func ChangeTheme() {
	contents, err := os.ReadFile(VimConfig)
	if err != nil {
		fmt.Println("ERR reading file:", VimConfig, err)
	}
	
	lines := strings.Split(string(contents), "\n")
	themeLine := fmt.Sprintf("colorscheme %s", CurrentTheme.Name)
	hasTheme := false

	for i, line := range lines {
		if strings.HasPrefix(line, "colorscheme ") {
			lines[i] = themeLine
			hasTheme = true
			break
		}
	}

	if !hasTheme {
		lines = append(lines, themeLine)
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(VimConfig, []byte(output), 0644)
	if err != nil {
		fmt.Println("ERR writing to file:", VimConfig, err)
	}
}


func createConfig() {
	cfg := ConfigFile{}
	
	cfg.Vimruntime = getRuntimePath()
	cfg.Themes = getThemes(cfg.Vimruntime)

	cfgToml, err := toml.Marshal(cfg)
	if err != nil {
		fmt.Println("TOML Marshal ERR:", err)
	}
	
	err2 := os.WriteFile(ConfigPath, []byte(cfgToml), 0644)
	if err2 != nil {
		fmt.Println("ERR writing to file:", ConfigPath, err)
	}

	Config = cfg
}


func getThemes(vimruntime string) []Theme {
	var themes []Theme
	var path string
	var custom bool

	// TODO:
	// Make this look better lol
	for i := range(2) {
		if i == 0 {
			custom = true
			path = fmt.Sprintf("%s/.vim/colors", os.Getenv("HOME"))
		} else {
			if vimruntime == "" {
				return themes
			}

			custom = false
			path = vimruntime
		}

		c, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("ERR reading dir:", path, err)
		}

		for _, entry := range c {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vim") {
				themes = append(themes, formatTheme(entry.Name(), custom))
			}
		}
	}
	
	return themes
}


func formatTheme(filename string, custom bool) Theme {
	return Theme{
		Name:		strings.TrimSuffix(filename, ".vim"),
		Color:		"Red",
		Custom:		custom,
	}
}


func getRuntimePath() string {

	// TODO:
	// Do this without opening bash since its slow as shit
	// (maybe 2 seconds to boot, just opening vim is pretty much instant)
	cmd := exec.Command("bash",
		"-c",
		`vim -T dumb --cmd 'exe "set t_cm=\<C-M>"|echo $VIMRUNTIME|quit' | tr -d '\015'`,
	)

	/*
	exec.Command(
		"ex",
		"-c",
		"call writefile([$VIMRUNTIME], 'vimruntime.txt')|q!",
	)
	*/
	
	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Runtime path error:", err)
		fmt.Println(err.Error())
	}

	return strings.TrimSpace(out.String()) + "/colors/"
}
