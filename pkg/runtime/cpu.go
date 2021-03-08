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
	"errors"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	s "github.com/carisa/pkg/strings"
)

var clkTck float64

func init() {
	clkTckStdout, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		panic(err)
	}
	aclkTckStdout := strings.Split(str(clkTckStdout), "\n")
	if len(aclkTckStdout) == 0 {
		panic(err)
	}
	clkTck = parseFloat(aclkTckStdout[0])
}

// CPU returns % use of CPU
func CPU() (float64, error) {
	buptime, err := ioutil.ReadFile(path.Join("/proc", "uptime"))
	if err != nil {
		return 0, err
	}
	auptime := strings.Split(str(buptime), " ")
	if len(auptime) == 0 {
		return 0, errors.New("can't read the uptime")
	}
	uptime := parseFloat(auptime[0])

	procStatFileBytes, err := ioutil.ReadFile(path.Join("/proc", strconv.Itoa(os.Getpid()), "stat"))
	if err != nil {
		return 0, err
	}
	splitAfter := strings.SplitAfter(str(procStatFileBytes), ")")
	if len(splitAfter) == 0 || len(splitAfter) == 1 {
		return 0, errors.New(s.Concat("can't find process with this PID: ", strconv.Itoa(os.Getpid())))
	}
	infos := strings.Split(splitAfter[1], " ")
	if len(infos) == 0 {
		return 0, errors.New("can't read the stat process")
	}
	utime := parseFloat(infos[12])
	stime := parseFloat(infos[13])
	cutime := parseFloat(infos[14])
	cstime := parseFloat(infos[15])
	start := parseFloat(infos[20]) / clkTck

	total := stime + utime + cstime + cutime
	seconds := uptime - start
	seconds = math.Abs(seconds)
	if seconds == 0 {
		seconds = 1
	}

	return 100 * ((total / clkTck) / seconds), nil
}

func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 32)
	return floatVal
}

func str(byt []byte) string {
	var b strings.Builder
	b.Grow(len(byt))
	b.Write(byt)
	return b.String()
}
