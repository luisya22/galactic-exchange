package uimodel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luisya22/galactic-exchange/cmd/tui/store"
)

type TradeHubMarketdataModel struct {
	width    int
	height   int
	isActive bool
	store    *store.Store
}

func (t TradeHubMarketdataModel) Init() tea.Cmd {
	return nil
}

func (t TradeHubMarketdataModel) Update(tea.Msg) (ContentModel, tea.Cmd) {
	return t, nil
}

func (t TradeHubMarketdataModel) View() string {
	return "Hello there 2!"
}

func (t TradeHubMarketdataModel) ID() string {
	return TradeTabMarketData
}

func (t TradeHubMarketdataModel) IsActive() bool {
	return t.isActive
}

func NewTradeHubMarketData(width int, height int, store *store.Store) ContentModel {
	return TradeHubMarketdataModel{
		width:  width,
		height: height,
		store:  store,
	}
}
