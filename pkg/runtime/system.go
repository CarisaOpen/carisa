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

// meassure gets a meassure of use of the CPU and memory (MB)
func Meassure(cpu uint8, mem uint32) uint32 {
	var u uint32 = 1
	if cpu > 20 && cpu <= 40 {
		u = 2
	}
	if cpu > 40 && cpu <= 60 {
		u = 3
	}
	if cpu > 60 && cpu <= 80 {
		u = 4
	}
	if cpu > 80 {
		u = 5
	}
	return (u * 1000000) + mem
}
