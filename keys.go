package walk

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	forceQuit key.Binding
	quit      key.Binding
	quitQ     key.Binding
	open      key.Binding
	back      key.Binding
	up        key.Binding
	down      key.Binding
	left      key.Binding
	right     key.Binding
	top       key.Binding
	bottom    key.Binding
	leftmost  key.Binding
	rightmost key.Binding
	pageUp    key.Binding
	pageDown  key.Binding
	home      key.Binding
	end       key.Binding
	vimUp     key.Binding
	vimDown   key.Binding
	vimLeft   key.Binding
	vimRight  key.Binding
	vimTop    key.Binding
	vimBottom key.Binding
	search    key.Binding
	preview   key.Binding
	delete    key.Binding
	undo      key.Binding
	yank      key.Binding
}

func newKeyMap() *keyMap { return new(keyMap).init() }

func (k *keyMap) init() *keyMap {
	k.forceQuit = key.NewBinding(key.WithKeys("ctrl+c"))
	k.quit = key.NewBinding(key.WithKeys("esc"))
	k.quitQ = key.NewBinding(key.WithKeys("q"))
	k.open = key.NewBinding(key.WithKeys("enter"))
	k.back = key.NewBinding(key.WithKeys("backspace"))
	k.up = key.NewBinding(key.WithKeys("up"))
	k.down = key.NewBinding(key.WithKeys("down"))
	k.left = key.NewBinding(key.WithKeys("left"))
	k.right = key.NewBinding(key.WithKeys("right"))
	k.top = key.NewBinding(key.WithKeys("shift+up"))
	k.bottom = key.NewBinding(key.WithKeys("shift+down"))
	k.leftmost = key.NewBinding(key.WithKeys("shift+left"))
	k.rightmost = key.NewBinding(key.WithKeys("shift+right"))
	k.pageUp = key.NewBinding(key.WithKeys("pgup"))
	k.pageDown = key.NewBinding(key.WithKeys("pgdown"))
	k.home = key.NewBinding(key.WithKeys("home"))
	k.end = key.NewBinding(key.WithKeys("end"))
	k.vimUp = key.NewBinding(key.WithKeys("k"))
	k.vimDown = key.NewBinding(key.WithKeys("j"))
	k.vimLeft = key.NewBinding(key.WithKeys("h"))
	k.vimRight = key.NewBinding(key.WithKeys("l"))
	k.vimTop = key.NewBinding(key.WithKeys("g"))
	k.vimBottom = key.NewBinding(key.WithKeys("G"))
	k.search = key.NewBinding(key.WithKeys("/"))
	k.preview = key.NewBinding(key.WithKeys(" "))
	k.delete = key.NewBinding(key.WithKeys("d"))
	k.undo = key.NewBinding(key.WithKeys("u"))
	k.yank = key.NewBinding(key.WithKeys("y"))
	return k
}
