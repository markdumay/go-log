// Copyright © 2021 Mark Dumay. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be found in the LICENSE file.

package log

//======================================================================================================================
// region Import Statements
//======================================================================================================================

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//======================================================================================================================
// endregion
//======================================================================================================================

//======================================================================================================================
// region Test Functions
//======================================================================================================================

func TestInitLoggerWithWriter(t *testing.T) {
	// define the tests
	type test struct {
		msg    string
		msgf   string
		format Format
		level  Level
		result string
		err    string
	}
	var tests = []test{
		// default logger tests
		{
			msg:    "debug message",
			msgf:   "%s message",
			format: Default,
			level:  DebugLevel,
			result: "DEBUG  debug message",
			err:    "DEBUG  debug message error=debug",
		},
		{
			msg:    "info message",
			msgf:   "%s message",
			format: Default,
			level:  InfoLevel,
			result: "info message",
			err:    "info message error=info",
		},
		{
			msg:    "warn message",
			msgf:   "%s message",
			format: Default,
			level:  WarnLevel,
			result: "WARN   warn message",
			err:    "WARN   warn message error=warn",
		},
		{
			msg:    "error message",
			msgf:   "%s message",
			format: Default,
			level:  ErrorLevel,
			result: "ERROR  error message",
			err:    "ERROR  error message error=error",
		},
		{
			msg:    "fatal message",
			msgf:   "%s message",
			format: Default,
			level:  FatalLevel,
			result: "FATAL  fatal message",
			err:    "FATAL  fatal message error=fatal",
		},

		// pretty logger tests
		{
			msg:    "debug message",
			msgf:   "%s message",
			format: Pretty,
			level:  DebugLevel,
			result: " | DEBUG  | debug message",
			err:    " | DEBUG  | debug message error=debug",
		},
		{
			msg:    "info message",
			msgf:   "%s message",
			format: Pretty,
			level:  InfoLevel,
			result: " | INFO   | info message",
			err:    " | INFO   | info message error=info",
		},
		{
			msg:    "warn message",
			msgf:   "%s message",
			format: Pretty,
			level:  WarnLevel,
			result: " | WARN   | warn message",
			err:    " | WARN   | warn message error=warn",
		},
		{
			msg:    "error message",
			msgf:   "%s message",
			format: Pretty,
			level:  ErrorLevel,
			result: " | ERROR  | error message",
			err:    " | ERROR  | error message error=error",
		},
		{
			msg:    "fatal message",
			msgf:   "%s message",
			format: Pretty,
			level:  FatalLevel,
			result: " | FATAL  | fatal message",
			err:    " | FATAL  | fatal message error=fatal",
		},

		// // json logger tests
		{
			msg:    "debug message",
			msgf:   "%s message",
			format: JSON,
			level:  DebugLevel,
			err:    "debug",
		},
		{
			msg:    "info message",
			msgf:   "%s message",
			format: JSON,
			level:  InfoLevel,
			err:    "info",
		},
		{
			msg:    "warn message",
			msgf:   "%s message",
			format: JSON,
			level:  WarnLevel,
			err:    "warn",
		},
		{
			msg:    "error message",
			msgf:   "%s message",
			format: JSON,
			level:  ErrorLevel,
			err:    "error",
		},
		{
			msg:    "fatal message",
			msgf:   "%s message",
			format: JSON,
			level:  FatalLevel,
			err:    "fatal",
		},
	}

	// run the tests
	_suppressExit = true
	for _, test := range tests {
		// redirect log output to buffer
		w := NewBufferedWriter(JSON, false)
		InitLoggerWithWriter(test.format, true, w)
		SetGlobalLevel(test.level)

		switch test.level {
		case DebugLevel:
			Debug(test.msg)
			Debugf(test.msgf, DebugLevel.String())
			DebugE(errors.New(DebugLevel.String()), "debug message")
			Msg(DebugLevel, test.msg)
			Msgf(DebugLevel, test.msgf, DebugLevel.String())
			MsgE(DebugLevel, errors.New(DebugLevel.String()), "debug message")

		case InfoLevel:
			Info(test.msg)
			Infof(test.msgf, InfoLevel.String())
			InfoE(errors.New(InfoLevel.String()), "info message")
			Msg(InfoLevel, test.msg)
			Msgf(InfoLevel, test.msgf, InfoLevel.String())
			MsgE(InfoLevel, errors.New(InfoLevel.String()), "info message")

		case WarnLevel:
			Warn(test.msg)
			Warnf(test.msgf, WarnLevel.String())
			WarnE(errors.New(WarnLevel.String()), "warn message")
			Msg(WarnLevel, test.msg)
			Msgf(WarnLevel, test.msgf, WarnLevel.String())
			MsgE(WarnLevel, errors.New(WarnLevel.String()), "warn message")

		case ErrorLevel:
			Error(test.msg)
			Errorf(test.msgf, ErrorLevel.String())
			ErrorE(errors.New(ErrorLevel.String()), "error message")
			Msg(ErrorLevel, test.msg)
			Msgf(ErrorLevel, test.msgf, ErrorLevel.String())
			MsgE(ErrorLevel, errors.New(ErrorLevel.String()), "error message")

		case FatalLevel:
			Fatal(test.msg)
			Fatalf(test.msgf, FatalLevel.String())
			FatalE(errors.New(FatalLevel.String()), "fatal message")
			Msg(FatalLevel, test.msg)
			Msgf(FatalLevel, test.msgf, FatalLevel.String())
			MsgE(FatalLevel, errors.New(FatalLevel.String()), "fatal message")
		}

		// test the log results
		got := w.Buffer()
		require.Len(t, got, 6)
		if test.format == JSON {
			m, e := UnmarshalLog([]byte(got[0]))
			require.Nil(t, e)
			assert.Equal(t, test.msg, m.Message)
			assert.Equal(t, int(test.level), int(m.Level))

			m, e = UnmarshalLog([]byte(got[1]))
			require.Nil(t, e)
			assert.Equal(t, test.msg, m.Message)
			assert.Equal(t, int(test.level), int(m.Level))

			m, e = UnmarshalLog([]byte(got[2]))
			require.Nil(t, e)
			assert.Equal(t, test.msg, m.Message)
			assert.Equal(t, int(test.level), int(m.Level))
			assert.Equal(t, test.err, m.Error)

			m, e = UnmarshalLog([]byte(got[3]))
			require.Nil(t, e)
			assert.Equal(t, test.msg, m.Message)
			assert.Equal(t, int(test.level), int(m.Level))

			m, e = UnmarshalLog([]byte(got[4]))
			require.Nil(t, e)
			assert.Equal(t, test.msg, m.Message)
			assert.Equal(t, int(test.level), int(m.Level))

			m, e = UnmarshalLog([]byte(got[5]))
			require.Nil(t, e)
			assert.Equal(t, test.msg, m.Message)
			assert.Equal(t, int(test.level), int(m.Level))
			assert.Equal(t, test.err, m.Error)
		} else {
			assert.Contains(t, got[0], test.result)
			assert.Contains(t, got[1], test.result)
			assert.Contains(t, got[2], test.err)
			assert.Contains(t, got[3], test.result)
			assert.Contains(t, got[4], test.result)
			assert.Contains(t, got[5], test.err)
		}
	}

	// restore the logger settings
	_suppressExit = false
	InitLogger(Default)
	SetGlobalLevel(InfoLevel)
}

