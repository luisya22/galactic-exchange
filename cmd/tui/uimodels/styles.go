package uimodel

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Sizes
	navbarHeight          = 5
	topNavbarMarginBottom = 1
	navbarMarginTop       = 1
	navbarMarginBottom    = 2

	// COLORS
	primaryColor = lipgloss.Color("#32C7C0")

	textPrimaryColorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#32C7C0"))

	// PADDING & MARGIN

	mainPaddingLeftStyle = lipgloss.NewStyle().PaddingLeft(2)

	// Navbar

	topNavbarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			MarginBottom(topNavbarMarginBottom)

	navbarStyle = lipgloss.NewStyle().
			MarginBottom(navbarMarginBottom).
			MarginTop(navbarMarginTop)

	// TABS

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tabStyle = lipgloss.NewStyle().
			Border(tabBorder, true).
			Padding(0, 1)

	tabGapStyle = lipgloss.NewStyle().
			Border(tabBorder, true).
			Padding(0, 1).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
)

// Justify content between two strings
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
