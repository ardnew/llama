package walk

import "github.com/charmbracelet/bubbles/key"

// keyMap defines all key bindings handled by walk.
//
// The type itself is unexported to prevent accidentally using a zero-valued
// keyMap with uninitialized bindings. Users can construct and modify a fully
// initialized keyMap returned from NewKeyMap. 
type keyMap struct {
	ForceQuit key.Binding
	Quit      key.Binding
	QuitQ     key.Binding
	Submit    key.Binding
	Next      key.Binding
	Prev      key.Binding
	Open      key.Binding
	Back      key.Binding
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Top       key.Binding
	Bottom    key.Binding
	Leftmost  key.Binding
	Rightmost key.Binding
	PageUp    key.Binding
	PageDown  key.Binding
	Home      key.Binding
	End       key.Binding
	VimUp     key.Binding
	VimDown   key.Binding
	VimLeft   key.Binding
	VimRight  key.Binding
	VimTop    key.Binding
	VimBottom key.Binding
	Search    key.Binding
	Preview   key.Binding
	Delete    key.Binding
	Undo      key.Binding
	Yank      key.Binding
}

func NewKeyMap() *keyMap { return new(keyMap).Default() }

func (k *keyMap) Default() *keyMap {
	if k == nil {
		k = new(keyMap)
	}
	k.ForceQuit = key.NewBinding(key.WithKeys("ctrl+c"))
	k.Quit = key.NewBinding(key.WithKeys("esc"))
	k.QuitQ = key.NewBinding(key.WithKeys("q"))
	k.Submit = key.NewBinding(key.WithKeys("enter"))
	k.Next = key.NewBinding(key.WithKeys("tab"))
	k.Prev = key.NewBinding(key.WithKeys("shift+tab"))
	k.Open = key.NewBinding(key.WithKeys("]", " ", "o", "n"))
	k.Back = key.NewBinding(key.WithKeys("[", "backspace", "b", "p"))
	k.Up = key.NewBinding(key.WithKeys("up"))
	k.Down = key.NewBinding(key.WithKeys("down"))
	k.Left = key.NewBinding(key.WithKeys("left"))
	k.Right = key.NewBinding(key.WithKeys("right"))
	k.Top = key.NewBinding(key.WithKeys("shift+up"))
	k.Bottom = key.NewBinding(key.WithKeys("shift+down"))
	k.Leftmost = key.NewBinding(key.WithKeys("shift+left"))
	k.Rightmost = key.NewBinding(key.WithKeys("shift+right"))
	k.PageUp = key.NewBinding(key.WithKeys("pgup"))
	k.PageDown = key.NewBinding(key.WithKeys("pgdown"))
	k.Home = key.NewBinding(key.WithKeys("home"))
	k.End = key.NewBinding(key.WithKeys("end"))
	k.VimUp = key.NewBinding(key.WithKeys("k"))
	k.VimDown = key.NewBinding(key.WithKeys("j"))
	k.VimLeft = key.NewBinding(key.WithKeys("h"))
	k.VimRight = key.NewBinding(key.WithKeys("l"))
	k.VimTop = key.NewBinding(key.WithKeys("g"))
	k.VimBottom = key.NewBinding(key.WithKeys("G"))
	k.Search = key.NewBinding(key.WithKeys("/"))
	k.Preview = key.NewBinding(key.WithKeys("`"))
	k.Delete = key.NewBinding(key.WithKeys("delete", "d"))
	k.Undo = key.NewBinding(key.WithKeys("u", "z"))
	k.Yank = key.NewBinding(key.WithKeys("y"))
	return k
}
