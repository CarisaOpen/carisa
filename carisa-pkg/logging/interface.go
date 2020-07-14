/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package logging

import (
	"github.com/pkg/errors"
)

// A Level is a logging priority.
type Level uint8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota + 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	PanicLevel
)

// Looger implements operations for several logging platforms
type Logger interface {
	// Level gets the configured level
	Level() Level

	// Info logs a message at InfoLevel. The message includes any fields passed and location
	Info(string, string)

	// Info1 logs a message at InfoLevel. The message includes any fields passed and location
	Info1(string, string, Field)

	// Info2 logs a message at InfoLevel. The message includes any fields passed and location
	Info2(string, string, Field, Field)

	// Info3 logs a message at InfoLevel. The message includes any fields passed and location
	Info3(string, string, Field, Field, Field)

	// Warn logs a message at WarnLevel. The message includes any fields passed and location
	Warn(string, string)

	// Warn1 logs a message at WarnLevel. The message includes any fields passed and location
	Warn1(string, string, Field)

	// Warn2 logs a message at WarnLevel. The message includes any fields passed and location
	Warn2(string, string, Field, Field)

	// Warn3 logs a message at WarnLevel. The message includes any fields passed and location
	Warn3(string, string, Field, Field, Field)

	// Debug logs a message at DebugLevel. The message includes any fields passed and location
	Debug(string, string)

	// Debug1 logs a message at DebugLevel. The message includes any fields passed and location
	Debug1(string, string, Field)

	// Debug2 logs a message at DebugLevel. The message includes any fields passed and location
	Debug2(string, string, Field, Field)

	// Debug3 logs a message at DebugLevel. The message includes any fields passed and location
	Debug3(string, string, Field, Field, Field)

	// Error logs a message at ErrorLevel. The message includes any fields passed and location
	Error(string, string)

	// Error1 logs a message at ErrorLevel. The message includes any fields passed and location
	Error1(string, string, Field)

	// Error2 logs a message at ErrorLevel. The message includes any fields passed and location
	Error2(string, string, Field, Field)

	// Error3 logs a message at ErrorLevel. The message includes any fields passed and location
	Error3(string, string, Field, Field, Field)

	// ErrWrap generates a trace in logging and return a wrap error
	ErrWrap(error, string, string) error

	// ErrWrap1 generates a trace in logging and return a wrap error
	ErrWrap1(error, string, string, Field) error

	// ErrWrap2 generates a trace in logging and return a wrap error
	ErrWrap2(error, string, string, Field, Field) error

	// ErrWrap3 generates a trace in logging and return a wrap error
	ErrWrap3(error, string, string, Field, Field, Field) error

	// ErrorE logs a message at ErrorLevel. The message includes any fields passed and location
	ErrorE(error, string)

	// ErrorE1 logs a message at ErrorLevel. The message includes any fields passed and location
	ErrorE1(error, string, Field)

	// ErrorE2 logs a message at ErrorLevel. The message includes any fields passed and location
	ErrorE2(error, string, Field, Field)

	// ErrorE3 logs a message at ErrorLevel. The message includes any fields passed and location
	ErrorE3(error, string, Field, Field, Field)

	// Panic logs a message at PanicLevel. The message includes any fields passed and location
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic(string, string)

	// Panic1 logs a message at PanicLevel. The message includes any fields passed and location
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic1(string, string, Field)

	// Panic2 logs a message at PanicLevel. The message includes any fields passed and location
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic2(string, string, Field, Field)

	// Panic3 logs a message at PanicLevel. The message includes any fields passed and location
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic3(string, string, Field, Field, Field)

	// Panic logs a message at PanicLevel. The message includes the location
	// The logger then panics, even if logging at PanicLevel is disabled.
	PanicE(err error, loc string)

	// Check checks if logging a message at the specified level
	// is enabled.
	// If disabled, it returns null, otherwise it returns the element to be written
	Check(Level, string) *CheckWrap

	// Write writes log. See Check
	Write(*CheckWrap, string)

	// Write1 writes log. See Check
	Write1(*CheckWrap, string, Field)

	// Write2 writes log. See Check
	Write2(*CheckWrap, string, Field, Field)

	// Write3 writes log. See Check
	Write3(*CheckWrap, string, Field, Field, Field)
}

// Looger composition
type loggerComp struct {
	log Logger
}

// ErrorE implements logging.Logger.ErrorE
func (l *loggerComp) ErrorE(err error, loc string) {
	l.log.Error(err.Error(), loc)
}

// ErrorE1 implements logging.Logger.ErrorE
func (l *loggerComp) ErrorE1(err error, loc string, f Field) {
	l.log.Error1(err.Error(), loc, f)
}

// ErrorE2 implements logging.Logger.ErrorE1
func (l *loggerComp) ErrorE2(err error, loc string, f Field, f1 Field) {
	l.log.Error2(err.Error(), loc, f, f1)
}

// ErrorE3 implements logging.Logger.ErrorE2
func (l *loggerComp) ErrorE3(err error, loc string, f Field, f1 Field, f2 Field) {
	l.log.Error3(err.Error(), loc, f, f1, f2)
}

// ErrWrap implements logging.Logger.ErrWrap
func (l *loggerComp) ErrWrap(err error, msg string, loc string) error {
	errW := errors.Wrap(err, msg)
	l.log.Error(errW.Error(), loc)
	return errW
}

// ErrWrap1 implements logging.Logger.ErrWrap1
func (l *loggerComp) ErrWrap1(err error, msg string, loc string, f Field) error {
	errW := errors.Wrap(err, Compose(msg, f))
	l.log.Error(errW.Error(), loc)
	return errW
}

// ErrWrap2 implements logging.Logger.ErrWrap2
func (l *loggerComp) ErrWrap2(err error, msg string, loc string, f Field, f1 Field) error {
	errW := errors.Wrap(err, Compose(msg, f, f1))
	l.log.Error(errW.Error(), loc)
	return errW
}

// ErrWrap3 implements logging.Logger.ErrWrap3
func (l *loggerComp) ErrWrap3(err error, msg string, loc string, f Field, f1 Field, f2 Field) error {
	errW := errors.Wrap(err, Compose(msg, f, f1, f2))
	l.log.Error(errW.Error(), loc)
	return errW
}

// PanicE implements logging.Logger.PanicE
func (l *loggerComp) PanicE(err error, loc string) {
	l.log.Panic(err.Error(), loc)
}
