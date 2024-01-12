package walk

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type Field struct {
	*Model

	value *FilePath
	key   string

	// error handling
	validate func(FilePath) error
	err      error

	// state
	isFocused  bool
	isFiltered bool
	filter     textinput.Model

	// customization
	heading string
	caption string

	// options
	accessible bool
	showAll    bool
	theme      *huh.Theme
}

// Theme(theme *huh.Theme) huh.Field
// Accessible(accessible bool) huh.Field
// KeyMap(keys *huh.KeyMap) huh.Field
// Width(width int) huh.Field

func Value(value string) Option[*Field] {
	return func(f *Field) *Field { return f.WithValue(value) }
}

func Key(key string) Option[*Field] {
	return func(f *Field) *Field { return f.WithKey(key) }
}

func Heading(heading string) Option[*Field] {
	return func(f *Field) *Field { return f.WithHeading(heading) }
}

func Caption(caption string) Option[*Field] {
	return func(f *Field) *Field { return f.WithCaption(caption) }
}

func Validate(validate func(FilePath) error) Option[*Field] {
	return func(f *Field) *Field { return f.WithValidate(validate) }
}

// Field returns a new Field of the receiver Model.
//
// Field implements both its epynomous interface and the Model interface used in
// the Bubble Tea framework (module packages "huh" & "bubbletea", respectively).
func (m *Model) Field(options ...Option[*Field]) *Field {
	filter := textinput.New()
	filter.Prompt = "/"

	return &Field{
		Model:    m,
		value:    new(FilePath),
		validate: func(FilePath) error { return nil },
		filter:   filter,
	}
}

// Init initializes the Field.
func (f *Field) Init() tea.Cmd {
	return f.Model.Init()
}

func (f *Field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := f.Model.Update(msg)
	return f, cmd
}

func (f *Field) View() string {
	return f.Model.View()
}

// Blur blurs the Field.
func (f *Field) Blur() tea.Cmd {
	f.isFocused = false
	f.err = f.validate(*f.value)
	return nil
}

// Focus focuses the Field.
func (f *Field) Focus() tea.Cmd {
	f.isFocused = true
	return nil
}

// Error returns the error of the Field.
func (f *Field) Error() error {
	return f.err
}

// Run runs the Field.
func (f *Field) Run() error {
	if f.accessible {
		return f.runAccessible()
	}
	return newRunError(huh.Run(f))
}

// KeyBinds returns the keybindings for the Field.
func (f *Field) KeyBinds() []key.Binding {
	return []key.Binding{} // f.keys.bindings()
}

// With returns the receiver with the given options applied.
func (f *Field) With(options ...Option[*Field]) *Field {
	for _, option := range options {
		f = option(f)
	}
	return f
}

// WithTheme sets the theme of the Field.
func (f *Field) WithTheme(theme *huh.Theme) huh.Field {
	f.theme = theme
	f.filter.Cursor.Style = f.theme.Focused.TextInput.Cursor
	f.filter.PromptStyle = f.theme.Focused.TextInput.Prompt
	return f
}

// WithAccessible sets the accessible mode of the Field.
func (f *Field) WithAccessible(accessible bool) huh.Field {
	f.accessible = accessible
	return f
}

// WithKeyMap sets the keymap on a Field.
func (f *Field) WithKeyMap(keys *huh.KeyMap) huh.Field {
	// TBD
	return f
}

// WithWidth sets the width of the Field.
func (f *Field) WithWidth(width int) huh.Field {
	f.width = width
	return f
}

// GetKey returns the key of the field.
func (f *Field) GetKey() string {
	return f.key
}

// GetValue returns the value of the field.
func (f *Field) GetValue() any {
	return f.value.path()
}

// Value sets the value of the Field.
func (f *Field) WithValue(value string) *Field {
	f.value = f.value.init(value)
	return f
}

// Key sets the key of the Field which can be used to retrieve the value
// after submission.
func (f *Field) WithKey(key string) *Field {
	f.key = key
	return f
}

// Heading sets the heading of the Field.
func (f *Field) WithHeading(heading string) *Field {
	f.heading = heading
	return f
}

// Caption sets the caption of the Field.
func (f *Field) WithCaption(caption string) *Field {
	f.caption = caption
	return f
}

// Validate sets the validation function of the Field.
func (f *Field) WithValidate(validate func(FilePath) error) *Field {
	f.validate = validate
	return f
}

func (f *Field) runAccessible() error {
	var sb strings.Builder
	sb.WriteString(f.theme.Focused.Title.Render(f.heading) + "\n")

	// for i, option := range t.option {
	// 	sb.WriteString(fmt.Sprintf("%d. %s", i+1, option.Key))
	// 	sb.WriteString("\n")
	// }
	//
	// fmt.Println(t.theme.Blurred.Base.Render(sb.String()))
	//
	// for {
	//	choice := accessibility.PromptInt("Choose: ", 1, len(t.options))
	//	option := t.options[choice-1]
	//	if err := t.validate(option.Value); err != nil {
	//		fmt.Println(err.Error())
	//		continue
	//	}
	//	fmt.Println(t.theme.Focused.SelectedOption.Render("Chose: " + option.Key + "\n"))
	//	*t.value = option.Value
	//	break
	//}

	return nil
}

func (f *Field) setIsFiltered(isFiltered bool) {
	f.isFiltered = isFiltered
}

func (f *Field) filterFunc(option string) bool {
	// XXX: remove diacritics or allow customization of filter function.
	return strings.Contains(
		strings.ToLower(option),
		strings.ToLower(f.filter.Value()),
	)
}
