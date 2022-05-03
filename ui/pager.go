package ui

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	lip "github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = true

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
	return fmt.Sprintf("%3.f%%", p*100)
}

func larger(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Preview(url, title string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	r, err := glamour.RenderBytes(body, "dark")
	if err != nil {
		return err
	}
	pager := tea.NewProgram(
		&model{
			content: string(r),
			title:   title,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	return pager.Start()
}
