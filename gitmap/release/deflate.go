// Package release — deflate.go provides a max-compression Deflate writer.
package release

import (
	"compress/flate"
	"io"
)

// maxDeflateWriter wraps a flate.Writer at BestCompression (level 9).
type maxDeflateWriter struct {
	fw *flate.Writer
}

// newMaxDeflateWriter creates a Deflate writer with level 9 compression.
func newMaxDeflateWriter(w io.Writer) io.WriteCloser {
	fw, _ := flate.NewWriter(w, flate.BestCompression)

	return &maxDeflateWriter{fw: fw}
}

func (m *maxDeflateWriter) Write(p []byte) (int, error) {
	return m.fw.Write(p)
}

func (m *maxDeflateWriter) Close() error {
	return m.fw.Close()
}
