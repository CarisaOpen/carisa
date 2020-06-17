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
	"github.com/carisa/pkg/strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const fieldsSize = 4

// zapWrap is a zap wrapper.
// it is defined a field size fixed to not escape to heap
type zapWrap struct {
	log *zap.Logger
	loggerComp
	locName string
}

// NewZapWrap creates zapWrap. If loc parameter is empty, loc is configured to "location"
func NewZapWrap(log *zap.Logger, loc string) Logger {
	if len(loc) == 0 {
		loc = "location"
	}
	zq := &zapWrap{
		log:     log,
		locName: loc,
	}
	zq.loggerComp.log = zq
	return zq
}

// Info implements logging.Logger.Info
func (z *zapWrap) Info(msg string, loc string, fields ...Field) {
	fSource := z.convertToZap(loc, fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Info(msg, fTarget...)
}

// Warn implements logging.Logger.Warn
func (z *zapWrap) Warn(msg string, loc string, fields ...Field) {
	fSource := z.convertToZap(loc, fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Debug(msg, fTarget...)
}

// Debug implements logging.Logger.Debug
func (z *zapWrap) Debug(msg string, loc string, fields ...Field) {
	fSource := z.convertToZap(loc, fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Debug(msg, fTarget...)
}

// Error implements logging.Logger.Error
func (z *zapWrap) Error(msg string, loc string, fields ...Field) {
	fSource := z.convertToZap(loc, fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Error(msg, fTarget...)
}

// Panic implements logging.Logger.Panic
func (z *zapWrap) Panic(msg string, loc string, fields ...Field) {
	fSource := z.convertToZap(loc, fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Panic(msg, fTarget...)
}

// Check implements logging.Logger.Check
func (z *zapWrap) Check(l Level, msg string) *CheckWrap {
	if ce := z.log.Check(zapcore.Level(l), msg); ce != nil {
		return &CheckWrap{ce}
	}
	return nil
}

// Write implements logging.Logger.Write
func (z *zapWrap) Write(wrap *CheckWrap, loc string, fields ...Field) {
	fSource := z.convertToZap(loc, fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	wrap.zapCheck.Write(fTarget...)
}

// The size of the array is fixed so that it does not escape to heap
func (z *zapWrap) convertToZap(loc string, fields ...Field) [fieldsSize]zap.Field {
	var fZap [fieldsSize]zap.Field
	actualFSize := len(fields) + 1 // loc is added
	if actualFSize > fieldsSize {
		panic(strings.Concat("Log fields and location size cannot be more than ", string(fieldsSize)))
	}

	fZap[0] = zap.String(z.locName, loc)
	for i := 1; i < fieldsSize; i++ {
		if i >= actualFSize { // The rest are skip until fill buffer
			fZap[i] = zap.Skip()
			continue
		}
		f := fields[i-1]
		switch f.tpy {
		case stringType:
			fZap[i] = zap.String(f.key, f.stringV)
		case boolType:
			fZap[i] = zap.Bool(f.key, f.boolV)
		}
	}
	return fZap
}
