package walk

import (
	"math"
	"os"
	. "strings"
)

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
			if max < strlen(name) {
				max = strlen(name)
			}
			names[i][j] = name
		}
		// Append spaces to make all names in one column of same size.
		for j := 0; j < rows; j++ {
			names[i][j] += Repeat(" ", max-strlen(names[i][j]))
		}
	}
	for j := 0; j < rows; j++ {
		row := make([]string, columns)
		for i := 0; i < columns; i++ {
			row[i] = names[i][j]
		}
		if strlen(Join(row, separator)) > width && columns > 1 {
			// Yep. No luck, let's decrease number of columns and try one more time.
			columns--
			goto start
		}
	}
	return names, rows, columns
}
