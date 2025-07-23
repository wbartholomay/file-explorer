package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) != 2 {
		panic("incorrect args")
	}
	outFile := os.Args[1]

	file, err := os.Create(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file: %v\n", err)
		os.Exit(1)
	}

	program := tea.NewProgram(initialModel(file))
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}

type Model struct {
	Entries []os.DirEntry
	Pipe    *os.File

	command   string
	cursorPos int
}

func initialModel(pipe *os.File) Model {
	var model Model
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	model.command = "cd ."
	model.Pipe = pipe
	model.Entries, err = os.ReadDir(cwd)
	if err != nil {
		panic(err)
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			fmt.Fprintln(model.Pipe, model.command)
			return model, tea.Quit
		case "enter":
			model = model.HandleEnter()
		case "up", "k":
			if model.cursorPos > 0 {
				model.cursorPos--
			}
		case "down", "j":
			if model.cursorPos < len(model.Entries) {
				model.cursorPos++
			}
		}
	}
	return model, nil
}

func (model Model) HandleEnter() Model {
	var fileName string
	var isDir bool
	if model.cursorPos == 0 {
		fileName = ".."
		isDir = true
	} else {
		entry := model.Entries[model.cursorPos-1]
		fileName = entry.Name()
		isDir = entry.IsDir()
	}
	if isDir {
		model.command += "/" + fileName
		err := os.Chdir(fileName)
		if err != nil {
			// TODO add error handling
			panic(err)
		}
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		model.Entries, err = os.ReadDir(cwd)
		if err != nil {
			panic(err)
		}
		model.cursorPos = 0
	}
	return model
}

func (model Model) View() string {
	res := ""
	if model.cursorPos == 0 {
		res += ">"
	}
	res += "../\n"
	for i, entry := range model.Entries {
		if model.cursorPos == i+1 {
			res += ">"
		}
		res += entry.Name() + "\n"
	}
	return res
}
