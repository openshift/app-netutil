// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2021 Red Hat, Inc.

//
// This module abstracts the logging so logging mdule can be
// swapped out without changing all the files.
//

package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Level type
type Level uint32

// ErrorLevel...MaxLevel indicates the logging level
const (
	ErrorLevel Level = iota
	WarningLevel
	InfoLevel
	DebugLevel
	MaxLevel
	UnknownLevel
)

var loggingStderr bool
var loggingW io.Writer
var loggingLevel Level

const defaultTimestampFormat = time.RFC3339

func (loglevel Level) String() string {
	switch loglevel {
	case ErrorLevel:
		return "error"
	case WarningLevel:
		return "warning"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	}
	return "unknown"
}

func printf(level Level, format string, a ...interface{}) {
	header := "%s [%s] "
	t := time.Now()
	if level > loggingLevel {
		return
	}

	if loggingStderr {
		fmt.Fprintf(os.Stderr, header, t.Format(defaultTimestampFormat), level)
		fmt.Fprintf(os.Stderr, format, a...)
		fmt.Fprintf(os.Stderr, "\n")
	}

	if loggingW != nil {
		fmt.Fprintf(loggingW, header, t.Format(defaultTimestampFormat), level)
		fmt.Fprintf(loggingW, format, a...)
		fmt.Fprintf(loggingW, "\n")
	}
}

// Debugf prints logging if logging level >= debug
func Debugf(format string, a ...interface{}) {
	printf(DebugLevel, format, a...)
}

// Verbosef prints logging if logging level >= info
func Infof(format string, a ...interface{}) {
	printf(InfoLevel, format, a...)
}

// Warningf prints logging if logging level >= warning
func Warningf(format string, a ...interface{}) error {
	printf(WarningLevel, format, a...)
	return fmt.Errorf(format, a...)
}

// Errorf prints logging if logging level >= error
func Errorf(format string, a ...interface{}) error {
	printf(ErrorLevel, format, a...)
	return fmt.Errorf(format, a...)
}

// GetLoggingLevel gets current logging level
func GetLoggingLevel() Level {
	return loggingLevel
}

func getLoggingLevel(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warning":
		return WarningLevel
	case "error":
		return ErrorLevel
	}
	fmt.Fprintf(os.Stderr, "app-netutil logging: cannot set logging level to %s\n", levelStr)
	return UnknownLevel
}

// SetLogLevel sets logging level
func SetLogLevel(levelStr string) {
	level := getLoggingLevel(levelStr)
	if level < MaxLevel {
		loggingLevel = level
	}
}

// SetLogStderr sets flag for logging stderr output
func SetLogStderr(enable bool) {
	loggingStderr = enable
}

// SetLogFile sets logging file
func SetLogFile(filename string) {
	if filename == "" {
		return
	}

	loggingW = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100, // megabytes
		MaxBackups: 5,
		MaxAge:     5, // days
		Compress:   true,
	}

}

func init() {
	loggingStderr = true
	loggingW = nil
	loggingLevel = WarningLevel
}
