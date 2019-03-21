package textwrap

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

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
