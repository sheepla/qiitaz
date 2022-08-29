package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	lip "github.com/charmbracelet/lipgloss"
	"github.com/sheepla/qiitaz/client"
)

const (
	useHighPerformanceRenderer = true
	glamourTheme               = "dark"
)

// nolint:gochecknoglobals
var (
	titleStyle = func() lip.Style {
		b := lip.NormalBorder()
		b.Right = "├"

		return lip.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lip.Style {
		b := lip.NormalBorder()
		b.Left = "┤"

		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	title    string
}

func (m *model) Init() tea.Cmd {
	return nil
}

// nolint:ireturn
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

		if msg.String() == "g" {
			m.viewport.GotoTop()
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

		if msg.String() == "G" {
			m.viewport.GotoBottom()
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case tea.WindowSizeMsg:
		headerHeight := lip.Height(m.headerView())
		footerHeight := lip.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		cmds = append(cmds, viewport.Sync(m.viewport))
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m *model) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", larger(0, m.viewport.Width-lip.Width(title)))

	return lip.JoinHorizontal(lip.Center, title, line)
}

func (m *model) footerView() string {
	info := infoStyle.Render(scrollPercent(m.viewport.ScrollPercent()))
	line := strings.Repeat("─", larger(0, m.viewport.Width-lip.Width(info)))

	return lip.JoinHorizontal(lip.Center, line, info)
}

func scrollPercent(p float64) string {
	// nolint:gomnd
	return fmt.Sprintf("%3.f%%", p*100)
}

func larger(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func NewPagerProgram(path string, title string) (*tea.Program, error) {
	body, err := client.FetchArticle(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the article page: %w", err)
	}
	defer body.Close()

	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	content, err := glamour.RenderBytes(bytes, glamourTheme)
	if err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	pager := tea.NewProgram(
		// nolint:exhaustivestruct,exhaustruct
		&model{
			title:   title,
			content: string(content),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return pager, nil
}
