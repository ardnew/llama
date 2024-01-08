package walk

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	. "strings"
	"time"
	"unicode/utf8"

	"github.com/antonmedv/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

var version = "v2.1.0"

func Version() string { return version }

const separator = "    " // Separator between columns.

var (
	fileSeparator = string(filepath.Separator)
	showIcons     = false
)

type Styles struct {
	Warning, Preview, Cursor, Bar, Search, Danger lipgloss.Style
}

func DefaultStyle() *Styles {
	return &Styles{
		Warning: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).PaddingLeft(1).PaddingRight(1),
		Preview: lipgloss.NewStyle().PaddingLeft(2),
		Cursor:  lipgloss.NewStyle().Background(lipgloss.Color("#825DF2")).Foreground(lipgloss.Color("#FFFFFF")),
		Bar:     lipgloss.NewStyle().Background(lipgloss.Color("#5C5C5C")).Foreground(lipgloss.Color("#FFFFFF")),
		Search:  lipgloss.NewStyle().Background(lipgloss.Color("#499F1C")).Foreground(lipgloss.Color("#FFFFFF")),
		Danger:  lipgloss.NewStyle().Background(lipgloss.Color("#FF0000")).Foreground(lipgloss.Color("#FFFFFF")),
	}
}

type Model struct {
	path              string              // Current dir path we are looking at.
	files             []fs.DirEntry       // Files we are looking at.
	err               error               // Error while listing files.
	kb                *keyMap             // Key bindings.
	st                *Styles             // Rendering attributes.
	cmdline           []string            // Command line to open files.
	c, r              int                 // Selector position in columns and rows.
	columns, rows     int                 // Displayed amount of rows and columns.
	width, height     int                 // Terminal size.
	offset            int                 // Scroll position.
	positions         map[string]position // Map of cursor positions per path.
	search            string              // Type to select files with this value.
	searchMode        bool                // Whether type-to-select is active.
	searchId          int                 // Search id to indicate what search we are currently on.
	matchedIndexes    []int               // List of char found indexes.
	prevName          string              // Base name of previous directory before "up".
	findPrevName      bool                // On View(), set c&r to point to prevName.
	exitCode          int                 // Exit code.
	previewMode       bool                // Whether preview is active.
	previewContent    string              // Content of preview.
	deleteCurrentFile bool                // Whether to delete current file.
	toBeDeleted       []toDelete          // Map of files to be deleted.
	yankSuccess       bool                // Show yank info
}

type position struct {
	c, r   int
	offset int
}

type toDelete struct {
	path string
	at   time.Time
}

type (
	clearSearchMsg int
	toBeDeletedMsg int
)

type Option func(*Model) *Model

func Path(path string) Option {
	return func(m *Model) *Model { return m.WithPath(path) }
}

func (m *Model) WithPath(path string) *Model {
	m.path = path
	return m
}

func Size(width, height int) Option {
	return func(m *Model) *Model { return m.WithSize(width, height) }
}

func (m *Model) WithSize(width, height int) *Model {
	m.width = width
	m.height = height
	return m
}

func Icons() Option {
	return func(m *Model) *Model { return m.WithIcons() }
}

func (m *Model) WithIcons() *Model {
	showIcons = true
	parseIcons()
	return m
}

func Command(cmd string) Option {
	return func(m *Model) *Model { return m.WithCommand(cmd) }
}

func (m *Model) WithCommand(cmd string) *Model {
	m.cmdline = Fields(cmd)
	return m
}

func Style(styles *Styles) Option {
	return func(m *Model) *Model { return m.WithStyle(styles) }
}

func (m *Model) WithStyle(styles *Styles) *Model {
	m.st = styles
	return m
}

func New(options ...Option) *Model {
	m := (&Model{
		kb:        newKeyMap(),
		positions: make(map[string]position),
	}).With(options...)

	if m.st == nil {
		m.st = DefaultStyle()
	}
	return m
}

func (m *Model) With(options ...Option) *Model {
	for _, option := range options {
		m = option(m)
	}
	return m
}

func Kill(status int)  { os.Exit(status) }
func (m *Model) Exit() { Kill(m.exitCode) }

