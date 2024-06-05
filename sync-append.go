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

// Append is an experiment in appending to slices without using the standard
// Go append function
func Append[V interface{}](slice []V, data ...V) []V {
	m := len(slice)     // current length
	n := m + len(data)  // needed length
	if n > cap(slice) { // current cap size
		grown := make([]V, (n+1)*2) // new double cap slice
		copy(grown, slice)          // transfer to new slice
		slice = grown               // grown becomes slice
	}
	slice = slice[0:n]     // truncate in case too many present
	copy(slice[m:n], data) // populate with data
	return slice
}

// AppendScaled is like Append but includes the scale argument for specifying
// the cap growth multiplier. The scale value is capped at a minimum of 1.5
// and for reference, the Append function uses a hard-coded scale of 2
func AppendScaled[V interface{}](scale float64, slice []V, data ...V) []V {
	if scale < 1.5 {
		scale = 1.5
	}
	m := len(slice)     // current length
	n := m + len(data)  // needed length
	if n > cap(slice) { // current cap size
		grown := make([]V, int(float64(n+1)*scale)) // new scaled cap slice
		copy(grown, slice)                          // transfer to new slice
		slice = grown                               // grown becomes slice
	}
	slice = slice[0:n]     // truncate in case too many present
	copy(slice[m:n], data) // populate with data
	return slice
}
