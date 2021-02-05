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

// ZapConfig defines the configuration for log framework
type ZapConfig struct {
	// Development mode. Common value: false
	Development bool `json:"development,omitempty"`
	// Level. See logging.Level. Common value: Depending of Development flag
	Level Level `json:"level,omitempty"`
	// Encoding type. Common value: Depending of Development flag
	// The values can be: j -> json format, c -> console format
	Encoding string `json:"encoding,omitempty"`
}

// zapWrap is a zap wrapper.
// Remark: Variadic params escape to heap, therefore the log methods is overloaded
type zapWrap struct {
	log   *zap.Logger
	level Level
	loggerComp
	locName string
}

// NewZapLogger builds a zap logger from config
func NewZapLogger(config ZapConfig) (Logger, *zap.Logger) {
	var log zap.Config
	var level Level
	if config.Development {
		log = zap.NewDevelopmentConfig()
		level = DebugLevel
	} else {
		log = zap.NewProductionConfig()
		level = InfoLevel
	}
	if config.Level > 0 {
		log.Level = zap.NewAtomicLevelAt(ConvertZapLevel(config.Level))
		level = config.Level
	}
	if len(config.Encoding) > 0 {
		log.Encoding = config.Encoding
	}

	l, err := log.Build()
	if err != nil {
		panic("Error creating zap logger")
	}
	return NewZapWrap(l, level, ""), l
}

// NewZapWrapDev creates ZapWrap for development
func NewZapWrapDev() (Logger, error) {
	log, err := zap.NewDevelopment()
	return NewZapWrap(log, DebugLevel, ""), err
}

// NewZapWrap creates zapWrap. If loc parameter is empty, loc is configured to "location"
func NewZapWrap(log *zap.Logger, level Level, loc string) Logger {
	if len(loc) == 0 {
		loc = "location"
	}
	zq := &zapWrap{
		log:     log,
		level:   level,
		locName: loc,
	}
	zq.loggerComp.log = zq
	return zq
}

// Level implements Logger.Level
func (z *zapWrap) Level() Level {
	return z.level
}

// Info implements Logger.Info
func (z *zapWrap) Info(msg string, loc string) {
	z.log.Info(msg, z.locZap(loc))
}

// Info1 implements Logger.Info1
func (z *zapWrap) Info1(msg string, loc string, f Field) {
	z.log.Info(msg, z.locZap(loc), convertToZap(f))
}

// Info2 implements Logger.Info2
func (z *zapWrap) Info2(msg string, loc string, f Field, f1 Field) {
	z.log.Info(msg, z.locZap(loc), convertToZap(f), convertToZap(f1))
}

// Info3 implements Logger.Info3
func (z *zapWrap) Info3(msg string, loc string, f Field, f1 Field, f2 Field) {
	z.log.Info(msg, z.locZap(loc), convertToZap(f), convertToZap(f1), convertToZap(f2))
}

// Warn implements Logger.Warn
func (z *zapWrap) Warn(msg string, loc string) {
	z.log.Warn(msg, z.locZap(loc))
}

// Warn1 implements Logger.Warn1
func (z *zapWrap) Warn1(msg string, loc string, f Field) {
	z.log.Warn(msg, z.locZap(loc), convertToZap(f))
}

// Warn2 implements Logger.Warn2
func (z *zapWrap) Warn2(msg string, loc string, f Field, f1 Field) {
	z.log.Warn(msg, z.locZap(loc), convertToZap(f), convertToZap(f1))
}

// Warn3 implements Logger.Warn3
func (z *zapWrap) Warn3(msg string, loc string, f Field, f1 Field, f2 Field) {
	z.log.Warn(msg, z.locZap(loc), convertToZap(f), convertToZap(f1), convertToZap(f2))
}

// Debug implements Logger.Debug
func (z *zapWrap) Debug(msg string, loc string) {
	z.log.Debug(msg, z.locZap(loc))
}

