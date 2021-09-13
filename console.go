// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

package log

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Types
//======================================================================================================================

// ConsoleWriter implements a log writer that supports different styles of formatting. It uses zerolog.ConsoleWriter
// under the hood.
type ConsoleWriter struct {
	format  Format
	noColor bool
	output  io.Writer
	writer  io.Writer
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Functions
//======================================================================================================================

// newWriter creates a new io.Writer that supports Default formatting and Pretty formatting, next to the default JSON
// formatting provided by zerolog.
func newWriter(format Format, noColor bool, out io.Writer) io.Writer {
	// customize the writer if default or pretty formatting is used
	switch format {
	case Format(Default):
		writer := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339, NoColor: noColor}
		writer.FormatTimestamp = func(i interface{}) string {
			return ""
		}
		writer.FormatLevel = func(i interface{}) string {
			v, ok := i.(string)
			if ok && v == "info" {
				return ""
			}
			return strings.ToUpper(fmt.Sprintf("%-6s", i))
		}
		return writer

	case Format(Pretty):
		writer := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339, NoColor: noColor}
		writer.FormatTimestamp = nil
		writer.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s |", i))
		}
		return writer

	default:
		return out
	}
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Functions
//======================================================================================================================

// NewConsoleWriter creates a new ConsoleWriter that supports Default formatting and Pretty formatting, next to the
// default JSON formatting provided by zerolog.
func NewConsoleWriter(format Format, noColor bool, out io.Writer) *ConsoleWriter {
	w := ConsoleWriter{
		format:  format,
		noColor: noColor,
		output:  out,
		writer:  newWriter(format, noColor, out),
	}

	return &w
}

// SetFormatting updates the log format and color coding of an existing ConsoleWriter.
func (w *ConsoleWriter) SetFormatting(f Format, noColor bool) {
	if w.format != f || w.noColor != noColor {
		w.format = f
		w.noColor = noColor
		w.writer = newWriter(f, noColor, w.output)
	}
}

// Write implements the io.Writer interface for ConsoleWriter.
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

//======================================================================================================================
// endregion
//======================================================================================================================
