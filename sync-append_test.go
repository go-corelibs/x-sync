// Copyright (c) 2024  The Go-CoreLibs Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sync

import (
	"fmt"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestAppend(t *testing.T) {
	c.Convey("Correctness", t, func() {

		for idx, test := range []struct {
			source   []int
			append   []int
			output   []int
			length   int
			capacity int
		}{
			{source: nil, append: []int{1}, output: []int{1}, length: 1, capacity: 4},
			{source: []int{1}, append: []int{2, 10}, output: []int{1, 2, 10}, length: 3, capacity: 8},
		} {

			output := Append(test.source, test.append...)

			c.SoMsg(
				fmt.Sprintf("test #%d (output)", idx),
				output,
				c.ShouldEqual,
				test.output,
			)

			c.SoMsg(
				fmt.Sprintf("test #%d (length)", idx),
				len(output),
				c.ShouldEqual,
				test.length,
			)

			c.SoMsg(
				fmt.Sprintf("test #%d (capacity)", idx),
				cap(output),
				c.ShouldEqual,
				test.capacity,
			)

		}

		for idx, test := range []struct {
			source   []int
			append   []int
			output   []int
			length   int
			capacity int
			scale    float64
		}{
			{source: nil, append: []int{1}, output: []int{1}, length: 1, capacity: 3, scale: 1.5},
			{source: nil, append: []int{1}, output: []int{1}, length: 1, capacity: 3, scale: 1.45}, // 1.5 is floor
			{source: nil, append: []int{1}, output: []int{1}, length: 1, capacity: 5, scale: 2.5},
			{source: []int{1}, append: []int{2, 10}, output: []int{1, 2, 10}, length: 3, capacity: 10, scale: 2.5},
		} {

			output := AppendScaled(test.scale, test.source, test.append...)

			c.SoMsg(
				fmt.Sprintf("test #%d (output)", idx),
				output,
				c.ShouldEqual,
				test.output,
			)

			c.SoMsg(
				fmt.Sprintf("test #%d (length)", idx),
				len(output),
				c.ShouldEqual,
				test.length,
			)

			c.SoMsg(
				fmt.Sprintf("test #%d (capacity)", idx),
				cap(output),
				c.ShouldEqual,
				test.capacity,
			)

		}

	})
}

func BenchmarkAppend_Int_X_Sync(b *testing.B) {
	var data []int
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			data = Append(data, i, i*2, i*10)
		} else {
			data = Append(data, i, i*-1, i*10*-1)
		}
	}
}

func BenchmarkAppend_Int_Golang(b *testing.B) {
	var data []int
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			data = append(data, i, i*2, i*10)
		} else {
			data = append(data, i, i*-1, i*10*-1)
		}
	}
}

func BenchmarkAppend_String_X_Sync(b *testing.B) {
	var data []string
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			data = AppendScaled(2, data, "even", "more")
			//data = Append(data, "even", "more")
		} else {
			data = AppendScaled(2, data, "odd", "stands", "out")
			//data = Append(data, "odd", "stands", "out")
		}
	}
	data = Append(data[:10], "truncated")
}

func BenchmarkAppend_String_Golang(b *testing.B) {
	var data []string
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			data = append(data, "even", "more")
		} else {
			data = append(data, "odd", "stands", "out")
		}
	}
	data = append(data[:10], "truncated")
}
