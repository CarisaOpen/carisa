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
	"os"
	"os/exec"
	"testing"

	"github.com/carisa/pkg/logging"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const prefix = "prefix"

func TestZapLogger_Prefix(t *testing.T) {
	_, l := newZapLogger(zapcore.InfoLevel)
	assert.Equal(t, prefix, l.Prefix())
}

func TestZapLogger_SetPrefix(t *testing.T) {
	_, l := newZapLogger(zapcore.InfoLevel)
	const pre = "prefix1"
	l.SetPrefix(pre)
	assert.Equal(t, pre, l.Prefix())
}

func TestZapLogger_Level(t *testing.T) {
	_, l := newZapLogger(zapcore.InfoLevel)
	assert.Equal(t, log.INFO, l.Level())
}

func TestZapLogger_SetLevel(t *testing.T) {
	_, l := newZapLogger(zapcore.InfoLevel)
	l.SetLevel(log.ERROR)
	assert.Equal(t, log.INFO, l.Level())
}

func TestZapLogger_Output(t *testing.T) {
	_, l := newZapLogger(zapcore.InfoLevel)
	assert.NotNil(t, l.Output())
}

func TestZapLogger_SetHeader(t *testing.T) {
	_, l := newZapLogger(zapcore.InfoLevel)
	l.SetHeader("")
}

func TestZapLogger_Log(t *testing.T) {
	j := make(log.JSON)
	j["param1"] = "a"
	j["param2"] = 1

	const jRes = "{\"param1\":\"a\",\"param2\":1}"

	tests := []struct {
		name  string
		level zapcore.Level
		log   func(echo.Logger)
		res   string
	}{
		{
			name:  "Print",
			level: zap.InfoLevel,
			log:   func(l echo.Logger) { l.Print("a", 1) },
			res:   "a1",
		},
		{
			name:  "Print",
			level: zap.InfoLevel,
			log:   func(l echo.Logger) { l.Printf("Msg: %s", "par") },
			res:   "Msg: [par]",
		},
		{
			name:  "Print",
			level: zap.InfoLevel,
			log:   func(l echo.Logger) { l.Printj(j) },
			res:   jRes,
		},
		{
			name:  "Debug",
			level: zap.DebugLevel,
			log:   func(l echo.Logger) { l.Debug("a", 1) },
			res:   "a1",
		},
		{
			name:  "Debug",
			level: zap.DebugLevel,
			log:   func(l echo.Logger) { l.Debugf("Msg: %s", "par") },
			res:   "Msg: [par]",
		},
		{
			name:  "Debug",
			level: zap.DebugLevel,
			log:   func(l echo.Logger) { l.Debugj(j) },
			res:   jRes,
		},
		{
			name:  "Info",
			level: zap.InfoLevel,
			log:   func(l echo.Logger) { l.Info("a", 1) },
			res:   "a1",
		},
		{
			name:  "Info",
			level: zap.InfoLevel,
			log:   func(l echo.Logger) { l.Infof("Msg: %s", "par") },
			res:   "Msg: [par]",
		},
		{
			name:  "Info",
			level: zap.InfoLevel,
			log:   func(l echo.Logger) { l.Infoj(j) },
			res:   jRes,
		},
		{
			name:  "Warn",
			level: zap.WarnLevel,
			log:   func(l echo.Logger) { l.Warn("a", 1) },
			res:   "a1",
		},
		{
			name:  "Warn",
			level: zap.WarnLevel,
			log:   func(l echo.Logger) { l.Warnf("Msg: %s", "par") },
			res:   "Msg: [par]",
		},
		{
			name:  "Warn",
			level: zap.WarnLevel,
			log:   func(l echo.Logger) { l.Warnj(j) },
			res:   jRes,
		},
		{
			name:  "Error",
			level: zap.ErrorLevel,
			log:   func(l echo.Logger) { l.Error("a", 1) },
			res:   "a1",
		},
		{
			name:  "Error",
			level: zap.ErrorLevel,
			log:   func(l echo.Logger) { l.Errorf("Msg: %s", "par") },
			res:   "Msg: [par]",
		},
		{
			name:  "Error",
			level: zap.ErrorLevel,
			log:   func(l echo.Logger) { l.Errorj(j) },
			res:   jRes,
		},
		{
			name:  "Panic",
			level: zap.PanicLevel,
			log:   func(l echo.Logger) { l.Panic("a", 1) },
			res:   "a1",
		},
		{
			name:  "Panic",
			level: zap.PanicLevel,
			log:   func(l echo.Logger) { l.Panicf("Msg: %s", "par") },
			res:   "Msg: [par]",
		},
		{
			name:  "Panic",
			level: zap.PanicLevel,
			log:   func(l echo.Logger) { l.Panicj(j) },
			res:   jRes,
		},
	}
	for _, tt := range tests {
		recorded, l := newZapLogger(tt.level)
		if tt.level == zap.PanicLevel {
			assert.Panics(t, func() { tt.log(l) })
		} else {
			tt.log(l)
		}
		for _, logs := range recorded.All() {
			assert.Equal(t, tt.res, logs.Message, tt.name)
		}
	}
}

func TestZapLogger_Fatal(t *testing.T) {
	j := make(log.JSON)
	j["param1"] = "a"

	tests := []struct {
		log func(echo.Logger)
	}{
		{
			log: func(l echo.Logger) { l.Fatal("a", 1) },
		},
		{
			log: func(l echo.Logger) { l.Fatalf("Msg: %s", "par") },
		},
		{
			log: func(l echo.Logger) { l.Errorj(j) },
		},
	}

	for _, tt := range tests {
		if os.Getenv("FLAG") == "1" {
			_, l := newZapLogger(zap.FatalLevel)
			tt.log(l)
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestZapLogger_Fatal") //nolint:gosec
		cmd.Env = append(os.Environ(), "FLAG=1")
		err := cmd.Run()

		e, ok := err.(*exec.ExitError)
		expectedErrorString := "exit status 1"
		assert.True(t, ok)
		assert.Equal(t, expectedErrorString, e.Error())
	}
}

func TestZapLogger_FormatJSONError(t *testing.T) {
	j := make(log.JSON)
	j["param"] = func() {}

	recorded, l := newZapLogger(zapcore.ErrorLevel)
	r := l.(*zapLogger).formatJSON(j)

	assert.Empty(t, r, "formatJSON result")
	for _, logs := range recorded.All() {
		assert.Equal(t, "Converting to JSON in echo log", logs.Message, "Error log")
	}
}

func TestZapLogger_ConvertLevel(t *testing.T) {
	tests := []struct {
		ls logging.Level
		lt log.Lvl
	}{
		{
			ls: logging.DebugLevel,
			lt: log.DEBUG,
		},
		{
			ls: logging.InfoLevel,
			lt: log.INFO,
		},
		{
			ls: logging.WarnLevel,
			lt: log.WARN,
		},
		{
			ls: logging.ErrorLevel,
			lt: log.ERROR,
		},
		{
			ls: logging.PanicLevel,
			lt: log.OFF,
		},
	}
	for _, tt := range tests {
		r := ConvertLevel(tt.ls)
		assert.Equal(t, tt.lt, r)
	}
}

func newZapLogger(level zapcore.Level) (*observer.ObservedLogs, echo.Logger) {
	core, obs := observer.New(level)
	return obs, NewLogging(prefix, log.INFO, zap.New(core))
}
