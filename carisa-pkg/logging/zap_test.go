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

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

const message = "message"

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

		l.Info(message, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapWarn(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.WarnLevel)

		l.Warn(message, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapDebug(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		l.Debug(message, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapError(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.ErrorLevel)

		l.Error(message, convertTo(tt.items)...)
		check(t, recorded, tt)
	}
}

func TestZapWrapPanic(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.PanicLevel)

		assert.Panics(t, func() { l.Panic(message, convertTo(tt.items)...) })
		check(t, recorded, tt)
	}
}

func TestZapWrapCheck(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		if ce := l.Check(InfoLevel, message); ce != nil {
			l.Write(ce, convertTo(tt.items)...)
		}
		check(t, recorded, tt)
	}
}

func TestZapWrapOutOfRangeFields(t *testing.T) {
	_, l := newLogger(zapcore.InfoLevel)

	assert.Panics(t,
		func() {
			l.Info(message,
				Bool("k", true),
				Bool("k", true),
				Bool("k", true),
				Bool("k", true),
				Bool("k", true))
		})
}

func check(t *testing.T, recorded *observer.ObservedLogs, tt tests) {
	for _, logs := range recorded.All() {
		assert.Equal(t, message, logs.Message, "Message")
		lenF := len(tt.items) - 1
		for idx, f := range logs.Context {
			if idx > lenF {
				assert.Equal(t, zapcore.SkipType, f.Type, "Key")
				continue
			}
			assert.Equal(t, tt.items[idx].f.key, f.Key, "Key")
			assert.Equal(t, tt.items[idx].t, f.Type, "Type")
		}
	}
}

func newLogger(level zapcore.Level) (*observer.ObservedLogs, Logger) {
	core, obs := observer.New(level)
	return obs, &ZapWrap{zap.New(core)}
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
				{
					f: Bool("key3", true),
					t: zapcore.BoolType,
				},
			},
		},
	}
}
