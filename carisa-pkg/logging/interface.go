/*
 * Copyright 2019-2022 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package logging

// A Level is a logging priority.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
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
	// Info logs a message at InfoLevel. The message includes any fields passed
	// The limit are 5 fields
	Info(msg string, fields ...Field)

	// Warn logs a message at WarnLevel. The message includes any fields passed
	// The limit are 5 fields
	Warn(msg string, fields ...Field)

	// Debug logs a message at DebugLevel. The message includes any fields passed
	// The limit are 5 fields
	Debug(msg string, fields ...Field)

	// Error logs a message at ErrorLevel. The message includes any fields passed
	// The limit are 5 fields
	Error(msg string, fields ...Field)

	// Panic logs a message at PanicLevel. The message includes any fields passed
	// The limit are 5 fields
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic(msg string, fields ...Field)

	// Check checks if logging a message at the specified level
	// is enabled.
	// If disabled, it returns null, otherwise it returns the element to be written
	Check(l Level, msg string) *checkWrap

	// Writes log. See Check
	Write(wrap *checkWrap, fields ...Field)
}
