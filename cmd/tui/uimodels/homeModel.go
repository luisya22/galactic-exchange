package uimodel

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	width  int
	height int
}

func (h HomeModel) Init() tea.Cmd {
	return nil
}

func (h HomeModel) SetSize(width, height int) (ContentModel, tea.Cmd) {
	h.width = width
	h.height = height

	return h, nil
}

func (h HomeModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return h.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	return h, nil
}

func (h HomeModel) View() string {
	column1 := lipgloss.NewStyle().
		Width((h.width-20)/3).
		Render(h.getSquadContent()) + "\n" // Add newline

	column2 := lipgloss.NewStyle().
		Width((h.width - 20) / 3).
		Render(h.getMainContent()) // Add newline

	column3 := lipgloss.NewStyle().
		Width((h.width - 20) / 3).
		Render(h.getEventContent()) // Add newline

	content := lipgloss.JoinHorizontal(lipgloss.Top, column1, column2, column3)

	return lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		Render(content)
}

func NewHomeModel(width, height int) ContentModel {
	return HomeModel{
		width:  width,
		height: height,
	}
}

func (h HomeModel) ID() string {
	return TabHome
}

func (h HomeModel) getSquadContent() string {
	activeSquadsTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("Active Squads")

	activeSquads := fmt.Sprintf(
		"%s\n%s\n%s\n",
		"Squad 1 - HARVESTING - 1:02 left",
		"Squad 2 - RETURNING - 0:30 left",
		"Squad 3 - ESCAPING FROM PIRATES - unknown",
	)

	activeSquadsInfo := fmt.Sprintf("%s\n%s\n", activeSquadsTitle, activeSquads)
	activeSquadsInfo = lipgloss.NewStyle().
		MarginBottom(2).
		Render(activeSquadsInfo)

	availableSquadsTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("Available Squads")

	availableSquads := fmt.Sprintf(
		"%s\n%s\n%s\n",
		"Squad 4 [ Send to Mission ]",
		"Squad 5 [ Send to Mission ]",
		"Squad 6 [ Send to Mission ]",
	)

	availableSquadsInfo := fmt.Sprintf("%s\n%s\n", availableSquadsTitle, availableSquads)

	return activeSquadsInfo + "\n" + availableSquadsInfo
}

func (h HomeModel) getMainContent() string {
	recentTradesTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("Recente Trades")

	recenteTrades := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n",
		"35x Iron sold to Planet XYZ",
		"Iron sold out on Sector A",
		"43x Water sold to Planet YYY",
		"100x Meteor Rock sold to Planet ABC",
	)

	tradesInfo := fmt.Sprintf("%s\n%s\n", recentTradesTitle, recenteTrades)
	tradesInfo = lipgloss.NewStyle().
		MarginBottom(2).
		Render(tradesInfo)

	rbTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("Research and Buildings")

	progress := progress.New(progress.WithScaledGradient("#32C7C0", "#A3FBF7"))

	rb := fmt.Sprintf(
		"%s %s\n%s %s\n%s %s\n",
		"Ship Construction    ",
		progress.ViewAs(.70),
		"Building Construction",
		progress.ViewAs(.50),
		"Research             ",
		progress.ViewAs(.10),
	)

	rbInfo := fmt.Sprintf("%s\n%s", rbTitle, rb)

	return tradesInfo + "\n" + rbInfo
}

func (h HomeModel) getEventContent() string {
	recentTradesTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32C7C0")).
		Render("Recent Events")

	recenteTrades := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n",
		"Senate has approved new oxygen sanctions",
		"Iron value increased 20% over the last months on Sector AZ24",
		"Rumors says that a new saction would be approved by Senate",
		"A recurring anomaly in sector XR45 is causing corporations to forfeit valuable assets",
	)

	tradesInfo := fmt.Sprintf("%s\n%s\n", recentTradesTitle, recenteTrades)
	tradesInfo = lipgloss.NewStyle().
		MarginBottom(2).
		Render(tradesInfo)

	return tradesInfo
}
