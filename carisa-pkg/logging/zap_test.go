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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	loc     = "test"
	message = "message"
)

type Item struct {
	t zapcore.FieldType
	f Field
}

type tests struct {
	items []Item
}

func TestZapWrap_NewZapLogger(t *testing.T) {
	tests := []struct {
		l      ZapConfig
		levelR zapcore.Level
		levelC Level
	}{
		{
			l: ZapConfig{
				Development: false,
				Level:       0,
				Encoding:    "",
			},
			levelC: InfoLevel,
			levelR: zap.InfoLevel,
		},
		{
			l: ZapConfig{
				Development: false,
				Level:       InfoLevel,
				Encoding:    "json",
			},
			levelC: InfoLevel,
			levelR: zap.InfoLevel,
		},
		{
			l: ZapConfig{
				Development: true,
				Level:       0,
				Encoding:    "console",
			},
			levelC: DebugLevel,
			levelR: zap.InfoLevel,
		},
		{
			l: ZapConfig{
				Development: true,
				Level:       PanicLevel,
				Encoding:    "json",
			},
			levelC: PanicLevel,
			levelR: zap.PanicLevel,
		},
	}
	for _, tt := range tests {
		l, zL := NewZapLogger(tt.l)
		assert.Equal(t, tt.levelC, l.Level())
		assert.True(t, zL.Core().Enabled(tt.levelR), "Level")
	}
}

func TestZapWrap_NewDev(t *testing.T) {
	log, err := NewZapWrapDev()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, DebugLevel, log.Level(), "Level")
	assert.NotNil(t, log.(*zapWrap).log, "Zap object")
}

func TestZapWrap_Info(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		switch i {
		case 0:
			l.Info(message, loc)
		case 1:
			l.Info1(message, loc, tt.items[0].f)
		case 2:
			l.Info2(message, loc, tt.items[0].f, tt.items[1].f)
		case 3:
			l.Info3(message, loc, tt.items[0].f, tt.items[1].f, tt.items[2].f)
		}

		check(t, recorded, tt)
	}
}

func TestZapWrap_Warn(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.WarnLevel)

		switch i {
		case 0:
			l.Warn(message, loc)
		case 1:
			l.Warn1(message, loc, tt.items[0].f)
		case 2:
			l.Warn2(message, loc, tt.items[0].f, tt.items[1].f)
		case 3:
			l.Warn3(message, loc, tt.items[0].f, tt.items[1].f, tt.items[2].f)
		}

		check(t, recorded, tt)
	}
}

func TestZapWrap_Debug(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		switch i {
		case 0:
			l.Debug(message, loc)
		case 1:
			l.Debug1(message, loc, tt.items[0].f)
		case 2:
			l.Debug2(message, loc, tt.items[0].f, tt.items[1].f)
		case 3:
			l.Debug3(message, loc, tt.items[0].f, tt.items[1].f, tt.items[2].f)
		}

		check(t, recorded, tt)
	}
}

func TestZapWrap_Error(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		switch i {
		case 0:
			l.Error(message, loc)
		case 1:
			l.Error1(message, loc, tt.items[0].f)
		case 2:
			l.Error2(message, loc, tt.items[0].f, tt.items[1].f)
		case 3:
			l.Error3(message, loc, tt.items[0].f, tt.items[1].f, tt.items[2].f)
		}

		check(t, recorded, tt)
	}
}

func TestZapWrap_ErrorE(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		switch i {
		case 0:
			l.ErrorE(errors.New(message), loc)
		case 1:
			l.ErrorE1(errors.New(message), loc, tt.items[0].f)
		case 2:
			l.ErrorE2(errors.New(message), loc, tt.items[0].f, tt.items[1].f)
		case 3:
			l.ErrorE3(errors.New(message), loc, tt.items[0].f, tt.items[1].f, tt.items[2].f)
		}

		check(t, recorded, tt)
	}
}

func TestZapWrap_Panic(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		switch i {
		case 0:
			assert.Panics(t, func() { l.Panic(message, loc) })
		case 1:
			assert.Panics(t, func() { l.Panic1(message, loc, tt.items[0].f) })
		case 2:
			assert.Panics(t, func() { l.Panic2(message, loc, tt.items[0].f, tt.items[1].f) })
		case 3:
			assert.Panics(t, func() { l.Panic3(message, loc, tt.items[0].f, tt.items[1].f, tt.items[2].f) })
		}

		check(t, recorded, tt)
	}
}

func TestZapWrap_PanicE(t *testing.T) {
	_, l := newLogger(zapcore.PanicLevel)

	assert.Panics(t, func() { l.PanicE(errors.New("panic"), "loc") })
}

