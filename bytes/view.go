package bytes

import (
	"errors"
	"io"
)

type View struct {
	b []byte
}

// NewView means you promise to not modify b afer this call and consider b read-only after this
func NewView(b []byte) *View { return &View{b} }

var errInvalidOffset = errors.New("invalid offset")

func (v *View) ReadAt(p []byte, off int64) (n int, err error) {
	limit := v.Size()
	if off < 0 || off > limit {
		return 0, errInvalidOffset
	}

	space := len(p)
	if space == 0 && off == limit {
		return 0, io.EOF
	}
	n = copy(p, v.b[off:])
	if n < space {
		return n, io.EOF
	}
	return n, nil
}

func (v *View) Size() int64 {
	return int64(len(v.b))
}
