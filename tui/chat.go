package model

// A simple program demonstrating the text area component from the Bubbles
// component library.

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var (
// 		tiCmd tea.Cmd
// 		vpCmd tea.Cmd
// 	)
//
// 	m.textarea, tiCmd = m.textarea.Update(msg)
// 	m.viewport, vpCmd = m.viewport.Update(msg)
//
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		m.viewport.Width = msg.Width
// 		m.textarea.SetWidth(msg.Width)
// 		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
//
// 		if len(m.messages) > 0 {
// 			// Wrap content before setting it.
// 			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
// 		}
// 		m.viewport.GotoBottom()
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			fmt.Println(m.textarea.Value())
// 			return m, tea.Quit
// 		case tea.KeyEnter:
// 			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
// 			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
// 			m.textarea.Reset()
// 			m.viewport.GotoBottom()
// 		}
//
// 	// We handle errors just like any other message
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
// 	}
//
// 	return m, tea.Batch(tiCmd, vpCmd)
// }

// func (m model) View() string {
// 	return fmt.Sprintf(
// 		"%s%s%s",
// 		m.viewport.View(),
// 		gap,
// 		m.textarea.View(),
// 	)
// }
