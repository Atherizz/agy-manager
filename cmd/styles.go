package cmd

import "github.com/charmbracelet/lipgloss"

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)