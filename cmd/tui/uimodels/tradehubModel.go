package uimodel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TradeHubModel struct {
	width     int
	height    int
	tabs      []TradeTab
	activeTab int
}

type TradeTab string

const (
	TradeTabInventory   = "Inventory"
	TradeTabMarketData  = "Market Data"
	TradeTabOrders      = "Trade Orders"
	TradeTabAgreements  = "Agreements"
	TradeTabSanctions   = "Sanctions"
	TradeTabStockMarket = "Stock Market"
)

func (t TradeHubModel) Init() tea.Cmd {
	return nil
}

func (t TradeHubModel) SetSize(width, height int) (ContentModel, tea.Cmd) {
	t.width = width
	t.height = height

	return t, nil
}

func (t TradeHubModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return t.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		}
	}
	return t, nil
}

func (t TradeHubModel) View() string {
	tabs := lipgloss.NewStyle().
		Render(t.getTabs())

	content := lipgloss.NewStyle().
		Height(40).
		Padding(0, 1).
		Width(t.width - 5).
		Border(lipgloss.NormalBorder()).
		Render()

	return tabs + content
}

func (t TradeHubModel) ID() string {
	return TabTradeHub
}

func NewTradeHubModel(width, height int) ContentModel {
	return TradeHubModel{
		activeTab: 0,
		width:     width,
		height:    height,
		tabs: []TradeTab{
			TradeTabInventory,
			TradeTabMarketData,
			TradeTabOrders,
			TradeTabAgreements,
			TradeTabSanctions,
			TradeTabStockMarket,
		},
	}
}

func (t TradeHubModel) getTabs() string {

	bottomLeftBorder := "┘"
	if t.activeTab == 0 {
		bottomLeftBorder = "|"
	}

	activeTabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  bottomLeftBorder,
		BottomRight: "└",
	}

	tabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab := lipgloss.NewStyle().
		Border(tabBorder, true).
		Padding(0, 1)

	activeTab := lipgloss.NewStyle().
		Border(activeTabBorder, true).
		Padding(0, 1)

	tabGap := lipgloss.NewStyle().
		Border(tabBorder, true).
		Padding(0, 1).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	tabsArr := []string{}
	for i, tabStr := range t.tabs {
		if i == t.activeTab {
			tabsArr = append(tabsArr, activeTab.Render(string(tabStr)))
			continue
		}

		tabsArr = append(tabsArr, tab.Render(string(tabStr)))
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabsArr...)

	gap := tabGap.Render(strings.Repeat(" ", max(0, t.width-5-lipgloss.Width(tabs)-2)))
	tabs = lipgloss.JoinHorizontal(lipgloss.Bottom, tabs, gap)

	return tabs

}
