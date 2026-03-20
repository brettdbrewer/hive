package tui

import "github.com/charmbracelet/lipgloss"

const (
	sidebarWidth = 22
	statusHeight = 1
	inputHeight  = 3
)

var (
	// Colors.
	accent    = lipgloss.Color("#7D56F4") // hive purple
	subtle    = lipgloss.Color("#626262")
	highlight = lipgloss.Color("#F25D94")
	dimText   = lipgloss.Color("#888888")
	white     = lipgloss.Color("#FAFAFA")

	// Sidebar.
	sidebarStyle = lipgloss.NewStyle().
			Width(sidebarWidth).
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			BorderForeground(subtle).
			Padding(1, 1)

	sidebarTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1)

	sidebarItem = lipgloss.NewStyle().
			Foreground(white)

	sidebarItemActive = lipgloss.NewStyle().
				Foreground(accent).
				Bold(true)

	sidebarSection = lipgloss.NewStyle().
			Foreground(dimText).
			MarginTop(1)

	// Chat.
	chatStyle = lipgloss.NewStyle().
			Padding(0, 1)

	humanMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89CFF0")).
			Bold(true)

	mindMsg = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	errorMsg = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true)

	inputBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	// Status bar.
	statusStyle = lipgloss.NewStyle().
			Foreground(dimText).
			Padding(0, 1)

	statusAccent = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	// Title bar.
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(accent).
			Padding(0, 1)
)
