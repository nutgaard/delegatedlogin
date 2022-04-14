package internal

import (
	"bytes"
	"io"
	"regexp"
)

type maskingWriter struct {
	io.Writer
	pattern     *regexp.Regexp
	replacement []byte
}

func CreateMaskingWriter(
	pattern string,
	replacement string,
	writer io.Writer,
) io.Writer {
	return &maskingWriter{
		Writer:      writer,
		pattern:     regexp.MustCompile(pattern),
		replacement: []byte(replacement),
	}
}

func (m maskingWriter) Write(src []byte) (n int, err error) {
	return m.Writer.Write(m.pattern.ReplaceAllFunc(src, func(match []byte) []byte {
		return bytes.Repeat(m.replacement, len(match))
	}))
}
