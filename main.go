package main

import (
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}

type Model struct {
	Entries []os.DirEntry
	Cursor  cursor.Model

	cursorPos int
}

func initialModel() Model {
	var model Model
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	model.Entries, err = os.ReadDir(cwd)
	if err != nil {
		panic(err)
	}

	model.Cursor = cursor.New()
	model.Cursor.Focus()

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
			return model, tea.Quit
		case "up", "k":
			if model.cursorPos > 0 {
				model.cursorPos--
			}
		case "down", "j":
			if model.cursorPos <= len(model.Entries) {
				model.cursorPos++
			}
		}
	}
	return model, nil
}

func (model Model) View() string {
	res := "../\n"
	res += model.Cursor.View()
	for _, entry := range model.Entries {
		res += entry.Name() + "\n"
	}
	return res
}