// Debug1 implements Logger.Debug1
func (z *zapWrap) Debug1(msg string, loc string, f Field) {
	z.log.Debug(msg, z.locZap(loc), convertToZap(f))
}

// Debug2 implements Logger.Debug2
func (z *zapWrap) Debug2(msg string, loc string, f Field, f1 Field) {
	z.log.Debug(msg, z.locZap(loc), convertToZap(f), convertToZap(f1))
}

// Debug3 implements Logger.Debug3
func (z *zapWrap) Debug3(msg string, loc string, f Field, f1 Field, f2 Field) {
	z.log.Debug(msg, z.locZap(loc), convertToZap(f), convertToZap(f1), convertToZap(f2))
}

// Error implements Logger.Error
func (z *zapWrap) Error(msg string, loc string) {
	z.log.Error(msg, z.locZap(loc))
}

// Error1 implements Logger.Error1
func (z *zapWrap) Error1(msg string, loc string, f Field) {
	z.log.Error(msg, z.locZap(loc), convertToZap(f))
}

// Error2 implements Logger.Error2
func (z *zapWrap) Error2(msg string, loc string, f Field, f1 Field) {
	z.log.Error(msg, z.locZap(loc), convertToZap(f), convertToZap(f1))
}

// Error3 implements Logger.Error3
func (z *zapWrap) Error3(msg string, loc string, f Field, f1 Field, f2 Field) {
	z.log.Error(msg, z.locZap(loc), convertToZap(f), convertToZap(f1), convertToZap(f2))
}

// Panic implements Logger.Panic
func (z *zapWrap) Panic(msg string, loc string) {
	z.log.Panic(msg, z.locZap(loc))
}

// Panic1 implements Logger.Panic1
func (z *zapWrap) Panic1(msg string, loc string, f Field) {
	z.log.Panic(msg, z.locZap(loc), convertToZap(f))
}

// Panic2 implements Logger.Panic2
func (z *zapWrap) Panic2(msg string, loc string, f Field, f1 Field) {
	z.log.Panic(msg, z.locZap(loc), convertToZap(f), convertToZap(f1))
}

// Panic3 implements Logger.Panic3
func (z *zapWrap) Panic3(msg string, loc string, f Field, f1 Field, f2 Field) {
	z.log.Panic(msg, z.locZap(loc), convertToZap(f), convertToZap(f1), convertToZap(f2))
}

// Check implements Logger.Check
func (z *zapWrap) Check(l Level, msg string) *CheckWrap {
	if ce := z.log.Check(ConvertZapLevel(l), msg); ce != nil {
		return &CheckWrap{ce}
	}
	return nil
}

// Write implements Logger.Write
func (z *zapWrap) Write(wrap *CheckWrap, loc string) {
	wrap.zapCheck.Write(z.locZap(loc))
}

// Write1 implements Logger.Write1
func (z *zapWrap) Write1(wrap *CheckWrap, loc string, f Field) {
	wrap.zapCheck.Write(z.locZap(loc), convertToZap(f))
}

// Write2 implements Logger.Write2
func (z *zapWrap) Write2(wrap *CheckWrap, loc string, f Field, f1 Field) {
	wrap.zapCheck.Write(z.locZap(loc), convertToZap(f), convertToZap(f1))
}

// Write3 implements Logger.Write3
func (z *zapWrap) Write3(wrap *CheckWrap, loc string, f Field, f1 Field, f2 Field) {
	wrap.zapCheck.Write(z.locZap(loc), convertToZap(f), convertToZap(f1), convertToZap(f2))
}

func (z *zapWrap) locZap(loc string) zap.Field {
	return zap.String(z.locName, loc)
}

func convertToZap(f Field) zap.Field {
	switch f.tpy {
	case stringType:
		return zap.String(f.key, f.stringV)
	case boolType:
		return zap.Bool(f.key, f.boolV)
	default:
		panic("can not convert to zap")
	}
}

// ConvertZapLevel convert level to zap level
func ConvertZapLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case PanicLevel:
		return zap.PanicLevel
	}
	panic("logging level is wrong")
}
