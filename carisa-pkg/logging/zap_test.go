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

type tests struct {
	fields []Field
}

func TestZapWrapInfo(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		l.Info(message, tt.fields...)
		check(t, recorded, tt)
	}
}

func TestZapWrapWarn(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.WarnLevel)

		l.Warn(message, tt.fields...)
		check(t, recorded, tt)
	}
}

func TestZapWrapDebug(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.DebugLevel)

		l.Debug(message, tt.fields...)
		check(t, recorded, tt)
	}
}

func TestZapWrapError(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.ErrorLevel)

		l.Error(message, tt.fields...)
		check(t, recorded, tt)
	}
}

func TestZapWrapPanic(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.PanicLevel)

		assert.Panics(t, func() { l.Panic(message, tt.fields...) })
		check(t, recorded, tt)
	}
}

func TestZapWrapCheck(t *testing.T) {
	tests := testdd()

	for _, tt := range tests {
		recorded, l := newLogger(zapcore.InfoLevel)

		l.Check(InfoLevel, message, tt.fields...)
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
		lenF := len(tt.fields) - 1
		for idx, f := range logs.Context {
			if idx > lenF {
				assert.Equal(t, zapcore.SkipType, f.Type, "Key")
				continue
			}
			assert.Equal(t, tt.fields[idx].key, f.Key, "Key")
		}
	}
}

func newLogger(level zapcore.Level) (*observer.ObservedLogs, Logger) {
	core, obs := observer.New(level)
	return obs, &ZapWrap{zap.New(core)}
}

func testdd() []tests {
	return []tests{
		{
			fields: []Field{
				Bool("key", true),
			},
		},
		{
			fields: []Field{
				Bool("key", true),
				Bool("key1", true),
				Bool("key2", true),
				Bool("key3", true),
			},
		},
	}
}
