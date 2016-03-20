package strings

import (
	"errors"
	"io"
)

// View contains a read-only, concurency safe, sharable string.
type View struct {
	s string
}

// NewView creates a new View of a string.
func NewView(s string) *View { return &View{s} }

var errInvalidOffset = errors.New("invalid offset")

// ReadAt satisfies io.ReaderAt
func (v *View) ReadAt(p []byte, off int64) (n int, err error) {
	limit := v.Size()
	if off < 0 || off > limit {
		return 0, errInvalidOffset
	}

	space := len(p)
	if space == 0 && off == limit {
		return 0, io.EOF
	}
	n = copy(p, v.s[off:])
	if n < space {
		return n, io.EOF
	}
	return n, nil
}

// Size returns the size
func (v *View) Size() int64 {
	return int64(len(v.s))
}
