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

package echo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/carisa/pkg/logging"
	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
)

// zapLogger defines a zap wrapper for echo http
type zapLogger struct {
	level  log.Lvl
	logB   *zap.Logger
	prefix string
	output io.Writer
}

func NewLogging(prefix string, level log.Lvl, log *zap.Logger) echo.Logger {
	return &zapLogger{
		level:  level,
		logB:   log,
		prefix: prefix,
		output: ioutil.Discard,
	}
}

func (l *zapLogger) Prefix() string {
	return l.prefix
}

func (l *zapLogger) SetPrefix(p string) {
	l.prefix = p
}

func (l *zapLogger) Level() log.Lvl {
	return l.level
}

// it is not removed by compatibility but is not used
func (l *zapLogger) SetLevel(level log.Lvl) {
}

func (l *zapLogger) Output() io.Writer {
	return l.output
}

// it is not removed by compatibility but is not used
func (l *zapLogger) SetOutput(w io.Writer) {
}

// it is not removed by compatibility but is not used
func (l *zapLogger) SetHeader(h string) {
}

func (l *zapLogger) Print(i ...interface{}) {
	l.logB.Info(fmt.Sprint(i...))
}

func (l *zapLogger) Printf(format string, args ...interface{}) {
	l.logB.Sugar().Infof(format, args)
}

func (l *zapLogger) Printj(j log.JSON) {
	l.logB.Info(l.formatJSON(j))
}

func (l *zapLogger) Debug(i ...interface{}) {
	l.logB.Debug(fmt.Sprint(i...))
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logB.Sugar().Debugf(format, args)
}

func (l *zapLogger) Debugj(j log.JSON) {
	l.logB.Debug(l.formatJSON(j))
}

func (l *zapLogger) Info(i ...interface{}) {
	l.logB.Info(fmt.Sprint(i...))
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logB.Sugar().Infof(format, args)
}

func (l *zapLogger) Infoj(j log.JSON) {
	l.logB.Info(l.formatJSON(j))
}

func (l *zapLogger) Warn(i ...interface{}) {
	l.logB.Warn(fmt.Sprint(i...))
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logB.Sugar().Warnf(format, args)
}

func (l *zapLogger) Warnj(j log.JSON) {
	l.logB.Warn(l.formatJSON(j))
}

func (l *zapLogger) Error(i ...interface{}) {
	l.logB.Error(fmt.Sprint(i...))
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logB.Sugar().Errorf(format, args)
}

func (l *zapLogger) Errorj(j log.JSON) {
	l.logB.Error(l.formatJSON(j))
}

func (l *zapLogger) Fatal(i ...interface{}) {
	l.logB.Fatal(fmt.Sprint(i...))
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.logB.Sugar().Fatalf(format, args)
}

func (l *zapLogger) Fatalj(j log.JSON) {
	l.logB.Fatal(l.formatJSON(j))
}

func (l *zapLogger) Panic(i ...interface{}) {
	l.logB.Panic(fmt.Sprint(i...))
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.logB.Sugar().Panicf(format, args)
}

func (l *zapLogger) Panicj(j log.JSON) {
	l.logB.Panic(l.formatJSON(j))
}

func (l *zapLogger) formatJSON(j log.JSON) string {
	b, err := json.Marshal(j)
	if err != nil {
		l.logB.Error("Converting to JSON in echo log", zap.String("error", err.Error()))
		return ""
	}
	return string(b)
}

// ConvertLevel converts logging level to log lvl
func ConvertLevel(level logging.Level) log.Lvl {
	switch level {
	case logging.DebugLevel:
		return log.DEBUG
	case logging.InfoLevel:
		return log.INFO
	case logging.WarnLevel:
		return log.WARN
	case logging.ErrorLevel:
		return log.ERROR
	default:
		return log.OFF
	}
}