func (m *Model) Init() tea.Cmd {
	if m.path == "" {
		var err error
		m.path, err = os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr,
				m.st.Danger.Render("error: failed to get working directory"))
		}
	}
	m.list()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reset position history as c&r changes.
		m.positions = make(map[string]position)
		// Keep cursor at same place.
		fileName, ok := m.fileName()
		if ok {
			m.prevName = fileName
			m.findPrevName = true
		}
		// Also, m.c&r no longer point to the correct indexes.
		m.c = 0
		m.r = 0
		return m, nil

	case tea.KeyMsg:
		if m.searchMode {
			if key.Matches(msg, m.kb.search) {
				m.searchMode = false
				return m, nil
			} else if key.Matches(msg, m.kb.back) {
				if len(m.search) > 0 {
					m.search = m.search[:len(m.search)-1]
					return m, nil
				}
			} else if msg.Type == tea.KeyRunes {
				m.search += string(msg.Runes)
				names := make([]string, len(m.files))
				for i, fi := range m.files {
					names[i] = fi.Name()
				}
				matches := fuzzy.Find(m.search, names)
				if len(matches) > 0 {
					m.matchedIndexes = matches[0].MatchedIndexes
					index := matches[0].Index
					m.c = index / m.rows
					m.r = index % m.rows
				}
				m.updateOffset()
				m.saveCursorPosition()
				// Save search id to clear only current search after delay.
				// User may have already started typing next search.
				searchId := m.searchId
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return clearSearchMsg(searchId)
				})
			}
		}

		switch {
		case key.Matches(msg, m.kb.forceQuit):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible after prompt.
			m.exitCode = 2
			m.dontDoPendingDeletions()
			return m, tea.Quit

		case key.Matches(msg, m.kb.quit, m.kb.quitQ):
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible after prompt.
			fmt.Println(m.path)            // Write to cd.
			m.exitCode = 0
			m.performPendingDeletions()
			return m, tea.Quit

		case key.Matches(msg, m.kb.open):
			m.searchMode = false
			filePath, ok := m.filePath()
			if !ok {
				return m, nil
			}
			if fi := fileInfo(filePath); fi.IsDir() {
				// Enter subdirectory.
				m.path = filePath
				if p, ok := m.positions[m.path]; ok {
					m.c = p.c
					m.r = p.r
					m.offset = p.offset
				} else {
					m.c = 0
					m.r = 0
					m.offset = 0
				}
				m.list()
			} else {
				// Open file. This will block until complete.
				return m, m.openCommand()
			}

		case key.Matches(msg, m.kb.back):
			m.searchMode = false
			m.prevName = filepath.Base(m.path)
			m.path = filepath.Join(m.path, "..")
			if p, ok := m.positions[m.path]; ok {
				m.c = p.c
				m.r = p.r
				m.offset = p.offset
			} else {
				m.findPrevName = true
			}
			m.list()
			return m, nil

		case key.Matches(msg, m.kb.up):
			m.moveUp()

		case key.Matches(msg, m.kb.top, m.kb.pageUp, m.kb.vimTop):
			m.moveTop()

		case key.Matches(msg, m.kb.bottom, m.kb.pageDown, m.kb.vimBottom):
			m.moveBottom()

		case key.Matches(msg, m.kb.leftmost):
			m.moveLeftmost()

		case key.Matches(msg, m.kb.rightmost):
			m.moveRightmost()

		case key.Matches(msg, m.kb.home):
			m.moveStart()

		case key.Matches(msg, m.kb.end):
			m.moveEnd()

		case key.Matches(msg, m.kb.vimUp):
			if !m.searchMode {
				m.moveUp()
			}

		case key.Matches(msg, m.kb.down):
			m.moveDown()

		case key.Matches(msg, m.kb.vimDown):
			if !m.searchMode {
				m.moveDown()
			}

		case key.Matches(msg, m.kb.left):
			m.moveLeft()

		case key.Matches(msg, m.kb.vimLeft):
			if !m.searchMode {
				m.moveLeft()
			}

		case key.Matches(msg, m.kb.right):
			m.moveRight()

		case key.Matches(msg, m.kb.vimRight):
			if !m.searchMode {
				m.moveRight()
			}

		case key.Matches(msg, m.kb.search):
			m.searchMode = true
			m.searchId++
			m.search = ""

		case key.Matches(msg, m.kb.preview):
			m.previewMode = !m.previewMode
			// Reset position history as c&r changes.
			m.positions = make(map[string]position)
			// Keep cursor at same place.
			fileName, ok := m.fileName()
			if !ok {
				return m, nil
			}
			m.prevName = fileName
			m.findPrevName = true

			if m.previewMode {
				return m, tea.EnterAltScreen
			}
			m.previewContent = ""
			return m, tea.ExitAltScreen

		case key.Matches(msg, m.kb.delete):
			filePathToDelete, ok := m.filePath()
			if ok {
				if m.deleteCurrentFile {
					m.deleteCurrentFile = false
					m.toBeDeleted = append(m.toBeDeleted, toDelete{
						path: filePathToDelete,
						at:   time.Now().Add(6 * time.Second),
					})
					m.list()
					m.previewContent = ""
					return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
						return toBeDeletedMsg(0)
					})
				}
				m.deleteCurrentFile = true
			}
			return m, nil

		case key.Matches(msg, m.kb.undo):
			if len(m.toBeDeleted) > 0 {
				m.toBeDeleted = m.toBeDeleted[:len(m.toBeDeleted)-1]
				m.list()
				m.previewContent = ""
				return m, nil
			}
		case key.Matches(msg, m.kb.yank):
			// copy path to clipboard
			clipboard.WriteAll(m.path)
			m.yankSuccess = true
			return m, nil
		} // End of switch statement for key presses.

		m.deleteCurrentFile = false
		m.yankSuccess = false
		m.updateOffset()
		m.saveCursorPosition()

	case clearSearchMsg:
		if m.searchId == int(msg) {
			m.searchMode = false
		}

	case toBeDeletedMsg:
		toBeDeleted := make([]toDelete, 0)
		for _, td := range m.toBeDeleted {
			if td.at.After(time.Now()) {
				toBeDeleted = append(toBeDeleted, td)
			} else {
				_ = os.RemoveAll(td.path)
			}
		}
		m.toBeDeleted = toBeDeleted
		if len(m.toBeDeleted) > 0 {
			return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
				return toBeDeletedMsg(0)
			})
		}
	}

	return m, nil
}

