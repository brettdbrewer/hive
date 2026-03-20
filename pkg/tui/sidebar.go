package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SidebarSection identifies which section of the sidebar is selected.
type SidebarSection int

const (
	SectionAgents SidebarSection = iota
	SectionActivity
	SectionTasks
)

// Sidebar shows agents, activity, and task summary.
type Sidebar struct {
	agents        []AgentInfo
	activity      []ActivityEntry
	conversations []ConversationInfo
	tasks         TaskSummary
	cursor        int
	section       SidebarSection
	width         int
	height        int
	focused       bool
}

// TaskSummary is a quick count for the sidebar.
type TaskSummary struct {
	Open       int
	InProgress int
	Done       int
}

// NewSidebar creates a new sidebar.
func NewSidebar() Sidebar {
	return Sidebar{
		agents: []AgentInfo{
			{Role: "mind", Name: "Mind", Status: "ready"},
		},
		width: sidebarWidth,
	}
}

// SetSize updates the sidebar height.
func (s *Sidebar) SetSize(h int) {
	s.height = h
}

// SetAgents updates the agent list.
func (s *Sidebar) SetAgents(agents []AgentInfo) {
	s.agents = agents
}

// SetActivity updates the activity list.
func (s *Sidebar) SetActivity(activity []ActivityEntry) {
	s.activity = activity
}

// SetConversations updates the conversation history list.
func (s *Sidebar) SetConversations(convs []ConversationInfo) {
	s.conversations = convs
}

// SetTasks updates the task summary.
func (s *Sidebar) SetTasks(summary TaskSummary) {
	s.tasks = summary
}

// Focus gives keyboard focus to the sidebar.
func (s *Sidebar) Focus() { s.focused = true }

// Blur removes keyboard focus.
func (s *Sidebar) Blur() { s.focused = false }

// SelectedAgent returns the currently selected agent role.
func (s *Sidebar) SelectedAgent() string {
	if s.section == SectionAgents && s.cursor < len(s.agents) {
		return s.agents[s.cursor].Role
	}
	return "mind"
}

// Update handles sidebar key events.
func (s Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			maxItems := len(s.agents)
			if s.cursor < maxItems-1 {
				s.cursor++
			}
		case "enter":
			return s, func() tea.Msg {
				return selectAgentMsg{Agent: s.SelectedAgent()}
			}
		case "tab":
			return s, func() tea.Msg { return focusChatMsg{} }
		}
	}

	return s, nil
}

// View renders the sidebar.
func (s Sidebar) View() string {
	var sb strings.Builder

	// Agents section.
	sb.WriteString(sidebarTitle.Render("Chat"))
	sb.WriteString("\n")
	for i, a := range s.agents {
		cursor := "  "
		style := sidebarItem
		if i == s.cursor && s.focused {
			cursor = "> "
			style = sidebarItemActive
		}
		status := ""
		if a.Status == "thinking" {
			status = " *"
		}
		sb.WriteString(style.Render(fmt.Sprintf("%s%s%s", cursor, a.Name, status)))
		sb.WriteString("\n")
	}

	// Conversation history section.
	sb.WriteString("\n")
	sb.WriteString(sidebarSection.Render("History"))
	sb.WriteString("\n")
	if len(s.conversations) == 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render("  (none)"))
		sb.WriteString("\n")
	}
	for _, c := range s.conversations {
		preview := c.Preview
		if len(preview) > 16 {
			preview = preview[:16] + ".."
		}
		if preview == "" {
			preview = "(empty)"
		}
		line := fmt.Sprintf("  %s %s", c.StartedAt, preview)
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render(line))
		sb.WriteString("\n")
	}

	// Activity section.
	sb.WriteString("\n")
	sb.WriteString(sidebarSection.Render("Activity"))
	sb.WriteString("\n")
	if len(s.activity) == 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render("  (none)"))
		sb.WriteString("\n")
	}
	for _, a := range s.activity {
		line := fmt.Sprintf("  %s>%s", a.From, a.To)
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render(line))
		sb.WriteString("\n")
	}

	// Tasks section.
	sb.WriteString("\n")
	sb.WriteString(sidebarSection.Render("Tasks"))
	sb.WriteString("\n")
	if s.tasks.Open+s.tasks.InProgress+s.tasks.Done == 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render("  (none)"))
	} else {
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render(
			fmt.Sprintf("  %d open  %d done", s.tasks.Open+s.tasks.InProgress, s.tasks.Done)))
	}
	sb.WriteString("\n")

	return sidebarStyle.Height(s.height).Render(sb.String())
}

// Internal messages.
type selectAgentMsg struct{ Agent string }
type focusChatMsg struct{}