func TestZapWrap_Check(t *testing.T) {
	tests := testdd()

	for i, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		if ce := l.Check(InfoLevel, message); ce != nil {
			switch i {
			case 0:
				l.Write(ce, loc)
			case 1:
				l.Write1(ce, loc, tt.items[0].f)
			case 2:
				l.Write2(ce, loc, tt.items[0].f, tt.items[1].f)
			case 3:
				l.Write3(ce, loc, tt.items[0].f, tt.items[1].f, tt.items[2].f)
			}
		}
		check(t, recorded, tt)
	}
}

func TestZapWrap_CheckNoMatchLevel(t *testing.T) {
	_, l := newLogger(zapcore.PanicLevel)
	ce := l.Check(InfoLevel, "message")
	assert.Nil(t, ce)
}

func TestZapWrap_ErrWrap(t *testing.T) {
	tests := []struct {
		fs []Field
		m  string
	}{
		{
			fs: []Field{},
			m:  "message: error",
		},
		{
			fs: []Field{
				String("key", "value"),
			},
			m: "message. key: value: error",
		},
		{
			fs: []Field{
				String("key", "value"),
				Bool("key1", true),
			},
			m: "message. key: value, key1: true: error",
		},
		{
			fs: []Field{
				String("key", "value"),
				Bool("key1", true),
				String("key2", "value2"),
			},
			m: "message. key: value, key1: true, key2: value2: error",
		},
	}

	_, l := newLogger(zapcore.InfoLevel)

	for i, tt := range tests {
		switch i {
		case 0:
			if err := l.ErrWrap(errors.New("error"), "message", "test"); err != nil {
				assert.Equal(t, tt.m, err.Error())
				return
			}
		case 1:
			if err := l.ErrWrap1(errors.New("error"), "message", "test", tt.fs[0]); err != nil {
				assert.Equal(t, tt.m, err.Error())
				return
			}
		case 2:
			if err := l.ErrWrap2(errors.New("error"), "message", "test", tt.fs[0], tt.fs[1]); err != nil {
				assert.Equal(t, tt.m, err.Error())
				return
			}
		case 3:
			if err := l.ErrWrap3(errors.New("error"), "message", "test", tt.fs[0], tt.fs[1], tt.fs[2]); err != nil {
				assert.Equal(t, tt.m, err.Error())
				return
			}
		}
		t.Error("err cannot be nil")
	}
}

func TestZapWrap_ConvertZapLevel(t *testing.T) {
	tests := []struct {
		l  Level
		zL zapcore.Level
	}{
		{
			l:  DebugLevel,
			zL: zap.DebugLevel,
		},
		{
			l:  InfoLevel,
			zL: zap.InfoLevel,
		},
		{
			l:  WarnLevel,
			zL: zap.WarnLevel,
		},
		{
			l:  ErrorLevel,
			zL: zap.ErrorLevel,
		},
		{
			l:  PanicLevel,
			zL: zap.PanicLevel,
		},
	}
	for _, tt := range tests {
		r := ConvertZapLevel(tt.l)
		assert.Equal(t, tt.zL, r)
	}
}

func check(t *testing.T, recorded *observer.ObservedLogs, tt tests) {
	for _, logs := range recorded.All() {
		assert.Equal(t, message, logs.Message, "Message")
		for idx, f := range logs.Context {
			if idx == 0 {
				assert.Equal(t, "location", f.Key, "location")
				assert.Equal(t, loc, f.String, "location")
				continue
			}
			assert.Equal(t, tt.items[idx-1].f.key, f.Key, "Key")
			assert.Equal(t, tt.items[idx-1].t, f.Type, "Type")
		}
	}
}

func newLogger(level zapcore.Level) (*observer.ObservedLogs, Logger) {
	core, obs := observer.New(level)
	return obs, NewZapWrap(zap.New(core), DebugLevel, "")
}

func convertTo(items []Item) []Field {
	fields := make([]Field, len(items))
	for i, item := range items {
		fields[i] = item.f
	}
	return fields
}

func testdd() []tests {
	return []tests{
		{
			items: []Item{},
		},
		{
			items: []Item{
				{
					f: Bool("key", true),
					t: zapcore.BoolType,
				},
			},
		},
		{
			items: []Item{
				{
					f: Bool("key", true),
					t: zapcore.BoolType,
				},
				{
					f: String("key1", "value"),
					t: zapcore.StringType,
				},
			},
		},
		{
			items: []Item{
				{
					f: Bool("key", true),
					t: zapcore.BoolType,
				},
				{
					f: String("key1", "value"),
					t: zapcore.StringType,
				},
				{
					f: Bool("key2", true),
					t: zapcore.BoolType,
				},
			},
		},
	}
}
