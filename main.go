package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
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
	Entries   []os.DirEntry
	Pipe      *os.File
	TextInput textinput.Model

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
	model.TextInput = textinput.New()

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
	shouldQuit := false
	if model.TextInput.Focused() {
		model, shouldQuit = model.UpdateTextInput(msg)
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				shouldQuit = true
			case "enter":
				model, shouldQuit = model.HandleEnter()
			case "up", "k":
				if model.cursorPos > 0 {
					model.cursorPos--
				}
			case "down", "j":
				if model.cursorPos < len(model.Entries) {
					model.cursorPos++
				}
			case "n", "t":
				model.TextInput.Focus()
			}
		}
	}

	if shouldQuit {
		fmt.Fprintln(model.Pipe, model.command)
		return model, tea.Quit
	} else {
		return model, nil
	}
}

func (model Model) UpdateTextInput(msg tea.Msg) (Model, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			fileName := model.TextInput.Value()
			model.command += "\nnvim " + fileName
			return model, true
		case "esc":
			model.TextInput.SetValue("")
			model.TextInput.Blur()
		case "ctrl+c":
			return model, true
		default:
			model.TextInput, _ = model.TextInput.Update(msg)
		}
	}
	return model, false
}

func (model Model) HandleEnter() (Model, bool) {
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
		return model, false
	} else {
		model.command += "\nnvim " + fileName
		return model, true
	}
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
	res += "\n\n" + model.TextInput.View()
	return res
}
