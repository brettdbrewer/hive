package tui

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// paneFocus tracks which pane has keyboard focus.
type paneFocus int

const (
	paneChat paneFocus = iota
	paneSidebar
)

// App is the root bubbletea model for the hive TUI.
type App struct {
	ctx     context.Context
	backend Backend
	focus   paneFocus
	sidebar Sidebar
	chat    ChatView
	info    SystemInfo
	width   int
	height  int
}

// New creates a new TUI application.
func New(ctx context.Context, backend Backend) App {
	chat := NewChatView()
	chat.Focus()

	sidebar := NewSidebar()

	return App{
		ctx:     ctx,
		backend: backend,
		focus:   paneChat,
		sidebar: sidebar,
		chat:    chat,
	}
}

// TUILogPath is where the TUI writes debug logs.
// Set by the caller before Run() to ensure a writable absolute path.
var TUILogPath string

// appendToLog writes a line to the TUI debug log.
func appendToLog(line string) {
	if TUILogPath == "" {
		return
	}
	f, err := os.OpenFile(TUILogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	f.WriteString(line)
	f.Close()
}

// Run starts the TUI.
func Run(ctx context.Context, backend Backend) error {
	// Clear previous log.
	os.Remove(TUILogPath)
	appendToLog("=== TUI started ===\n")

	app := New(ctx, backend)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Init returns initial commands.
func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.chat.Init(),
		a.loadInfo(),
		a.loadAgents(),
		a.loadHistory(),
		a.loadConversations(),
	)
}

// Update handles all messages.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateLayout()
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "ctrl+l":
			path := conversationLogPath()
			if err := a.dumpConversation(); err != nil {
				a.chat.AddMessage("error", fmt.Sprintf("Could not save: %v", err))
			} else {
				a.chat.AddMessage("system", fmt.Sprintf("Conversation saved to %s", path))
			}
			return a, nil
		case "tab":
			a.cycleFocus()
			return a, nil
		}

	// Chat sent a message — relay to backend.
	case sendMsg:
		a.chat.SetThinking(true)
		return a, a.chatCmd(msg.Text)

	// Backend responded.
	case chatResponseMsg:
		a.chat.SetThinking(false)
		a.chat.AddMessage("mind", msg.Content)
		return a, nil

	case chatErrorMsg:
		a.chat.SetThinking(false)
		a.chat.AddMessage("error", msg.Err.Error())
		return a, nil

	// Sidebar wants to focus chat.
	case focusChatMsg:
		a.setFocus(paneChat)
		return a, nil

	case blurChatMsg:
		a.setFocus(paneSidebar)
		return a, nil

	case selectAgentMsg:
		// For now, all agents route to mind.
		a.chat.AddMessage("system", fmt.Sprintf("Switched to %s", msg.Agent))
		a.setFocus(paneChat)
		return a, nil

	// Info loaded.
	case infoMsg:
		a.info = SystemInfo(msg)
		return a, nil

	case agentsMsg:
		a.sidebar.SetAgents([]AgentInfo(msg))
		return a, nil

	// History loaded from previous session.
	case historyMsg:
		for _, m := range msg {
			a.chat.AddMessage(m.Role, m.Content)
		}
		return a, nil

	// Conversation list loaded.
	case conversationsMsg:
		a.sidebar.SetConversations(msg)
		return a, nil
	}

	// Route to focused component.
	var cmd tea.Cmd
	switch a.focus {
	case paneChat:
		a.chat, cmd = a.chat.Update(msg)
	case paneSidebar:
		a.sidebar, cmd = a.sidebar.Update(msg)
	}
	return a, cmd
}

// View renders the full TUI.
func (a App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	// Title bar.
	title := titleStyle.Width(a.width).Render(
		fmt.Sprintf(" hive mind   %s", a.modelLabel()))

	// Sidebar + chat.
	sidebar := a.sidebar.View()
	chat := a.chat.View()
	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chat)

	// Status bar.
	status := a.renderStatus()

	return lipgloss.JoinVertical(lipgloss.Left, title, content, status)
}

