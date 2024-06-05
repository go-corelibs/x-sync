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
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	c.Convey("NewStringBuilderPool", t, func() {

		p := NewStringBuilderPool(1)

		c.So(p.Scale(), c.ShouldEqual, 1)
		c.So(p.Ready(), c.ShouldEqual, 1)

		v0 := p.Get()
		c.So(v0, c.ShouldNotBeNil)
		c.So(p.Ready(), c.ShouldEqual, 1) // seeded
		c.So(v0.String(), c.ShouldEqual, "")
		v0.WriteString(`life at the speed of thought`)

		p.Put(v0)
		c.So(p.Ready(), c.ShouldEqual, 2)
		c.So(v0.String(), c.ShouldEqual, `life at the speed of thought`)

		v1 := p.Get()
		c.So(v1, c.ShouldNotBeNil)
		c.So(p.Ready(), c.ShouldEqual, 1)
		c.So(v1.String(), c.ShouldEqual, ``)

		v2 := p.Get()
		c.So(p.Ready(), c.ShouldEqual, 1)
		v2.WriteString(strings.Repeat("0", 64001))
		c.So(v2.Len(), c.ShouldEqual, 64001)
		p.Put(v2)
		c.So(p.Ready(), c.ShouldEqual, 1)

	})
}
