package uimodel

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NavbarModel struct {
	activeTab int
	tabs      []NavbarTab
	width     int
	height    int
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

	userNameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("Darthtrooper")

	leftSide := lipgloss.NewStyle().
		PaddingLeft(2).
		Render(fmt.Sprintf("%s, %s", greetings, userNameStyle))

	rightSide := ""
	for i, tab := range n.tabs {
		if i == n.activeTab {
			rightSide += lipgloss.NewStyle().
				Foreground(lipgloss.Color("#32C7C0")).
				Render(fmt.Sprintf(" %s ", tab))
		} else {
			rightSide += fmt.Sprintf(" %s ", tab)
		}

	}

	rightSide = lipgloss.NewStyle().
		PaddingRight(2).
		Render(rightSide)

	topNavbar := justifyBetween(leftSide, rightSide, n.width)
	topNavbar = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		MarginBottom(1).
		MarginTop(1).
		Render(topNavbar)

	bottomNavbar := n.getBottomNavbar()

	return lipgloss.NewStyle().
		MarginBottom(2).
		Render(fmt.Sprintf("%s\n%s", topNavbar, bottomNavbar))
}

func (n NavbarModel) getBottomNavbar() string {
	year := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("2024")

	month := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("3")

	day := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("4")

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

func InitNavbar() NavbarModel {
	return NavbarModel{
		activeTab: 0,
		tabs:      []NavbarTab{TabHome, TabTradeHub, TabCommandCenter, TabSpaceport, TabGalacticHub, TabSettings},
	}

}

func (n NavbarModel) SetSize(width, height int) (tea.Model, tea.Cmd) {
	n.width = width
	n.height = height

	return n, nil
}

func justifyBetween(str1, str2 string, totalWidth int) string {
	str1Len := lipgloss.Width(str1)
	str2Len := lipgloss.Width(str2)

	strSum := str1Len + str2Len

	spaceCount := totalWidth - strSum
	if spaceCount < 0 {
		spaceCount = 0
	}

	spaces := strings.Repeat(" ", spaceCount)

	return str1 + spaces + str2
}
