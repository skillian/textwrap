package textwrap

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

// SliceLines receives a slice of individual values that are wrapped into a
// slice of lines.
func SliceLines(vs []string, columns int, separator string) []string {
	lines := make([][]string, 1, len(vs))
	c := 0
	for _, v := range vs {
		length := len(v) + len(separator)
		if length+c > columns {
			lines = append(lines, []string{v})
			c = length
		} else {
			line := &lines[len(lines)-1]
			*line = append(*line, v)
			c += length
		}
	}
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = strings.Join(line, separator)
	}
	return out
}

// Slice wraps the given slice of string values.  The values cannot be broken
// across lines even if they contain spaces.
func Slice(vs []string, columns int) string {
	b := strings.Builder{}
	c := 0
	for _, v := range vs {
		length := len(v)
		c += length
		if c > columns {
			b.WriteRune('\n')
			c = length
		}
		b.WriteString(v)
	}
	return b.String()
}

// String wraps a string to the given number of columns
func String(v string, columns int) string {
	src := bytes.NewBufferString(v)
	trg := new(bytes.Buffer)
	r, err := NewReader(ioutil.NopCloser(src), columns)
	if err != nil {
		panic(err)
	}
	if _, err = io.Copy(trg, r); err != nil {
		panic(err)
	}
	return trg.String()
}

// Reader wraps a ReadCloser and wraps the text it reads from that ReadCloser
// to a certain number of lines.
type Reader struct {
	// ReadCloser holds the underlying ReadCloser whose bytes are read from.
	io.ReadCloser

	// Columns holds the number of columns where the reader should wrap
	// the underlying text.
	Columns int

	// lastcol holds the column number of the last incomplete line read.
	lastcol int
}

var (
	errReadCloserNil   = errors.New("io.ReadCloser cannot be nil")
	errColumnsNegative = errors.New("columns cannot be negative")
)

// NewReader constructs a Reader with the given parameters.
//
// Note that right now, you could just use the type directly because
// its private members are valid in their zero-state, but I'm still making
// changes, so that might not always be the case.
func NewReader(rc io.ReadCloser, columns int) (*Reader, error) {
	if rc == nil {
		return nil, errReadCloserNil
	}
	if columns < 0 {
		return nil, errColumnsNegative
	}
	r := &Reader{
		ReadCloser: rc,
		Columns:    columns,
		lastcol:    0,
	}
	return r, nil
}

// Read implements the io.Reader interface.
//
// Read calls the underlying ReadCloser's read function and then replaces
// space (' ') characters with newlines ('\n') to make sure that every line
// is less than the Reader's Columns value.
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	lastcol := r.lastcol
	lastspc := lastcol
	for i, b := range p[:n] {
		switch b {
		case '\n':
			lastcol = 0
		case ' ':
			lastspc = i
			fallthrough
		default:
			lastcol++
		}
		if lastcol > r.Columns {
			p[lastspc] = '\n'
			lastcol = i - lastspc
		}
	}
	r.lastcol = lastcol
	return
}
