package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChatView is the main chat pane: message viewport + text input.
type ChatView struct {
	viewport viewport.Model
	input    textarea.Model
	spinner  spinner.Model
	messages []ChatMessage
	width    int
	height   int
	focused  bool
	thinking bool
}

// NewChatView creates a new chat view.
func NewChatView() ChatView {
	ta := textarea.New()
	ta.Placeholder = "Type a message... (Enter to send)"
	ta.CharLimit = 4096
	ta.SetHeight(inputHeight)
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(accent)

	vp := viewport.New(80, 20)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 3

	return ChatView{
		viewport: vp,
		input:    ta,
		spinner:  sp,
	}
}

// SetSize updates the chat view dimensions.
func (c *ChatView) SetSize(w, h int) {
	c.width = w
	c.height = h
	c.input.SetWidth(w - 2)
	c.viewport.Width = w - 2
	c.viewport.Height = h - inputHeight - 3 // leave room for input + border
	c.refreshViewport()
}

// Focus gives keyboard focus to the chat input.
func (c *ChatView) Focus() {
	c.focused = true
	c.input.Focus()
}

// Blur removes keyboard focus.
func (c *ChatView) Blur() {
	c.focused = false
	c.input.Blur()
}

// AddMessage appends a message and scrolls to bottom.
// Also appends to a log file for debugging.
func (c *ChatView) AddMessage(role, content string) {
	c.messages = append(c.messages, ChatMessage{Role: role, Content: content})
	c.refreshViewport()
	c.viewport.GotoBottom()

	// Append to log file for debugging (TUI makes copy/paste hard).
	appendToLog(fmt.Sprintf("[%s] %s\n", role, content))
}

// SetThinking toggles the thinking indicator.
func (c *ChatView) SetThinking(v bool) {
	c.thinking = v
}

// Init returns initial commands.
func (c ChatView) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, c.spinner.Tick)
}

// Update handles messages.
func (c ChatView) Update(msg tea.Msg) (ChatView, tea.Cmd) {
	var cmds []tea.Cmd

	if c.thinking {
		var cmd tea.Cmd
		c.spinner, cmd = c.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !c.focused {
			return c, nil
		}
		switch msg.Type {
		case tea.KeyEnter:
			text := strings.TrimSpace(c.input.Value())
			if text != "" && !c.thinking {
				c.AddMessage("human", text)
				c.input.Reset()
				return c, func() tea.Msg { return sendMsg{Text: text} }
			}
			return c, nil
		case tea.KeyEsc:
			return c, func() tea.Msg { return blurChatMsg{} }

		// Scroll keys → viewport only (don't let textarea consume them).
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			c.viewport, cmd = c.viewport.Update(msg)
			cmds = append(cmds, cmd)
			return c, tea.Batch(cmds...)

		// Ctrl+Home / Ctrl+End for top/bottom.
		case tea.KeyHome:
			c.viewport.GotoTop()
			return c, tea.Batch(cmds...)
		case tea.KeyEnd:
			c.viewport.GotoBottom()
			return c, tea.Batch(cmds...)

		default:
			// All other keys → textarea (typing).
			if !c.thinking {
				var cmd tea.Cmd
				c.input, cmd = c.input.Update(msg)
				cmds = append(cmds, cmd)
			}
			return c, tea.Batch(cmds...)
		}

	case tea.MouseMsg:
		// Mouse events → viewport (scroll wheel).
		var cmd tea.Cmd
		c.viewport, cmd = c.viewport.Update(msg)
		cmds = append(cmds, cmd)
		return c, tea.Batch(cmds...)
	}

	// Non-key messages (window resize, etc.) go to both.
	if c.focused && !c.thinking {
		var cmd tea.Cmd
		c.input, cmd = c.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	c.viewport, cmd = c.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

// View renders the chat pane.
func (c ChatView) View() string {
	var sb strings.Builder
	sb.WriteString(c.viewport.View())
	sb.WriteString("\n")

	// Scroll indicator between viewport and input.
	scrollInfo := c.scrollIndicator()
	if scrollInfo != "" {
		sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Render(scrollInfo))
		sb.WriteString("\n")
	}

	if c.thinking {
		thinkLine := fmt.Sprintf(" %s thinking...", c.spinner.View())
		sb.WriteString(inputBorder.Width(c.width - 2).Render(thinkLine))
	} else {
		sb.WriteString(inputBorder.Width(c.width - 2).Render(c.input.View()))
	}

	return chatStyle.Width(c.width).Height(c.height).Render(sb.String())
}

// scrollIndicator shows position when not at the bottom.
func (c ChatView) scrollIndicator() string {
	if c.viewport.TotalLineCount() <= c.viewport.Height {
		return "" // all content visible, no indicator needed
	}
	pct := c.viewport.ScrollPercent()
	if pct >= 1.0 {
		return "" // at bottom, no indicator
	}
	return fmt.Sprintf(" --- scroll: %.0f%% (PgUp/PgDown to scroll, End to jump to latest) ---", pct*100)
}

// refreshViewport rebuilds the viewport content from messages.
func (c *ChatView) refreshViewport() {
	w := c.viewport.Width
	if w < 10 {
		w = 80
	}

	var sb strings.Builder
	for _, m := range c.messages {
		// Word-wrap content to viewport width.
		wrapped := lipgloss.NewStyle().Width(w).Render(m.Content)
		switch m.Role {
		case "human":
			sb.WriteString(humanMsg.Render("Matt") + ": " + wrapped + "\n\n")
		case "mind":
			sb.WriteString(mindMsg.Render("Mind") + ": " + wrapped + "\n\n")
		case "error":
			sb.WriteString(errorMsg.Render("Error") + ": " + wrapped + "\n\n")
		case "system":
			sb.WriteString(lipgloss.NewStyle().Foreground(dimText).Width(w).Render(m.Content) + "\n\n")
		}
	}
	c.viewport.SetContent(sb.String())
}

// Internal messages.
type sendMsg struct{ Text string }
type blurChatMsg struct{}
