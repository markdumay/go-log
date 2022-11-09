// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

// Package log is a simplified logger package for Go applications. Using the Zero Allocation JSON Logger
// (zerolog) under the hood, it simplifies the logging of application-wide messages. It supports three logging modes:
// Default, Pretty, and JSON. Logs are directed to the console by default, but can be buffered or redirected to a log
// file instead.
package log

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Constants
//======================================================================================================================

// Defines a pseudo enumeration of possible logging formats.
const (
	// Default prints logs as standard console output (no timestamp and level prefixes), for example:
	// 		// Listing snapshots
	Default Format = iota

	// Pretty prints logs as semi-structured messages with a timestamp and level prefix, for example:
	// 		// 2020-12-17T07:12:57+01:00 | INFO   | Listing snapshots
	Pretty

	// JSON prints logs as JSON strings, for example:
	// 		// {"level":"info","time":"2020-12-17T07:12:57+01:00","message":"Listing snapshots"}
	JSON
)

// Defines a pseudo enumeration of possible logging levels, copied from zerolog to hide implementation details.
const (
	// DebugLevel defines the debugging log level.
	DebugLevel Level = iota

	// InfoLevel defines the info log level.
	InfoLevel

	// WarnLevel defines the warning log level.
	WarnLevel

	// ErrorLevel defines the error log level.
	ErrorLevel

	// FatalLevel defines the fatal log level.
	FatalLevel

	// PanicLevel defines the panic log level.
	PanicLevel

	// NoLevel defines an absent log level.
	NoLevel

	// Disabled disables the logger.
	Disabled

	// TraceLevel defines the trace log level.
	TraceLevel Level = -1
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Variables
//======================================================================================================================

// _logger is used as internal handler for any logs to be created by the functions Info(), Debug(), et al.
var _logger = NewLogger(Default, false)

// _suppressExit suppresses Fatal logs from exiting the program. Used for testing.
var _suppressExit bool

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Types
//======================================================================================================================

// Writer defines the interface for writers to be compatible with Logger. It adds the SetFormatting function to
// ensure custom writers are synchronized with the application settings.
type Writer interface {
	io.Writer
	SetFormatting(format Format, noColor bool)
}

// Logger is a simplified logger that uses zerolog under the hood. It supports three logging modes, being Default,
// Pretty, and JSON. In default mode, all logs are printed using simplified formatting. This format omits timestamps and
// puts a simple keyword in front of the message to indicate the level. For Info logs, the level is omitted. Pretty mode
// structures the logs using a timestamp (RFC 3339) and level indicator, separated by the symbol '|'. Finally, JSON mode
// formats the log as a JSON message, consisting of the attributes timestamp (RFC 3339), level, and message.
//
// A default logger is instantiated by default. The following examples illustrate how to use the package.
//
//	package main
//
//	import (
//		"go.markdumay.org/log"
//	)
//
//	func main() {
//		// show an info message using default formatting, expected output:
//		// This is an info log
//		log.Info("This is an info log")
//
//		// show an error message using default formatting, expected output:
//		// ERROR  Error message
//		log.Info("Error message")
//
//		// switch to pretty formatting
//		log.InitLogger(log.Pretty)
//
//		// show a warning using pretty formatting, expected output:
//		// 2006-01-02T15:04:05Z07:00 | WARN   | Warning
//		log.Warn("Warning")
//
//		// switch to JSON formatting
//		log.InitLogger(log.JSON)
//
//		// switch to debug level as minimum level
//		log.SetGlobalLevel(log.DebugLevel)
//
//		// show a debug message using JSON formatting, expected output:
//		// {"level":"debug","time":"2006-01-02T15:04:05Z07:00","message":"Testing level debug"}
//		log.Debugf("Testing level %s", "debug")
//	}
type Logger struct {
	format  Format
	level   Level
	handler *zerolog.Logger
	writers []Writer
	noColor bool
	buffer  []Message
	hold    bool
}

// Format defines the type of logging format to use, either Default, Pretty, or JSON.
type Format int

// Level defines the minimum level of logs to display. Supported levels are DebugLevel, InfoLevel, WarnLevel,
// ErrorLevel, FatalLevel, and PanicLevel. Level is an abstraction of a type with the same name provided by the
// underlying zerolog package.
type Level int8

// Message defines the structure of JSON-formatted log messages produced by zerolog.
type Message struct {
	Level   Level
	Time    time.Time
	Message string
	Error   string
	err     error
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Private Functions
//======================================================================================================================

// getWriterIndex returns the index of the Writer within the list of writers known by Logger. It returns -1 if the
// writer cannot be found.
func getWriterIndex(w Writer) int {
	for index, curr := range _logger.writers {
		if w == curr {
			return index
		}
	}

	return -1
}

// log is an internal function to redirect logging requests to either the handler or local buffer.
func log(level Level, msg string, err error, v ...interface{}) {
	var m string
	if v != nil {
		m = fmt.Sprintf(msg, v...)
	} else {
		m = msg
	}

	if _logger.hold {
		var log Message
		log.Level = level
		log.Time = time.Now()
		log.Message = m
		log.err = err
		if err != nil {
			log.Error = err.Error()
		}
		_logger.buffer = append(_logger.buffer, log)
	} else {
		if err != nil {
			_logger.handler.WithLevel(zerolog.Level(level)).Err(err).Msg(m)
		} else {
			_logger.handler.WithLevel(zerolog.Level(level)).Msg(m)
		}
	}
}

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Public Functions
//======================================================================================================================

// NewLogger initializes a new logger with the desired format.
func NewLogger(format Format, noColor bool, writer ...Writer) *Logger {
	var writers []Writer

	// add a default console writer if needed
	if len(writer) == 0 {
		writers = append(writers, NewConsoleWriter(format, noColor, os.Stdout))
	} else {
		// update the formatting of the existing writers
		writers = append(writers, writer...)
		for _, w := range writers {
			w.SetFormatting(format, noColor) // instruct the writers to use the defined log format
		}
	}

	// init a zerologger with either a single writer or a multi-level writer
	var l = new(Logger)
	var handler zerolog.Logger
	if len(writers) == 1 {
		handler = zerolog.New(writers[0]).With().Timestamp().Logger()
	} else {
		// Note: compiler complains when using variadic expansion "writers...", therefore convert to []io.Writer first
		var export []io.Writer
		for _, w := range writers {
			export = append(export, w)
		}
		multi := zerolog.MultiLevelWriter(export...)
		handler = zerolog.New(multi).With().Timestamp().Logger()
	}

	// init the logger and return the reference
	l.format = format
	l.writers = writers
	l.noColor = noColor
	l.handler = &handler
	l.buffer = make([]Message, 0)

	return l
}

// Write implements the io.Writer interface for Logger.
func (l *Logger) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		// skip empty lines when not using default logging format
		if line != "" || Format(zerolog.GlobalLevel()) == Format(Default) {
			l.handler.WithLevel(zerolog.Level(l.level)).Msg(line)
		}
	}
	return len(p), nil
}

