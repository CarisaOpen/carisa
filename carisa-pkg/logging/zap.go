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

// If the len is modified, change info, error, etc.
const fieldsSize = 4

// ZapConfig defines the configuration for log framework
type ZapConfig struct {
	// Development mode. Default value: false
	Development bool `json:"development,omitempty"`
	// Level. See logging.Level. Default value: Depending of Development flag
	Level Level `json:"level,omitempty"`
	// Encoding type. Default value: Depending of Development flag
	// The values can be: j -> json format, c -> console format
	Encoding string `json:"encoding,omitempty"`
}

// zapWrap is a zap wrapper.
// Remark: it is defined a fields size fixed to not escape to heap
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

// Level implements logging.Logger.Level
func (z *zapWrap) Level() Level {
	return z.level
}

// Info implements logging.Logger.Info
func (z *zapWrap) Info(msg string, loc string, fields ...Field) {
	f := z.convertToZap(loc, fields...)
	z.log.Info(msg, f[0], f[1], f[2], f[3])
}

// Warn implements logging.Logger.Warn
func (z *zapWrap) Warn(msg string, loc string, fields ...Field) {
	f := z.convertToZap(loc, fields...)
	z.log.Debug(msg, f[0], f[1], f[2], f[3])
}

// Debug implements logging.Logger.Debug
func (z *zapWrap) Debug(msg string, loc string, fields ...Field) {
	f := z.convertToZap(loc, fields...)
	z.log.Debug(msg, f[0], f[1], f[2], f[3])
}

// Error implements logging.Logger.Error
func (z *zapWrap) Error(msg string, loc string, fields ...Field) {
	f := z.convertToZap(loc, fields...)
	z.log.Error(msg, f[0], f[1], f[2], f[3])
}

// Panic implements logging.Logger.Panic
func (z *zapWrap) Panic(msg string, loc string, fields ...Field) {
	f := z.convertToZap(loc, fields...)
	z.log.Panic(msg, f[0], f[1], f[2], f[3])
}

// Check implements logging.Logger.Check
func (z *zapWrap) Check(l Level, msg string) *CheckWrap {
	if ce := z.log.Check(ConvertZapLevel(l), msg); ce != nil {
		return &CheckWrap{ce}
	}
	return nil
}

// Write implements logging.Logger.Write
func (z *zapWrap) Write(wrap *CheckWrap, loc string, fields ...Field) {
	f := z.convertToZap(loc, fields...)
	wrap.zapCheck.Write(f[0], f[1], f[2], f[3])
}

// The size of the array is fixed so that it does not escape to heap
func (z *zapWrap) convertToZap(loc string, fields ...Field) [fieldsSize]zap.Field {
	var fZap [fieldsSize]zap.Field
	actualFSize := len(fields) + 1 // loc is added
	if actualFSize > fieldsSize {
		panic(strings.Concat("log fields and location size cannot be more than ", string(fieldsSize)))
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
