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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const fieldsSize = 4

// zapWrap is a zap wrapper.
// I define a field size fixed to not escape to heap
type ZapWrap struct {
	log *zap.Logger
}

// Info implements logging.Logger.Info
func (z *ZapWrap) Info(msg string, fields ...Field) {
	fSource := convertToZap(fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Info(msg, fTarget...)
}

// Warn implements logging.Logger.Warn
func (z *ZapWrap) Warn(msg string, fields ...Field) {
	fSource := convertToZap(fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Debug(msg, fTarget...)
}

// Debug implements logging.Logger.Debug
func (z *ZapWrap) Debug(msg string, fields ...Field) {
	fSource := convertToZap(fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Debug(msg, fTarget...)
}

// Error implements logging.Logger.Error
func (z *ZapWrap) Error(msg string, fields ...Field) {
	fSource := convertToZap(fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Error(msg, fTarget...)
}

// Panic implements logging.Logger.Panic
func (z *ZapWrap) Panic(msg string, fields ...Field) {
	fSource := convertToZap(fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	z.log.Panic(msg, fTarget...)
}

// Check implements logging.Logger.Check
func (z *ZapWrap) Check(l Level, msg string) *CheckWrap {
	if ce := z.log.Check(zapcore.Level(l), msg); ce != nil {
		return &CheckWrap{ce}
	}
	return nil
}

func (z *ZapWrap) Write(wrap *CheckWrap, fields ...Field) {
	fSource := convertToZap(fields...)
	fTarget := make([]zap.Field, fieldsSize)
	for i := 0; i < fieldsSize; i++ {
		fTarget[i] = fSource[i]
	}
	wrap.zapCheck.Write(fTarget...)
}

// The size of the array is fixed so that it does not escape to heap
func convertToZap(fields ...Field) [fieldsSize]zap.Field {
	var fZap [fieldsSize]zap.Field
	actualFSize := len(fields)
	if actualFSize > fieldsSize {
		panic("Log fields size cannot be more than 5")
	}
	for i := 0; i < fieldsSize; i++ {
		if i >= actualFSize { // The rest are skip until fill buffer
			fZap[i] = zap.Skip()
			continue
		}
		f := fields[i]
		switch f.tpy {
		case stringType:
			fZap[i] = zap.String(f.key, f.stringV)
		case boolType:
			fZap[i] = zap.Bool(f.key, f.boolV)
		}
	}
	return fZap
}
