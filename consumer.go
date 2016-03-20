package views

import (
	"errors"
	"io"
)

// SizeReaderAt is a io.ReaderAt that also know the upper boundary (it's size) of where to read.
// Users can assume it will never change it's size again.
type SizeReaderAt interface {
	io.ReaderAt
	Size() int64
}

var errClosed = errors.New("stream already closed")

type Consumer struct {
	r      SizeReaderAt
	i      int64
	limit  int64
	closed bool
}

// NewConsumer turns a SizeReaderAt into a consumable stream satisfying io.Reader.
func NewConsumer(r SizeReaderAt) *Consumer {
	return &Consumer{
		r:     r,
		i:     0,
		limit: r.Size(),
	}
}

func (c *Consumer) Read(p []byte) (n int, err error) {
	if c.closed {
		return 0, errClosed
	}
	n, err = c.r.ReadAt(p, c.i)
	c.i += int64(n)
	return n, err
}

// Reset starts reading from the beginning of r.
func (c *Consumer) Reset(r SizeReaderAt) {
	c.limit = r.Size()
	c.r = r
	c.i = 0
	c.closed = false
}

// Size returns the intial size of the underlying SizeReaderAt.
// It continues to work after calling Close, but the next Read will return an error.
func (c *Consumer) Size() int64 {
	return c.limit
}

// Rewind positions the stream at the beginning.
// It continues to work after calling Close, but the next Read will return an error.
func (c *Consumer) Rewind() {
	c.i = 0
}

// Close closes this stream. Any further calls to Read or will return an error now.
// Calling Close on a closed stream is an error, too.
// It also closes the underlying SizeReaderAt, if it is satisfies io.Closer.
func (c *Consumer) Close() error {
	if c.closed {
		return errClosed
	}
	c.closed = true

	cl, ok := c.r.(io.Closer)
	if !ok {
		return nil
	}
	return cl.Close()
}
