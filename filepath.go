package walk

import (
	"container/list"
	"os"
	"strings"
)

type FilePath struct {
	*list.List
}

func (f *FilePath) init(path string) *FilePath {
	f.List = list.New()
	for _, v := range strings.Split(path, string(os.PathSeparator)) {
		f.PushBack(v)
	}
	return f
}

func (f *FilePath) push(s string) bool {
	f.PushBack(s)
	return true
}

func (f *FilePath) pop() (s string, ok bool) {
	if f.Len() > 0 {
		s, ok = f.Remove(f.Back()).(string)
	}
	return
}

func (f *FilePath) basename() string {
	if f.Len() > 0 {
		return f.Back().Value.(string)
	}
	return ""
}

func (f *FilePath) path() string {
	numElems := f.Len()
	if numElems == 0 {
		return ""
	}
	var sb strings.Builder
	totalLen := 0
	for e := f.Front(); e != nil; e = e.Next() {
		totalLen += len(e.Value.(string))
	}
	sb.Grow(totalLen + (numElems - 1))
	// Append each element in f.List joined by os.PathSeparator
	for e := f.Front(); e != nil; {
		sb.WriteString(e.Value.(string))
		e = e.Next()
		if e != nil {
			sb.WriteByte(os.PathSeparator)
		}
	}
	return sb.String()
}
