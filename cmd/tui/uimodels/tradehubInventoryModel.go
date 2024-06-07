package uimodel

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luisya22/galactic-exchange/cmd/tui/store"
)

type TradeHubInventoryModel struct {
	width          int
	height         int
	inventoryTable table.Model
	isActive       bool
	store          *store.Store
}

func (t TradeHubInventoryModel) Init() tea.Cmd {
	return nil
}

func (t TradeHubInventoryModel) SetSize(width, height int) (ContentModel, tea.Cmd) {
	t.width = width
	t.height = height

	return t, nil
}

func (t TradeHubInventoryModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return t.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		t.inventoryTable, cmd = t.inventoryTable.Update(msg)
	}

	return t, cmd
}

func (t TradeHubInventoryModel) View() string {
	return t.inventoryTable.View()
}

func (t TradeHubInventoryModel) ID() string {
	return TradeTabInventory
}

func (t TradeHubInventoryModel) IsActive() bool {
	return t.isActive
}

func NewTradeHubInventory(width, height int, s *store.Store) ContentModel {
	t := TradeHubInventoryModel{
		width:    width,
		height:   height,
		isActive: true,
		store:    s,
	}

	t.inventoryTable = t.initializeInventoryTable()

	return t
}

// Inventory Table
func (t TradeHubInventoryModel) initializeInventoryTable() table.Model {
	columns := []table.Column{
		{Title: "Resource", Width: 20},
		{Title: "Quantity", Width: 20},
	}

	rows := []table.Row{
		{"Neutronium Ore", "104435"},
		{"Solar Plase", "400443"},
	}

	tb := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(t.store.ContentHeight-10),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	tb.SetStyles(s)

	return tb

}
