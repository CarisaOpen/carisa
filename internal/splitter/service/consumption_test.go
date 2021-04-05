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

	"github.com/carisa/pkg/encoding"

	"github.com/rs/xid"

	"github.com/stretchr/testify/assert"
)

func TestConsumption_Renew(t *testing.T) {
	c := newConsumption(1 * time.Second)
	_ = c.renew()
	assert.GreaterOrEqual(t, c.measure(), uint32(1000000))
}

func TestConsumption_Meassure(t *testing.T) {
	tests := []struct {
		name     string
		cpu      uint8
		meassure uint32
	}{
		{
			name:     "10% CPU",
			cpu:      10,
			meassure: 1000000,
		},
		{
			name:     "40% CPU",
			cpu:      40,
			meassure: 2000000,
		},
		{
			name:     "60% CPU",
			cpu:      60,
			meassure: 3000000,
		},
		{
			name:     "70% CPU",
			cpu:      70,
			meassure: 4000000,
		},
		{
			name:     "81% CPU",
			cpu:      81,
			meassure: 5000000,
		},
	}

	c := newConsumption(1 * time.Second)

	for _, tt := range tests {
		c.cpu = tt.cpu
		assert.Equal(t, tt.meassure, c.measure())
	}
}

func TestConsumption_SaveMeasure(t *testing.T) {
	c := newConsumption(1 * time.Second)
	_ = c.renew()
	c.saveMeasure()
	assert.Equal(t, c.pmeasure, c.measure())
}

func TestConsumption_Wake(t *testing.T) {
	c := newConsumption(1)
	time.Sleep(1 * time.Second)
	assert.True(t, c.wake())
}

func TestConsumption_No_Wake(t *testing.T) {
	c := newConsumption(5)
	assert.False(t, c.wake())
}

func TestConsumption_Reg(t *testing.T) {
	c := newConsumption(1)
	c.cpu = 5
	c.mem = 1024
	id := xid.New()
	reg := c.reg(id, 1024)

	d := encoding.NewSimpleDecoder([]byte(reg))
	idb, _ := d.ReadBytes()
	idr, _ := xid.FromBytes(idb)
	assert.Equal(t, id, idr, "id")
	cpu, _ := d.ReadUint8()
	assert.Equal(t, c.cpu, cpu, "cpu")
	mem, _ := d.ReadUint32()
	assert.Equal(t, c.mem, mem, "mem")
	enteMemAvg, _ := d.ReadUint32()
	assert.Equal(t, uint32(1024), enteMemAvg, "enteMemAvg")
}
