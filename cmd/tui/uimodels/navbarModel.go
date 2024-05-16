package uimodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luisya22/galactic-exchange/cmd/tui/store"
)

type NavbarModel struct {
	activeTab int
	tabs      []NavbarTab
	width     int
	height    int
	store     *store.Store
}

type NavbarTab string

const (
	TabHome          = "Home"
	TabTradeHub      = "Trade Hub"
	TabCommandCenter = "Command Center"
	TabSpaceport     = "Spaceport"
	TabGalacticHub   = "Galactic Hub"
	TabSettings      = "Settings"
)

func (n NavbarModel) Init() tea.Cmd {
	return nil
}

func (n NavbarModel) View() string {

	greetings := "Welcome Back"

	userNameStyle := textPrimaryColorStyle.Render("Darthtrooper")

	leftSide := mainPaddingLeftStyle.Render(fmt.Sprintf("%s, %s", greetings, userNameStyle))

	rightSide := ""
	for i, tab := range n.tabs {
		if i == n.activeTab {
			rightSide += textPrimaryColorStyle.
				Render(fmt.Sprintf(" %s ", tab))
		} else {
			rightSide += fmt.Sprintf(" %s ", tab)
		}

	}

	rightSide = mainPaddingLeftStyle.Render(rightSide)

	topNavbar := justifyBetween(leftSide, rightSide, n.width)
	topNavbar = topNavbarStyle.Render(topNavbar)

	bottomNavbar := n.getBottomNavbar()

	return navbarStyle.Render(fmt.Sprintf("%s\n%s", topNavbar, bottomNavbar))
}

func (n NavbarModel) getBottomNavbar() string {
	year := textPrimaryColorStyle.Render("2024")
	month := textPrimaryColorStyle.Render("3")
	day := textPrimaryColorStyle.Render("4")

	date := fmt.Sprintf("YEAR %s - MONTH %s - DAY %s", year, month, day)
	credits := "20,234 CR"

	bottomBar := justifyBetween(date, credits, n.width)

	return bottomBar
}

func (n NavbarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return n.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if n.activeTab > 0 {
				n.activeTab--
			}
		case "right", "l":
			if n.activeTab < len(n.tabs)-1 {
				n.activeTab++
			}
		}
	}

	return n, nil
}

func InitNavbar(s *store.Store) NavbarModel {
	return NavbarModel{
		activeTab: 0,
		tabs:      []NavbarTab{TabHome, TabTradeHub, TabCommandCenter, TabSpaceport, TabGalacticHub, TabSettings},
		store:     s,
	}

}

func (n NavbarModel) SetSize(width, height int) (tea.Model, tea.Cmd) {
	n.width = width
	n.height = height

	return n, nil
}
