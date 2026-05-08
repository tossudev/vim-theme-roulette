package main

import (
	"os"
	"os/exec"
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
	cmd := exec.Command("bash",
		"-c",
		`vim -T dumb --cmd 'exe "set t_cm=\<C-M>"|echo $VIMRUNTIME|quit' | tr -d '\015'`,
	)
	
	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Runtime path error:", err)
		fmt.Println(err.Error())
	}

	return strings.TrimSpace(out.String()) + "/colors/"
}
