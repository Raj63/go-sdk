package qrcode

import "bytes"

// ByteArrayWriteCloser is a wrapper around a byte array that satisfies the io.WriteCloser interface.
type ByteArrayWriteCloser struct {
	buf *bytes.Buffer
}

// NewByteArrayWriteCloser creates a new ByteArrayWriteCloser with the provided byte array.
func NewByteArrayWriteCloser() *ByteArrayWriteCloser {
	return &ByteArrayWriteCloser{buf: bytes.NewBuffer(nil)}
}

// Write writes data to the underlying byte array.
func (b *ByteArrayWriteCloser) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

// Close closes the underlying byte array.
func (b *ByteArrayWriteCloser) Close() error {
	return nil
}

// Bytes returns the underlying byte array.
func (b *ByteArrayWriteCloser) Bytes() []byte {
	return b.buf.Bytes()
}
