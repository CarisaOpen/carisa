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

	"github.com/rs/xid"

	"github.com/carisa/pkg/encoding"

	"github.com/carisa/pkg/runtime"
)

const consumptionSchema = "SC"

// consumption updates the actual consumption each wutime in seconds
type consumption struct {
	actual time.Time
	wutime time.Duration

	pmeasure uint32
	cpu      uint8
	mem      uint32
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
	c.mem = uint32(m.Alloc / 1024 / 1024) // MB

	cpu, err := runtime.CPU()
	if err != nil {
		return err
	}
	c.cpu = (c.cpu + uint8(cpu)) / 2 // Average
	return nil
}

// reg returns a struct of consumption
func (c *consumption) reg(id xid.ID, entesMemAvg uint32) string {
	sc := encoding.SimpleEncoder{}
	sc.WriteBytes(id.Bytes())
	sc.WriteUint8(c.cpu)
	sc.WriteUint32(c.mem)
	sc.WriteUint32(entesMemAvg)
	return sc.String()
}

// measure getting measure based on the cpu range and memory
func (c *consumption) measure() uint32 {
	return runtime.Meassure(c.cpu, c.mem)
}

// saveMeasure saves the actual measure
func (c *consumption) saveMeasure() {
	c.pmeasure = c.measure()
}
