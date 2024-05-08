package uimodel

import tea "github.com/charmbracelet/bubbletea"

type BlankModel struct {
}

func (h BlankModel) Init() tea.Cmd {
	return nil
}

func (h BlankModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	return h, nil
}

func (h BlankModel) View() string {
	return "Blank Model"
}

func NewBlankModel() ContentModel {
	return BlankModel{}
}

func (h BlankModel) ID() string {
	return TabTradeHub
}