// MarshalText implements the TextMarshaler interface for Format.
func (f Format) MarshalText() (text []byte, err error) {
	return []byte(f.String()), nil
}

// String converts a typed log format to it's string representation.
func (f Format) String() string {
	if f < Default || f > JSON {
		return ""
	}

	return [...]string{"default", "pretty", "json"}[f]
}

// MarshalText implements the TextMarshaler interface for Level.
func (l Level) MarshalText() (text []byte, err error) {
	return []byte(l.String()), nil
}

// String converts a typed log level to its string representation.
func (l Level) String() string {
	z := zerolog.Level(l)
	return zerolog.Level.String(z)
}

// AppendWriter appends a writer to the list of writers known by Logger. Logs are duplicated for each known writer.
func AppendWriter(w Writer) {
	writers := make([]Writer, len(_logger.writers))
	copy(writers, _logger.writers)
	writers = append(writers, w)
	InitLoggerWithWriter(_logger.format, _logger.noColor, writers...)
}

// Bypass logs an info message using a default logging format, bypassing the current level and format. Use this
// function to ensure custom logs are written as-is to the standardized logging stream(s). If multiple writers are
// specified, the message is duplicated for all writers.
func Bypass(msg string) {
	// back up the current level and format
	level := zerolog.GlobalLevel()
	format := _logger.format
	noColor := _logger.noColor

	// ensure to restore the logger when done
	defer zerolog.SetGlobalLevel(level)
	defer SetFormatting(format, noColor)

	// log a info message with default format
	SetFormatting(Default, true)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	_logger.handler.Info().Msg(msg)
}

// Debug logs a debugging message.
func Debug(msg string) {
	log(DebugLevel, msg, nil)
}

// DebugE logs a debugging error.
func DebugE(e error, msg string) {
	log(DebugLevel, msg, e)
}

// Debugf logs a formatted debugging message.
func Debugf(format string, v ...interface{}) {
	log(DebugLevel, format, nil, v...)
}

// Error logs an error message.
func Error(msg string) {
	log(ErrorLevel, msg, nil)
}

// ErrorE logs an error.
func ErrorE(e error, msg string) {
	log(ErrorLevel, msg, e)
}

// Errorf logs a formatted error message.
func Errorf(format string, v ...interface{}) {
	log(ErrorLevel, format, nil, v...)
}

// Fatal logs a fatal message. It exits the program with exit code 1. Fatal messages are never buffered.
func Fatal(msg string) {
	_logger.handler.WithLevel(zerolog.FatalLevel).Msg(msg)
	if !_suppressExit {
		os.Exit(1)
	}
}

// FatalE logs a fatal error. It exits the program with exit code 1. Fatal messages are never buffered.
func FatalE(e error, msg string) {
	_logger.handler.WithLevel(zerolog.FatalLevel).Err(e).Msg(msg)
	if !_suppressExit {
		os.Exit(1)
	}
}

