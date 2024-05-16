package uimodel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luisya22/galactic-exchange/cmd/tui/store"
)

type TradeHubModel struct {
	width     int
	height    int
	tabs      []TradeTab
	activeTab int
	isActive  bool
	state     TradeState
	store     *store.Store
}

type TradeState int

const (
	TradeStateTopMenu TradeState = iota
	TradeStateSubMenu
)

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

func (t TradeHubModel) IsActive() bool {
	return t.isActive
}

func (t TradeHubModel) SetSize(width, height int) (ContentModel, tea.Cmd) {
	t.width = width
	t.height = height

	return t, nil
}

func (t TradeHubModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {

	if !t.IsActive() && msg == "Activate" {
		t.isActive = true
		t.state = TradeStateTopMenu
		return t, nil
	}

	if !t.isActive {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return t.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if t.activeTab > 0 && t.state == TradeStateTopMenu {
				t.activeTab--
			}
		case "right", "l":
			if t.activeTab < len(t.tabs)-1 {
				t.activeTab++
			}
		case "enter":
			// if t.state == TradeStateTopMenu {
			//
			// }
		case "esc":
			// TODO: Only if child is correct value
			if t.state == TradeStateSubMenu {
				t.state = TradeStateTopMenu
			} else {
				t.isActive = false
			}
		}
	}
	return t, nil
}

func (t TradeHubModel) View() string {
	tabs := lipgloss.NewStyle().
		Render(t.getTabs())

	h := t.height - t.store.NavBarHeight - navbarMarginTop - topNavbarMarginBottom - navbarMarginBottom - 5

	content := lipgloss.NewStyle().
		Height(h).
		Padding(0, 1).
		Width(t.width - 5).
		Border(lipgloss.NormalBorder()).
		Render()

	return tabs + content
}

func (t TradeHubModel) ID() string {
	return TabTradeHub
}

func NewTradeHubModel(width, height int, s *store.Store) ContentModel {
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
		store: s,
	}
}

func (t TradeHubModel) getTabs() string {

	bottomLeftBorder := "┘"
	if t.activeTab == 0 {
		bottomLeftBorder = "|"
	}

	activeBorder := activeTabBorder
	activeBorder.BottomLeft = bottomLeftBorder

	activeTab := lipgloss.NewStyle().
		Border(activeBorder, true).
		Padding(0, 1)

	if t.isActive {
		activeTab.Foreground(primaryColor)
	}

	tabsArr := []string{}
	for i, tabStr := range t.tabs {
		if i == t.activeTab {
			tabsArr = append(tabsArr, activeTab.Render(string(tabStr)))
			continue
		}

		tabsArr = append(tabsArr, tabStyle.Render(string(tabStr)))
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabsArr...)

	gap := tabGapStyle.Render(strings.Repeat(" ", max(0, t.width-5-lipgloss.Width(tabs)-2)))
	tabs = lipgloss.JoinHorizontal(lipgloss.Bottom, tabs, gap)

	return tabs

}