func (a *App) updateLayout() {
	contentHeight := a.height - 2 // title + status
	chatWidth := a.width - sidebarWidth - 2
	if chatWidth < 20 {
		chatWidth = 20
	}
	a.sidebar.SetSize(contentHeight)
	a.chat.SetSize(chatWidth, contentHeight)
}

func (a *App) cycleFocus() {
	switch a.focus {
	case paneChat:
		a.setFocus(paneSidebar)
	case paneSidebar:
		a.setFocus(paneChat)
	}
}

func (a *App) setFocus(f paneFocus) {
	a.focus = f
	switch f {
	case paneChat:
		a.chat.Focus()
		a.sidebar.Blur()
	case paneSidebar:
		a.chat.Blur()
		a.sidebar.Focus()
	}
}

func (a App) modelLabel() string {
	if a.info.Model != "" {
		return statusAccent.Render(a.info.Model)
	}
	return statusAccent.Render("Opus")
}

func (a App) renderStatus() string {
	var parts []string
	if a.info.MCPConnected {
		parts = append(parts, "MCP: connected")
	} else {
		parts = append(parts, "MCP: disconnected")
	}
	if a.info.EventCount > 0 {
		parts = append(parts, fmt.Sprintf("Events: %d", a.info.EventCount))
	}
	if a.info.TaskCount > 0 {
		parts = append(parts, fmt.Sprintf("Tasks: %d", a.info.TaskCount))
	}
	if a.info.CostUSD > 0 {
		parts = append(parts, fmt.Sprintf("$%.2f", a.info.CostUSD))
	}

	status := strings.Join(parts, "  |  ")
	return statusStyle.Width(a.width).Render(status)
}

// conversationLogPath returns the path for the conversation dump file.
func conversationLogPath() string {
	return fmt.Sprintf("%s%chive-conversation.log", os.TempDir(), os.PathSeparator)
}

// dumpConversation writes the full chat to a log file for copy/paste.
func (a *App) dumpConversation() error {
	var sb strings.Builder
	for _, m := range a.chat.messages {
		label := m.Role
		if label == "human" {
			label = "Matt"
		} else if label == "mind" {
			label = "Mind"
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n\n", label, m.Content))
	}
	return os.WriteFile(conversationLogPath(), []byte(sb.String()), 0644)
}

// Commands for async backend calls.

type chatResponseMsg struct{ Content string }
type chatErrorMsg struct{ Err error }
type infoMsg SystemInfo
type agentsMsg []AgentInfo
type historyMsg []ChatMessage
type conversationsMsg []ConversationInfo

func (a *App) chatCmd(text string) tea.Cmd {
	ctx := a.ctx
	backend := a.backend
	return func() tea.Msg {
		resp, err := backend.Chat(ctx, "mind", text)
		if err != nil {
			return chatErrorMsg{Err: err}
		}
		return chatResponseMsg{Content: resp}
	}
}

func (a *App) loadInfo() tea.Cmd {
	ctx := a.ctx
	backend := a.backend
	return func() tea.Msg {
		info, err := backend.Info(ctx)
		if err != nil {
			return infoMsg{}
		}
		return infoMsg(info)
	}
}

func (a *App) loadAgents() tea.Cmd {
	ctx := a.ctx
	backend := a.backend
	return func() tea.Msg {
		agents, err := backend.ListAgents(ctx)
		if err != nil {
			return agentsMsg(nil)
		}
		return agentsMsg(agents)
	}
}

func (a *App) loadHistory() tea.Cmd {
	ctx := a.ctx
	backend := a.backend
	return func() tea.Msg {
		msgs, err := backend.History(ctx)
		if err != nil || len(msgs) == 0 {
			return historyMsg(nil)
		}
		return historyMsg(msgs)
	}
}

func (a *App) loadConversations() tea.Cmd {
	ctx := a.ctx
	backend := a.backend
	return func() tea.Msg {
		convs, err := backend.ListConversations(ctx, 20)
		if err != nil {
			return conversationsMsg(nil)
		}
		return conversationsMsg(convs)
	}
}
