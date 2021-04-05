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

package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystem_Meassure(t *testing.T) {
	tests := []struct {
		name     string
		cpu      uint8
		mem      uint32
		meassure uint32
	}{
		{
			name:     "10% CPU",
			cpu:      10,
			mem:      1024,
			meassure: 1001024,
		},
		{
			name:     "40% CPU",
			cpu:      40,
			mem:      256,
			meassure: 2000256,
		},
		{
			name:     "60% CPU",
			cpu:      60,
			mem:      2048,
			meassure: 3002048,
		},
		{
			name:     "70% CPU",
			cpu:      70,
			mem:      512,
			meassure: 4000512,
		},
		{
			name:     "81% CPU",
			cpu:      81,
			mem:      1024,
			meassure: 5001024,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.meassure, Meassure(tt.cpu, tt.mem))
	}
}
