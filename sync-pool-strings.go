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
	"strings"
)

// NewStringBuilderPool is a convenience wrapper around NewPool, configured
// for the *strings.Builder type (all Pool values must be pointers), and
// includes both getter and setter PoolHookFn implementations. The getter
// function resets the buffer before returning it and the setter function
// won't allow recycling buffers that are larger than 65k
func NewStringBuilderPool(scale int) (pool Pool[*strings.Builder]) {
	return NewPool[*strings.Builder](scale, func() *strings.Builder {
		return new(strings.Builder)
	}, func(v *strings.Builder) *strings.Builder {
		// getter
		v.Reset()
		return v
	}, func(v *strings.Builder) *strings.Builder {
		// setter
		if v.Len() < 64000 {
			return v
		}
		return nil
	})
}
