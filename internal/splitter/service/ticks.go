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
	"strconv"
	"time"

	"github.com/carisa/pkg/strings"
)

const ticksSchema = "ST"

// ticks defines a timestamp to let the controller know if the splitter service is dead
// This timestamp by splitter must be updated each n time
type ticks struct {
	timeStamp         time.Time
	previousTimeStamp time.Time
}

func newTicks() ticks {
	return ticks{
		timeStamp: time.Now(),
	}
}

// renew renews the timestamp with the actual now
func (t *ticks) renew() {
	t.previousTimeStamp = t.timeStamp
	t.timeStamp = time.Now()
}

// undo returns the previous timestamp
func (t *ticks) undo() {
	t.timeStamp = t.previousTimeStamp
}

// tstring converts the timestamp to string
func (t *ticks) tstring() string {
	return strings.Concat(
		strconv.Itoa(t.timeStamp.Year()),
		strings.Lpad(strconv.Itoa(int(t.timeStamp.Month())), 2, "0"),
		strings.Lpad(strconv.Itoa(t.timeStamp.Day()), 2, "0"),
		strings.Lpad(strconv.Itoa(t.timeStamp.Hour()), 2, "0"),
		strings.Lpad(strconv.Itoa(t.timeStamp.Minute()), 2, "0"),
		strings.Lpad(strconv.Itoa(t.timeStamp.Second()), 2, "0"))
}
