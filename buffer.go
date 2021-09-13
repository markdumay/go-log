// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

package log

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"regexp"
	"strings"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Types
//======================================================================================================================

// Buffer defines a simple buffer to store logs in memory.
type Buffer []string

// BufferedWriter captures application logs and stores them in a local buffer. Log lines are separated by newline
// characters and are added one at a time.
type BufferedWriter struct {
	writer *ConsoleWriter
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Functions
//======================================================================================================================

// NewBufferedWriter creates a log writer that buffers logs in memory.
func NewBufferedWriter(format Format, noColor bool) *BufferedWriter {
	b := BufferedWriter{}
	buffer := make(Buffer, 0)
	b.writer = NewConsoleWriter(format, noColor, &buffer)
	return &b
}

// Write implements the io.Writer interface for Buffer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	// remove multiple line feeds
	input := strings.TrimSuffix(string(p), "\n")
	re := regexp.MustCompile("^\n{2,}")
	input = re.ReplaceAllString(input, "")
	lines := strings.Split(input, "\n")

	// capture the log lines
	for _, line := range lines {
		if _logger.format == Default || line != "" {
			*b = append(*b, line)
		}
	}
	return len(p), nil
}

// Buffer retrieves a copy of the local buffer managed by BufferedWriter.
func (b *BufferedWriter) Buffer() Buffer {
	if b.writer != nil && b.writer.output != nil {
		if v, ok := b.writer.output.(*Buffer); ok {
			return *v
		}
	}

	return make(Buffer, 0)
}

// Reset removes all existing logs from the local buffer.
func (b *BufferedWriter) Reset() {
	if b.writer != nil {
		buffer := make(Buffer, 0)
		format := b.writer.format
		noColor := b.writer.noColor
		b.writer = NewConsoleWriter(format, noColor, &buffer)
	}
}

// SetFormatting updates the log format and color coding of an existing BufferedWriter.
func (b *BufferedWriter) SetFormatting(format Format, noColor bool) {
	b.writer.SetFormatting(format, noColor)
}

// Write implements the io.Writer interface for BufferedWriter.
func (b *BufferedWriter) Write(p []byte) (n int, err error) {
	return b.writer.Write(p)
}

// Flush writes all buffered logs to the active logger and empties the buffer. Subsequent logs are no longer buffered.
func Flush() {
	_logger.hold = false // remove hold to display next message immediately

	// flush the buffered logs
	if len(_logger.buffer) > 0 {
		Debugf("Flushing buffer with %d log(s)", len(_logger.buffer))
		for _, l := range _logger.buffer {
			log(l.Level, l.Message, l.err)
		}
	}

	// clear the buffer
	_logger.buffer = make([]Message, 0)
}

// Hold instructs the active logger to buffer all incoming logs instead of writing them to current output stream. Use
// Flush to write the buffered logs and to empty the buffer.
func Hold() {
	_logger.hold = true
}

//======================================================================================================================
// endregion
//======================================================================================================================