// Fatalf logs a formatted fatal error. It exits the program with exit code 1. Fatal messages are never buffered.
func Fatalf(format string, v ...interface{}) {
	_logger.handler.WithLevel(zerolog.FatalLevel).Msgf(format, v...)
	if !_suppressExit {
		os.Exit(1)
	}
}

// GlobalLevel retrieves the logging level of all loggers.
func GlobalLevel() Level {
	return Level(zerolog.GlobalLevel())
}

// Info logs a message.
func Info(msg string) {
	log(InfoLevel, msg, nil)
}

// InfoE logs an error.
func InfoE(e error, msg string) {
	log(InfoLevel, msg, e)
}

// Infof logs a formatted message.
func Infof(format string, v ...interface{}) {
	log(InfoLevel, format, nil, v...)
}

// InitLogger initializes the global logger with the desired format. Output is written to STDOUT with color coding.
func InitLogger(format Format) {
	InitLoggerWithWriter(format, true)
}

// InitLoggerWithWriter initializes the global logger with the desired format, writer(s), and color coding.
func InitLoggerWithWriter(format Format, noColor bool, writer ...Writer) {
	b := _logger.buffer
	_logger = NewLogger(format, noColor, writer...)
	_logger.buffer = b
}

// Msg logs a message at the desired level.
func Msg(level Level, msg string) {
	log(level, msg, nil)
}

// MsgE logs an error at the desired level.
func MsgE(level Level, e error, msg string) {
	log(level, msg, e)
}

// Msgf logs a formatted message at the desired level.
func Msgf(level Level, format string, v ...interface{}) {
	log(level, format, nil, v...)
}

// ParseFormat converts a format string into a typed Format value. It returns an error if the input string does not
// match known values.
func ParseFormat(formatStr string) (Format, error) {
	switch strings.ToLower(formatStr) {
	case "default":
		return Format(Default), nil

	case "pretty":
		return Format(Pretty), nil

	case "json":
		return Format(JSON), nil
	}
	return Format(Default), fmt.Errorf("unknown log format: '%s'", formatStr)
}

// ParseLevel converts a level string into a typed Level value. It returns an error if the input string does not
// match known values.
func ParseLevel(levelStr string) (Level, error) {
	l, e := zerolog.ParseLevel(levelStr)
	if e != nil {
		return InfoLevel, e
	}

	return Level(l), nil
}

// RemoveWriter removes a writer from the list of writers known by Logger. The request is ignored when the writer cannot
// be found.
func RemoveWriter(w Writer) {
	index := getWriterIndex(w)
	if index >= 0 {
		writers := append(_logger.writers[:index], _logger.writers[index+1:]...)
		InitLoggerWithWriter(_logger.format, _logger.noColor, writers...)
	}
}

// SetFormatting adjusts the logging format of the current logger.
func SetFormatting(format Format, noColor bool) {
	_logger.format = format
	_logger.noColor = noColor
	for _, w := range _logger.writers {
		w.SetFormatting(format, noColor)
	}
}

// SetGlobalLevel sets the logging level for all loggers.
func SetGlobalLevel(l Level) {
	zerolog.SetGlobalLevel(zerolog.Level(l))
}

// UpdateWriter replaces an old writer from the list of writers known by Logger with a new writer. UpdateWriter returns
// an error if the old writer cannot be found.
func UpdateWriter(old Writer, new Writer) error {
	index := getWriterIndex(old)
	if index < 0 || index >= len(_logger.writers) {
		return errors.New("Cannot update logger stream, current stream not found")
	}

	writers := _logger.writers
	writers[index] = new
	InitLoggerWithWriter(_logger.format, _logger.noColor, writers...)

	return nil
}

// UnmarshalLog converts json bytes into a Message instance.
func UnmarshalLog(bytes []byte) (*Message, error) {
	const layout = "2006-01-02T15:04:05Z07:00"

	// construct a placeholder with looser typing
	raw := struct {
		Level   string `json:"level"`
		Time    string `json:"time"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}{}

	// convert json input to placeholder type
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, err
	}

	// convert input to typed timestamp, fail on error
	timestamp, err := time.Parse(layout, raw.Time)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse datetime format, got %s, want %s", raw.Time, layout)
	}

	// parse Level
	level, err := zerolog.ParseLevel(raw.Level)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse level: %s", raw.Level)
	}

	// convert placeholder type to final type
	log := &Message{
		Level:   Level(level),
		Time:    timestamp,
		Message: raw.Message,
		Error:   raw.Error,
	}

	return log, nil
}

// Warn logs a warning.
func Warn(msg string) {
	log(WarnLevel, msg, nil)
}

// WarnE logs an error as warning.
func WarnE(e error, msg string) {
	log(WarnLevel, msg, e)
}

// Warnf logs a formatted warning.
func Warnf(format string, v ...interface{}) {
	log(WarnLevel, format, nil, v...)
}

//======================================================================================================================
// endregion
//======================================================================================================================