func (m *Model) View() string {
	width := m.width
	if m.previewMode {
		width = m.width / 2
	}
	height := m.listHeight()

	var names [][]string
	names, m.rows, m.columns = wrap(m.files, width, height, func(name string, i, j int) {
		if m.findPrevName && m.prevName == name {
			m.c = i
			m.r = j
		}
	})

	// If we need to select previous directory on "up".
	if m.findPrevName {
		m.findPrevName = false
		m.updateOffset()
		m.saveCursorPosition()
	}

	// After we have updated offset and saved cursor position, we can
	// preview currently selected file.
	m.preview()

	// Get output rows width before coloring.
	outputWidth := len(path.Base(m.path)) // Use current dir name as default.
	if m.previewMode {
		row := make([]string, m.columns)
		for i := 0; i < m.columns; i++ {
			if len(names[i]) > 0 {
				row[i] = names[i][0]
			} else {
				outputWidth = width
			}
		}
		outputWidth = max(outputWidth, len(Join(row, separator)))
	} else {
		outputWidth = width
	}

	// Let's add colors to file names.
	output := make([]string, m.rows)
	for j := 0; j < m.rows; j++ {
		row := make([]string, m.columns)
		for i := 0; i < m.columns; i++ {
			if i == m.c && j == m.r {
				if m.deleteCurrentFile {
					row[i] = m.st.Danger.Render(names[i][j])
				} else {
					row[i] = m.st.Cursor.Render(names[i][j])
				}
			} else {
				row[i] = names[i][j]
			}
		}
		output[j] = Join(row, separator)
	}

	if len(output) >= m.offset+height {
		output = output[m.offset : m.offset+height]
	}

	// Preview pane.
	fileName, _ := m.fileName()
	previewPane := m.st.Bar.Render(fileName) + "\n"
	previewPane += m.previewContent

	// Location bar (grey).
	location := m.path
	if userHomeDir, err := os.UserHomeDir(); err == nil {
		location = Replace(m.path, userHomeDir, "~", 1)
	}
	if runtime.GOOS == "windows" {
		location = ReplaceAll(Replace(location, "\\/", fileSeparator, 1), "/", fileSeparator)
	}

	// Filter bar (green).
	filter := ""
	if m.searchMode {
		location = TrimSuffix(location, fileSeparator)
		filter = fileSeparator + m.search
	}
	barLen := len(location) + len(filter)
	if barLen > outputWidth {
		location = location[min(barLen-outputWidth, len(location)):]
	}
	barStr := m.st.Bar.Render(location) + m.st.Search.Render(filter)

	main := barStr + "\n" + Join(output, "\n")

	if m.err != nil {
		main = barStr + "\n" + m.st.Warning.Render(m.err.Error())
	} else if len(m.files) == 0 {
		main = barStr + "\n" + m.st.Warning.Render("No files")
	}

	// Delete bar.
	if len(m.toBeDeleted) > 0 {
		toDelete := m.toBeDeleted[len(m.toBeDeleted)-1]
		timeLeft := int(time.Until(toDelete.at).Seconds())
		deleteBar := fmt.Sprintf("%v deleted. (u)ndo %v", path.Base(toDelete.path), timeLeft)
		main += "\n" + m.st.Danger.Render(deleteBar)
	}

	// Yank success.
	if m.yankSuccess {
		yankBar := fmt.Sprintf("yanked path to clipboard: %v", m.path)
		main += "\n" + m.st.Bar.Render(yankBar)
	}

	if m.previewMode {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			main,
			m.st.Preview.
				MaxHeight(m.height).
				Render(previewPane),
		)
	}
	return main
}

