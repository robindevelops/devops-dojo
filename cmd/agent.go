package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/sensei"
)

type aiResponseMsg struct {
	text string
	err  error
}

func fetchAIResponse(prompt string) tea.Cmd {
	return func() tea.Msg {
		resp, err := sensei.AskAgent(prompt)
		return aiResponseMsg{text: resp, err: err}
	}
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Talk to Dojo AI Agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}

// Styling
var (
	rootTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("27")). // blue
			Padding(0, 1).
			Bold(true)

	userTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("239")). // dark gray
			Padding(0, 1).
			Bold(true)

	senderStyle = lipgloss.NewStyle().MarginBottom(1)

	welcomeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
)

type keyMap struct {
	Quit     key.Binding
	Switch   key.Binding
	Tab      key.Binding
	PrevNext key.Binding
	Commands key.Binding
	Help     key.Binding
	Newline  key.Binding
	EditVi   key.Binding
	History  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit, k.Switch, k.Tab, k.PrevNext, k.Commands, k.Help, k.Newline, k.EditVi, k.History,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Switch, k.Tab, k.PrevNext, k.Commands},
		{k.Help, k.Newline, k.EditVi, k.History},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("Ctrl+c", "quit"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "switch focus"),
	),
	Tab: key.NewBinding(
		key.WithKeys("ctrl+t", "ctrl+w"),
		key.WithHelp("Ctrl+t/w", "new/close tab"),
	),
	PrevNext: key.NewBinding(
		key.WithKeys("ctrl+p", "ctrl+n"),
		key.WithHelp("Ctrl+p/n", "prev/next tab"),
	),
	Commands: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("Ctrl+k", "commands"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("Ctrl+h", "help"),
	),
	Newline: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("Ctrl+j", "newline"),
	),
	EditVi: key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("Ctrl+g", "edit in Vi"),
	),
	History: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("Ctrl+r", "history search"),
	),
}

type errMsg error

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	help        help.Model
	senderStyle lipgloss.Style
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "type your message here.."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 1000

	ta.SetWidth(80)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	welcomeMsgText := "Welcome to DevOps Dojo AI!\nType your questions about incidents, commands, or DevOps concepts.\n--------------------------------------------------------------\n\n"
	welcomeMsg := welcomeStyle.Render(welcomeMsgText)
	vp.SetContent(
		welcomeMsg +
			senderStyle.Render(rootTagStyle.Render("root")) + "\n" +
			"Ready to start fresh. What do you need help with?\n",
	)

	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	h.Styles.FullSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)

	return model{
		textarea:    ta,
		help:        h,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().MarginBottom(1),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.textarea.Value()
			// KeyEnter in textarea actually adds a newline, we want to submit on enter
			// BUT bubbles/textarea is multi-line by default. 
			// Let's trim trailing newline
			val = strings.TrimSpace(val)

			if val == "clear" {
				m.messages = []string{}
				welcomeMsgText := "Welcome to DevOps Dojo AI!\nType your questions about incidents, commands, or DevOps concepts.\n--------------------------------------------------------------\n\n"
				welcomeMsg := welcomeStyle.Render(welcomeMsgText)
				m.viewport.SetContent(
					welcomeMsg +
						senderStyle.Render(rootTagStyle.Render("root")) + "\n" +
						"Ready to start fresh. What do you need help with?\n",
				)
				m.textarea.Reset()
				return m, nil
			}

			if val != "" {
				userMsg := senderStyle.Render(userTagStyle.Render("you")) + "\n" + val + "\n"
				m.messages = append(m.messages, userMsg)

				// Append a loading/thinking message
				agentMsg := senderStyle.Render(rootTagStyle.Render("root")) + "\n" + "Thinking...\n"
				m.messages = append(m.messages, agentMsg)

				m.viewport.SetContent(m.buildContent())
				m.textarea.Reset()
				m.viewport.GotoBottom()

				return m, fetchAIResponse(val)
			}
		}

	case aiResponseMsg:
		if msg.err != nil {
			m.messages[len(m.messages)-1] = senderStyle.Render(rootTagStyle.Render("root")) + "\n" + "Error: " + msg.err.Error() + "\n"
		} else {
			m.messages[len(m.messages)-1] = senderStyle.Render(rootTagStyle.Render("root")) + "\n" + msg.text + "\n"
		}
		m.viewport.SetContent(m.buildContent())
		m.viewport.GotoBottom()
		return m, nil

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - 3 // -3 for margin and help bar
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
		m.help.View(keys),
	)
}

func (m *model) buildContent() string {
	var content string
	welcomeMsgText := "Welcome to DevOps Dojo AI!\nType your questions about incidents, commands, or DevOps concepts.\n--------------------------------------------------------------\n\n"
	welcomeMsg := welcomeStyle.Render(welcomeMsgText)
	content += welcomeMsg + senderStyle.Render(rootTagStyle.Render("root")) + "\n" + "Ready to start fresh. What do you need help with?\n\n"
	for _, msg := range m.messages {
		content += msg + "\n"
	}
	return content
}
