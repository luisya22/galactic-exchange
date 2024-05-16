package uimodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luisya22/galactic-exchange/cmd/tui/store"
)

type MainModel struct {
	navbar  NavbarModel
	content ContentModel
	state   MainState
	width   int
	height  int
	store   *store.Store
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
	IsActive() bool
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

		case "esc":
			m, cmd = m.manageContentChange(msg)
		case "enter":
			if m.state == MainStateNavbarControl {
				m.state = MainStateContentControl
				m, cmd = m.manageContentChange("Activate")

			} else {
				m, cmd = m.manageContentChange(msg)
			}
		default:
			if m.state == MainStateNavbarControl {
				m, cmd = m.manageNavbarChange(msg)

			} else if m.state == MainStateContentControl {
				m, cmd = m.manageContentChange(msg)
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
			m.content = NewTradeHubModel(m.width, m.height, m.store)
		default:
			m.content = NewBlankModel()
		}
	}

	return m, cmd
}

func (m MainModel) manageContentChange(msg tea.Msg) (MainModel, tea.Cmd) {

	newContent, cmd := m.content.Update(msg)
	m.content = newContent

	if !m.content.IsActive() {
		m.state = MainStateNavbarControl
	}

	return m, cmd
}

func InitialModel(s *store.Store) MainModel {

	navbar := InitNavbar(s)
	home := HomeModel{
		store: s,
	}

	return MainModel{
		navbar:  navbar,
		content: home,
		state:   MainStateNavbarControl,
		store:   s,
	}
}
