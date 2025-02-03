package key

import "github.com/charmbracelet/bubbles/key"

type Map struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding

	PageUp   key.Binding
	PageDown key.Binding

	Help key.Binding
	Quit key.Binding
}

var Default = &Map{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),

	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("<Tab>", "tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("<S-Tab>", "shift tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("<Enter>/<Space>", "enter"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+f"),
		key.WithHelp("<C-f>", "preview page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+b"),
		key.WithHelp("<C-b>", "preview page down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