func (m *Model) moveUp() {
	m.r--
	if m.r < 0 {
		m.r = m.rows - 1
		m.c--
	}
	if m.c < 0 {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.files))
		m.c = m.columns - 1
	}
}

func (m *Model) moveDown() {
	m.r++
	if m.r >= m.rows {
		m.r = 0
		m.c++
	}
	if m.c >= m.columns {
		m.c = 0
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.files) {
		m.r = 0
		m.c = 0
	}
}

func (m *Model) moveLeft() {
	m.c--
	if m.c < 0 {
		m.c = m.columns - 1
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.files) {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.files))
		m.c = m.columns - 1
	}
}

func (m *Model) moveRight() {
	m.c++
	if m.c >= m.columns {
		m.c = 0
	}
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.files) {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.files))
		m.c = m.columns - 1
	}
}

func (m *Model) moveTop() {
	m.r = 0
}

func (m *Model) moveBottom() {
	m.r = m.rows - 1
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.files) {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.files))
	}
}

func (m *Model) moveLeftmost() {
	m.c = 0
}

func (m *Model) moveRightmost() {
	m.c = m.columns - 1
	if m.c == m.columns-1 && (m.columns-1)*m.rows+m.r >= len(m.files) {
		m.r = m.rows - 1 - (m.columns*m.rows - len(m.files))
		m.c = m.columns - 1
	}
}

func (m *Model) moveStart() {
	m.moveLeftmost()
	m.moveTop()
}

func (m *Model) moveEnd() {
	m.moveRightmost()
	m.moveBottom()
}

func (m *Model) list() {
	var err error
	m.files = nil

	// ReadDir already returns files and dirs sorted by filename.
	files, err := os.ReadDir(m.path)
	if err != nil {
		m.err = err
		return
	}
	m.err = nil

files:
	for _, file := range files {
		for _, toDelete := range m.toBeDeleted {
			if path.Join(m.path, file.Name()) == toDelete.path {
				continue files
			}
		}
		m.files = append(m.files, file)
	}
}

func (m *Model) listHeight() int {
	h := m.height - 1 // Subtract 1 for location bar.
	if len(m.toBeDeleted) > 0 {
		h-- // Subtract 1 for delete bar.
	}
	return h
}

func (m *Model) updateOffset() {
	height := m.listHeight()
	// Scrolling down.
	if m.r >= m.offset+height {
		m.offset = m.r - height + 1
	}
	// Scrolling up.
	if m.r < m.offset {
		m.offset = m.r
	}
	// Don't scroll more than there are rows.
	if m.offset > m.rows-height && m.rows > height {
		m.offset = m.rows - height
	}
}

func (m *Model) saveCursorPosition() {
	m.positions[m.path] = position{
		c:      m.c,
		r:      m.r,
		offset: m.offset,
	}
}

func (m *Model) fileName() (string, bool) {
	i := m.c*m.rows + m.r
	if i >= len(m.files) || i < 0 {
		return "", false
	}
	return m.files[i].Name(), true
}

func (m *Model) filePath() (string, bool) {
	fileName, ok := m.fileName()
	if !ok {
		return fileName, false
	}
	return path.Join(m.path, fileName), true
}

func (m *Model) openCommand() tea.Cmd {
	filePath, ok := m.filePath()
	if !ok {
		return nil
	}

	cmdline := m.cmdline
	if len(cmdline) == 0 || cmdline[0] == "" {
		cmdline = Fields(lookup([]string{"LK_COMMAND", "EDITOR"}, "less"))
	}
	var replace bool
	for i, s := range cmdline {
		if replace = Contains(s, "{}"); replace {
			cmdline[i] = ReplaceAll(s, "{}", filePath)
			break
		}
	}
	if !replace {
		cmdline = append(cmdline, filePath)
	}

	execCmd := exec.Command(cmdline[0], cmdline[1:]...)
	return tea.ExecProcess(execCmd, func(err error) tea.Msg {
		// Note: we could return a message here indicating that editing is
		// finished and altering our application about any errors. For now,
		// however, that's not necessary.
		return nil
	})
}

