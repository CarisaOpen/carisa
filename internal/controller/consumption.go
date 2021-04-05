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

package controller

import (
	"github.com/carisa/pkg/encoding"
	"github.com/rs/xid"
)

// consumption defines the struct of consumptions of the splitters
type consumption struct {
	splitterID  xid.ID
	cpu         uint8
	mem         uint32
	entesMemAvg uint32
}

// unmarshallConsumption unmarshalls a array of bytes into consumption
func unmarshallConsumption(buffer []byte) (consumption, error) {
	decoder := encoding.NewSimpleDecoder(buffer)
	idb, err := decoder.ReadBytes()
	if err != nil {
		return consumption{}, err
	}
	id, err := xid.FromBytes(idb)
	if err != nil {
		return consumption{}, err
	}
	cpu, err := decoder.ReadUint8()
	if err != nil {
		return consumption{}, err
	}
	mem, err := decoder.ReadUint32()
	if err != nil {
		return consumption{}, err
	}
	enteMemAvg, err := decoder.ReadUint32()
	if err != nil {
		return consumption{}, err
	}
	return consumption{
		splitterID:  id,
		cpu:         cpu,
		mem:         mem,
		entesMemAvg: enteMemAvg,
	}, nil
}
