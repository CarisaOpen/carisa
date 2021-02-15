/*
 *  Copyright 2019-2022 the original author or authors.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations under the License.
 *
 */

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicks_Renew(t *testing.T) {
	ti := newTicks()
	ts := ti.timeStamp
	ti.renew()
	assert.NotEqual(t, ts, ti.timeStamp)
	assert.Equal(t, ts, ti.previousTimeStamp)
}

func TestTicks_undo(t *testing.T) {
	ti := newTicks()
	ts := ti.timeStamp
	ti.renew()
	ti.undo()
	assert.Equal(t, ts, ti.timeStamp)
}

func TestTicks_String(t *testing.T) {
	tests := []struct {
		name string
		ts   time.Time
		res  string
	}{
		{
			name: "Ticks - 20210101010000",
			ts:   time.Date(2021, time.January, 1, 1, 0, 0, 0, time.UTC),
			res:  "20210101010000",
		},
		{
			name: "Ticks - 20211115225959",
			ts:   time.Date(2021, time.November, 15, 22, 59, 59, 0, time.UTC),
			res:  "20211115225959",
		},
	}

	for _, tt := range tests {
		ti := newTicks()
		ti.timeStamp = tt.ts
		assert.Equal(t, tt.res, ti.tstring())
	}
}
