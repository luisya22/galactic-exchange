package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luisya22/galactic-exchange/cmd/tui/store"
	uimodel "github.com/luisya22/galactic-exchange/cmd/tui/uimodels"
)

func main() {
	s := store.NewStore(3)
	p := tea.NewProgram(uimodel.InitialModel(s), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
