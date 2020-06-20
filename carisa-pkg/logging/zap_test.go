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

func TestZapWrapInfo(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		l.Info(message, loc, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapWarn(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.WarnLevel)

		l.Warn(message, loc, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapDebug(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		l.Debug(message, loc, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapError(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.ErrorLevel)

		l.Error(message, loc, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapErrorE(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.ErrorLevel)

		l.ErrorE(errors.New(message), loc, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapPanic(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.PanicLevel)

		assert.Panics(t, func() { l.Panic(message, loc, convertTo(tt.items)...) })
		check(t, recorded, tt)
	}
}

func TestZapWrapCheck(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		if ce := l.Check(InfoLevel, message); ce != nil {
			l.Write(ce, loc, convertTo(tt.items)...)
		}
		check(t, recorded, tt)
	}
}

func TestZapWrapCheckNoMatchLevel(t *testing.T) {
	_, l := newLogger(zapcore.PanicLevel)
	ce := l.Check(InfoLevel, "message")
	assert.Nil(t, ce)
}

func TestZapWrapOutOfRangeFields(t *testing.T) {
	_, l := newLogger(zapcore.InfoLevel)

	assert.Panics(t,
		func() {
			l.Info(message, loc,
				Bool("k", true),
				Bool("k", true),
				Bool("k", true),
				Bool("k", true))
		})
}

func TestErrWrap(t *testing.T) {
	_, l := newLogger(zapcore.InfoLevel)

	if err := l.ErrWrap(errors.New("error"), "message", "test", String("key", "value")); err != nil {
		assert.Equal(t, "message. key: value: error", err.Error())
		return
	}
	t.Error("err cannot be nil")
}

func TestConvertZapLevel(t *testing.T) {
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

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		l      ZapConfig
		panic  bool
		levelR zapcore.Level
	}{
		{
			l: ZapConfig{
				Development: false,
				Level:       0,
				Encoding:    "",
			},
			panic:  false,
			levelR: zap.InfoLevel,
		},
		{
			l: ZapConfig{
				Development: false,
				Level:       InfoLevel,
				Encoding:    "json",
			},
			panic:  false,
			levelR: zap.InfoLevel,
		},
		{
			l: ZapConfig{
				Development: true,
				Level:       0,
				Encoding:    "console",
			},
			panic:  false,
			levelR: zap.InfoLevel,
		},
		{
			l: ZapConfig{
				Development: true,
				Level:       PanicLevel,
				Encoding:    "json",
			},
			panic:  true,
			levelR: zap.PanicLevel,
		},
	}
	for _, tt := range tests {
		if tt.panic {
			assert.Panics(t, func() { NewZapLogger(tt.l) })
		} else {
			zL := NewZapLogger(tt.l).(*zapWrap)
			assert.True(t, zL.log.Core().Enabled(tt.levelR), "Level")
		}
	}
}

func check(t *testing.T, recorded *observer.ObservedLogs, tt tests) {
	for _, logs := range recorded.All() {
		assert.Equal(t, message, logs.Message, "Message")
		lenF := len(tt.items)
		for idx, f := range logs.Context {
			if idx == 0 {
				assert.Equal(t, "location", f.Key, "location")
				assert.Equal(t, loc, f.String, "location")
				continue
			}
			if idx > lenF {
				assert.Equal(t, zapcore.SkipType, f.Type, "Key")
				continue
			}
			assert.Equal(t, tt.items[idx-1].f.key, f.Key, "Key")
			assert.Equal(t, tt.items[idx-1].t, f.Type, "Type")
		}
	}
}

func newLogger(level zapcore.Level) (*observer.ObservedLogs, Logger) {
	core, obs := observer.New(level)
	return obs, NewZapWrap(zap.New(core), "")
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
				{
					f: Bool("key2", true),
					t: zapcore.BoolType,
				},
			},
		},
	}
}