func TestLogDirect(t *testing.T) {
	// redirect log output to buffer
	w := NewBufferedWriter(JSON, false)
	InitLoggerWithWriter(JSON, true, w)
	SetGlobalLevel(WarnLevel)

	// log a direct message
	Bypass("Direct message")

	// test the log results and logger settings
	got := w.Buffer()
	require.Len(t, got, 1)
	assert.Equal(t, "Direct message", got[0])
	assert.Equal(t, WarnLevel, GlobalLevel())
}

func TestParseFormat(t *testing.T) {
	type test struct {
		input    string
		expected Format
		err      string
	}

	var tests = []test{
		{input: "default", expected: Default, err: ""},
		{input: "pretty", expected: Pretty, err: ""},
		{input: "json", expected: JSON, err: ""},
		{input: "DEFAULT", expected: Default, err: ""},
		{input: "PRETTY", expected: Pretty, err: ""},
		{input: "JSON", expected: JSON, err: ""},
		{input: "unknown", expected: Default, err: "unknown log format: 'unknown'"},
	}

	for _, test := range tests {
		r, e := ParseFormat(test.input)
		assert.Equal(t, test.expected, r)
		if test.err != "" {
			assert.Equal(t, test.err, e.Error())
		}
	}
}

func TestSetLogFormat(t *testing.T) {
	InitLogger(Default)
	SetFormatting(Pretty, false)

	assert.Equal(t, Pretty, _logger.format)

	// restore the logger settings
	InitLogger(Default)
	SetGlobalLevel(InfoLevel)
}

func TestLogFormatString(t *testing.T) {
	assert.Equal(t, "default", Default.String())
	assert.Equal(t, "pretty", Pretty.String())
	assert.Equal(t, "json", JSON.String())
}

func TestWrite(t *testing.T) {
	buffer := Buffer{}

	// validate simple input
	input1 := "input"
	n, e := buffer.Write([]byte(input1))
	require.Nil(t, e)
	assert.Equal(t, len(input1), n)

	// validate multiline input
	input2 := "multiline\n\ninput"
	n, e = buffer.Write([]byte(input2))
	require.Nil(t, e)
	assert.Equal(t, len(input2), n)

	// validate buffer
	expected := []string{input1, "multiline", "", "input"}
	assert.Equal(t, expected, []string(buffer))
}

//======================================================================================================================
// endregion
//======================================================================================================================
