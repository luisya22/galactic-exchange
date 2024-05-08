package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	uimodel "github.com/luisya22/galactic-exchange/cmd/tui/uimodels"
)

func main() {
	p := tea.NewProgram(uimodel.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
