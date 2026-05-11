package main

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
	"encoding/json"
)

var (
	VimConfig string
	ConfigPath string = "config.json"
	Config ConfigFile
)

type ConfigFile struct {
	Vimruntime string
	Themes []Theme
}


func FetchConfig() {
	VimConfig = fmt.Sprintf("%s/.vimrc", os.Getenv("HOME"))
	
	_, err := os.Stat(ConfigPath)
	if err != nil {
		createConfig()
		return
	}

	jsonData, err2 := os.ReadFile(ConfigPath)
	if err2 != nil {
		fmt.Println("ERR reading file:", ConfigPath, err)
	}
	
	cfgFile := ConfigFile{}
	err3 := json.Unmarshal(jsonData, &cfgFile)
	if err3 != nil {
		fmt.Println("JSON Unmarshal ERR:", err3)
	}
	
	Config = cfgFile
}


func ChangeTheme() {
	contents, err := os.ReadFile(VimConfig)
	if err != nil {
		fmt.Println("ERR reading file:", VimConfig, err)
	}
	
	lines := strings.Split(string(contents), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "colorscheme ") {
			lines[i] = fmt.Sprintf("colorscheme %s", CurrentTheme.Name)
			break
		}
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

	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("JSON Marshal ERR:", err)
	}
	
	err2 := os.WriteFile(ConfigPath, []byte(cfgJson), 0644)
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
		Color:		"\033[32m",
		Custom:		custom,
		Favorite:	false,
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
