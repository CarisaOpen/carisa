#!/usr/bin/env bash

#
# Copyright 2019-2022 the original author or authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software  distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and  limitations under the License.
#
#

MOD=$(grep "module" go.mod)
MOD="${MOD/module /}"

EXCLUDE=""
while read n; do
  EXCLUDE+=$"$MOD/$n|"
done < ./coverage/exclude

EXCLUDE=${EXCLUDE%?};

go test -coverprofile=coverage.out -coverpkg ./... ./...
grep -vwE "($EXCLUDE)" coverage.out > coverage-final.out
go tool cover -func=coverage-final.out

rm coverage.out
rm coverage-final.out