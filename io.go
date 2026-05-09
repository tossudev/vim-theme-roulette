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
	_, err := os.Stat(ConfigPath)
	if err != nil {
		CreateConfig()
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


func CreateConfig() {
	cfg := ConfigFile{}
	
	cfg.Vimruntime = getRuntimePath()
	cfg.Themes = getThemes()

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


func getThemes() []Theme {
	VimConfig = fmt.Sprintf("%s/.vimrc", os.Getenv("HOME"))

	var themes []Theme
	var path string

	// TODO: add both theme directories at same time
	/*
	if custom {
		path = fmt.Sprintf("%s/.vim/colors", os.Getenv("HOME"))
	} else {
		path = getRuntimePath()
	}
	*/
	path = fmt.Sprintf("%s/.vim/colors", os.Getenv("HOME"))

	if path == "" {
		fmt.Println("ERR: Couldn't find themes path!")
		return []Theme{}
	}

	c, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("ERR reading dir:", path, err)
	}

	for _, entry := range c {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".vim") {
			themes = append(themes, FormatTheme(entry.Name(), true))
		}
	}
	
	return themes
}

func FormatTheme(filename string, builtin bool) Theme {
	return Theme{
		Name:		strings.TrimSuffix(filename, ".vim"),
		Color:		"\033[32m",
		Builtin:	builtin,
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