func (m *Model) preview() {
	if !m.previewMode {
		return
	}
	filePath, ok := m.filePath()
	if !ok {
		return
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}

	width := m.width / 2
	height := m.height - 1 // Subtract 1 for name bar.

	if fileInfo.IsDir() {
		files, err := os.ReadDir(filePath)
		if err != nil {
			m.previewContent = err.Error()
		}

		names, rows, columns := wrap(files, width, height, nil)

		output := make([]string, rows)
		for j := 0; j < rows; j++ {
			row := make([]string, columns)
			for i := 0; i < columns; i++ {
				row[i] = names[i][j]
			}
			output[j] = Join(row, separator)
		}
		if len(output) >= height {
			output = output[0:height]
		}
		m.previewContent = Join(output, "\n")
		return
	}

	if isImageExt(filePath) {
		img, err := drawImage(filePath, width, height)
		if err != nil {
			m.previewContent = m.st.Warning.Render("No image preview available")
			return
		}
		m.previewContent = img
		return
	}

	var content []byte
	// If file is too big (> 100kb), read only first 100kb.
	if fileInfo.Size() > 100*1024 {
		file, err := os.Open(filePath)
		if err != nil {
			m.previewContent = err.Error()
			return
		}
		defer file.Close()
		content = make([]byte, 100*1024)
		_, err = file.Read(content)
		if err != nil {
			m.previewContent = err.Error()
			return
		}
	} else {
		content, err = os.ReadFile(filePath)
		if err != nil {
			m.previewContent = err.Error()
			return
		}
	}

	switch {
	case utf8.Valid(content):
		m.previewContent = leaveOnlyAscii(content)
	default:
		m.previewContent = m.st.Warning.Render("No preview available")
	}
}

func leaveOnlyAscii(content []byte) string {
	var result []byte

	for _, b := range content {
		if b == '\t' {
			result = append(result, ' ', ' ', ' ', ' ')
		} else if b == '\r' {
			continue
		} else if (b >= 32 && b <= 127) || b == '\n' { // '\n' is kept if newline needs to be retained
			result = append(result, b)
		}
	}

	return string(result)
}

func wrap(files []os.DirEntry, width int, height int, callback func(name string, i, j int)) ([][]string, int, int) {
	// If it's possible to fit all files in one column on a third of the screen,
	// just use one column. Otherwise, let's squeeze listing in half of screen.
	columns := len(files) / (height / 3)
	if columns <= 0 {
		columns = 1
	}

start:
	// Let's try to fit everything in terminal width with this many columns.
	// If we are not able to do it, decrease column number and goto start.
	rows := int(math.Ceil(float64(len(files)) / float64(columns)))
	names := make([][]string, columns)
	n := 0

	for i := 0; i < columns; i++ {
		names[i] = make([]string, rows)
		// Columns size is going to be of max file name size.
		max := 0
		for j := 0; j < rows; j++ {
			name := ""
			if n < len(files) {
				if showIcons {
					info, err := files[n].Info()
					if err == nil {
						icon := icons.getIcon(info)
						if icon != "" {
							name += icon + " "
						}
					}
				}
				name += files[n].Name()
				if callback != nil {
					callback(files[n].Name(), i, j)
				}
				if files[n].IsDir() {
					// Dirs should have a slash at the end.
					name += fileSeparator
				}
				n++
			}
			if max < len(name) {
				max = len(name)
			}
			names[i][j] = name
		}
		// Append spaces to make all names in one column of same size.
		for j := 0; j < rows; j++ {
			names[i][j] += Repeat(" ", max-len(names[i][j]))
		}
	}
	for j := 0; j < rows; j++ {
		row := make([]string, columns)
		for i := 0; i < columns; i++ {
			row[i] = names[i][j]
		}
		if len(Join(row, separator)) > width && columns > 1 {
			// Yep. No luck, let's decrease number of columns and try one more time.
			columns--
			goto start
		}
	}
	return names, rows, columns
}

func (m *Model) dontDoPendingDeletions() {
	for _, toDelete := range m.toBeDeleted {
		fmt.Fprintf(os.Stderr, "Was not deleted: %v\n", toDelete.path)
	}
}

func (m *Model) performPendingDeletions() {
	for _, toDelete := range m.toBeDeleted {
		_ = os.RemoveAll(toDelete.path)
	}
	m.toBeDeleted = nil
}

func fileInfo(path string) os.FileInfo {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fi
}

func lookup(names []string, val string) string {
	for _, name := range names {
		val, ok := os.LookupEnv(name)
		if ok && val != "" {
			return val
		}
	}
	return val
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
