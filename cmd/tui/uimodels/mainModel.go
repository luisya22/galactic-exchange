package uimodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	navbar  NavbarModel
	content ContentModel
	state   MainState
	width   int
	height  int
}

type MainState int

const (
	MainStateNavbarControl MainState = iota
	MainStateContentControl
)

type ContentModel interface {
	Update(msg tea.Msg) (ContentModel, tea.Cmd)
	View() string
	ID() string
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		newNavBar, _ := m.navbar.Update(msg)
		m.navbar = newNavBar.(NavbarModel)

		newContent, _ := m.content.Update(msg)
		m.content = newContent
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			if m.state == MainStateNavbarControl {
				m, cmd = m.manageNavbarChange(msg)

			} else if m.state == MainStateContentControl {
				// var content tea.Model
				// content, cmd = m.manageContentChange(msg)
				// m.content = content
			}
		}
	default:
	}

	return m, cmd
}

func (m MainModel) View() string {

	str := fmt.Sprintf("%s\n%s", m.navbar.View(), m.content.View())

	return str
}

func (m MainModel) manageNavbarChange(msg tea.Msg) (MainModel, tea.Cmd) {
	newNavBar, cmd := m.navbar.Update(msg)
	m.navbar = newNavBar.(NavbarModel)

	activeTab := m.navbar.tabs[m.navbar.activeTab]

	if m.content.ID() != string(activeTab) {
		switch activeTab {
		case TabHome:
			m.content = NewHomeModel(m.width, m.height)
		case TabTradeHub:
			m.content = NewTradeHubModel(m.width, m.height)
		default:
			m.content = NewBlankModel()
		}
	}

	return m, cmd
}

func (m MainModel) manageContentChange(msg tea.Msg) (MainModel, tea.Cmd) {
	return m, nil

}

func InitialModel() MainModel {

	navbar := InitNavbar()
	home := HomeModel{}

	return MainModel{
		navbar:  navbar,
		content: home,
	}
}
