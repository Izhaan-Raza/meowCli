package tui

import (
	"fmt"
	"log" // <-- Import log
	"meowCli/gemini"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour" // <-- Import glamour
	"github.com/charmbracelet/lipgloss"
)

type (
	geminiResponseMsg string
	errorMsg          error
)

var (
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("35")).Bold(true)
	appNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("94")).Bold(true)
)

type model struct {
	textInput textinput.Model
	viewport  viewport.Model
	spinner   spinner.Model
	isLoading bool
}

// Helper function to render markdown with glamour
func renderMarkdown(content string) (string, error) {
	// Create a new glamour renderer with a dark style
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100), // Adjust word wrap to your liking
	)
	if err != nil {
		return "", err
	}
	// Render the markdown content
	out, err := renderer.Render(content)
	if err != nil {
		return "", err
	}
	return out, nil
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask me ......"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(80, 20)

	// Render the initial message with markdown
	initialContent, _ := renderMarkdown("I'm ready to help! Ask me to write some code.")
	vp.SetContent(initialContent)

	return model{
		textInput: ti,
		spinner:   s,
		viewport:  vp,
		isLoading: false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func getGeminiResponse(prompt string) tea.Cmd {
	return func() tea.Msg {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return errorMsg(fmt.Errorf("key not found"))
		}
		resp, err := gemini.GenerateContent(apiKey, prompt)
		if err != nil {
			return errorMsg(err)
		}
		return geminiResponseMsg(resp)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 8
		m.textInput.Width = msg.Width - 4

	case tea.KeyMsg:
		if !m.textInput.Focused() || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.isLoading {
				return m, nil
			}
			m.isLoading = true
			m.viewport.SetContent(fmt.Sprintf("Let me check Sir......\n\x1b[90m%s\x1b[0m", "ruk na lode"))
			prompt := m.textInput.Value()
			m.textInput.Reset()
			cmds = append(cmds, getGeminiResponse(prompt))
			return m, tea.Batch(cmds...)
		}

	case spinner.TickMsg:
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	// CHANGED THIS PART
	case geminiResponseMsg:
		m.isLoading = false
		// Render the response as markdown before setting it
		renderedContent, err := renderMarkdown(string(msg))
		if err != nil {
			log.Printf("Error rendering markdown: %v", err)
			m.viewport.SetContent(fmt.Sprintf("Error rendering response: %s", err))
		} else {
			m.viewport.SetContent(renderedContent)
		}

	case errorMsg:
		m.isLoading = false
		m.viewport.SetContent(fmt.Sprintf("Error: %s", msg))
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Welcome to ") + appNameStyle.Render("MeowCLi🐈🚀") + "\n\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n\n")

	if m.isLoading {
		b.WriteString(m.spinner.View() + " ")
	}

	b.WriteString(m.textInput.View())
	b.WriteString("\n\n\t\t press esc or CTRL+C to exit")

	return b.String()
}
