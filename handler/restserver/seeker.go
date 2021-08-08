package restserver

import (
	"bytes"
	"errors"
	"io"
)

type readSeeker struct {
	io.Reader
	buf         *bytes.Buffer
	i           int64
	max         int64
	completed   bool
	bufferError bool
}

const (
	DefaultMaxBufferSize int64 = 1024 * 1024
)

var (
	ErrInvalidWhence    = errors.New("invalid whence")
	ErrNegativePosition = errors.New("negative position")
	ErrMaxSizeExceeded  = errors.New("max size exceeded")
)

func (r readSeeker) availableBuffer() int64 {
	return r.max - int64(r.buf.Len())
}

func (r *readSeeker) Read(p []byte) (n int, err error) {
	if r.bufferError {
		return 0, ErrMaxSizeExceeded
	}
	needToRead := int64(len(p))
	i := 0 // index of starting point to write
	if r.i < int64(r.buf.Len()) {
		// The content to read is already in buffer
		rawBuf := r.buf.Bytes()
		sizeToRead := copy(p, rawBuf[r.i:])
		i += sizeToRead
		r.i += int64(sizeToRead)
		needToRead -= int64(sizeToRead)
	} else if int64(r.buf.Len()) < r.i {
		// Start point is far from location that are already read
		needToSkip := r.i - int64(r.buf.Len())
		size, err := io.CopyN(r.buf, r.Reader, needToSkip)
		if err == io.EOF || size < needToSkip {
			r.completed = true
			return 0, io.EOF
		}
	}
	if needToRead > 0 {
		var read int
		readSize := int64(len(p) - i)
		if readSize > r.availableBuffer() {
			return 0, ErrMaxSizeExceeded
		}
		read, err = r.Reader.Read(p[i : i+int(readSize)])
		if err == io.EOF {
			r.completed = true
		} else if r.availableBuffer() == 0 && read < len(p)-i {
			return 0, ErrMaxSizeExceeded
		}
		if int64(read) < needToRead {
			r.completed = true
		}
		r.buf.Write(p[i : i+read])
		r.i += int64(read)
		i += read
	}
	return i, err
}

func (r *readSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
		r.bufferError = false
	case io.SeekCurrent:
		abs = r.i + offset
		r.bufferError = false
	case io.SeekEnd:
		if !r.completed {
			io.CopyN(r.buf, r.Reader, r.max-int64(r.buf.Len()))
			if r.availableBuffer() == 0 {
				r.bufferError = true // bufferError is only for SeekEnd
			}
			r.completed = true
		}
		abs = int64(r.buf.Len()) + offset
	default:
		return 0, ErrInvalidWhence
	}
	if abs < 0 {
		return 0, ErrNegativePosition
	}
	r.i = abs
	return abs, nil
}
