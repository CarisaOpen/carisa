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
	rx "runtime"
	"time"

	"github.com/carisa/pkg/runtime"
)

// consumption updates the actual consumption each wutime in seconds
type consumption struct {
	actual time.Time
	wutime time.Duration

	pmeasure int
	cpu      uint8
	mem      int
}

func newConsumption(wutime time.Duration) consumption {
	return consumption{
		actual: time.Now(),
		wutime: wutime,
	}
}

// wake warns if the renewal period is met
func (c *consumption) wake() bool {
	if time.Since(c.actual) > c.wutime*time.Second {
		c.actual = time.Now()
		return true
	}
	return false
}

// renew updates the cpu (avarage) and memory measure
func (c *consumption) renew() error {
	var m rx.MemStats
	rx.ReadMemStats(&m)
	c.mem = int(m.Alloc / 1024 / 1024) // MB

	cpu, err := runtime.CPU()
	if err != nil {
		return err
	}
	c.cpu = (c.cpu + uint8(cpu)) / 2 // Average
	return nil
}

// measure getting measure based on the cpu range and memory
func (c *consumption) measure() int {
	u := 1
	if c.cpu > 20 && c.cpu <= 40 {
		u = 2
	}
	if c.cpu > 40 && c.cpu <= 60 {
		u = 3
	}
	if c.cpu > 60 && c.cpu <= 80 {
		u = 4
	}
	if c.cpu > 80 {
		u = 5
	}
	return (u * 1000000) + c.mem
}

// saveMeasure saves the actual measure
func (c *consumption) saveMeasure() {
	c.pmeasure = c.measure()
}
