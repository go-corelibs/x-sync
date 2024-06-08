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
	"sync"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestPool(t *testing.T) {

	c.Convey("Internals", t, func() {

		p := NewPool(0, func() *int {
			var v int
			return &v
		})

		c.So(p.Scale(), c.ShouldEqual, 1)
		c.So(p.Ready(), c.ShouldEqual, 1)

		v := p.(*cPool[*int]).new()
		c.So(v, c.ShouldNotBeNil)
		c.So(p.Ready(), c.ShouldEqual, 1)

		p.Seed(0) // seed clamps to scale
		c.So(p.Ready(), c.ShouldEqual, 2)
		_ = p.Get()

		wg := &sync.WaitGroup{}

		wg.Add(2)

		go func() {
			p.Seed(10)
			wg.Done()
		}()
		go func() {
			p.Seed(100)
			wg.Done()
		}()

		wg.Wait() // need to wait for the competing funcs to return
		c.So(p.Ready(), c.ShouldEqual, 111)

	})

	c.Convey("Correctness", t, func() {

		p := NewPool(1, func() *strings.Builder {
			return new(strings.Builder)
		})

		c.So(p.Scale(), c.ShouldEqual, 1)
		c.So(p.Ready(), c.ShouldEqual, 1)

		v0 := p.Get()
		c.So(v0, c.ShouldNotBeNil)
		c.So(p.Ready(), c.ShouldEqual, 1) // seeded
		c.So(v0.String(), c.ShouldEqual, "")
		v0.WriteString(`life at the speed of thought`)

		p.Put(v0)
		c.So(p.Ready(), c.ShouldEqual, 2)

		v1 := p.Get()
		c.So(v1, c.ShouldNotBeNil)
		c.So(p.Ready(), c.ShouldEqual, 1)
		// unreliable validation
		//c.So(v1.String(), c.ShouldEqual, "life at the speed of thought")

	})

	c.Convey("PoolHookFn", t, func() {

		c.Convey("getter", func() {

			p := NewPool(1, func() *strings.Builder {
				return new(strings.Builder)
			}, func(v *strings.Builder) *strings.Builder {
				v.Reset()
				return v
			})

			c.So(p.Scale(), c.ShouldEqual, 1)
			c.So(p.Ready(), c.ShouldEqual, 1)

			v0 := p.Get()
			c.So(v0, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1) // seeded
			c.So(v0.String(), c.ShouldEqual, "")
			v0.WriteString(`life at the speed of thought`)

			p.Put(v0)
			c.So(p.Ready(), c.ShouldEqual, 2)

			v1 := p.Get()
			c.So(v1, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1)
			c.So(v1.String(), c.ShouldEqual, "")

		})

		c.Convey("setter", func() {

			p := NewPool(1, func() *strings.Builder {
				return new(strings.Builder)
			}, func(v *strings.Builder) *strings.Builder {
				v.WriteString("!")
				return v
			})

			c.So(p.Scale(), c.ShouldEqual, 1)
			c.So(p.Ready(), c.ShouldEqual, 1)

			v0 := p.Get()
			c.So(v0, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1) // seeded
			c.So(v0.String(), c.ShouldEqual, "!")
			v0.WriteString(`life at the speed of thought`)
			c.So(v0.String(), c.ShouldEqual, `!life at the speed of thought`)

			p.Put(v0)
			c.So(p.Ready(), c.ShouldEqual, 2)

			v1 := p.Get()
			c.So(v1, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1)
			c.So(v1.String(), c.ShouldBeIn, `!`, `!life at the speed of thought!`)

			v2 := p.Get()
			c.So(v2, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1)
			c.So(v2.String(), c.ShouldBeIn, `!`, `!life at the speed of thought!`)

			v3 := p.Get()
			c.So(v3, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1)
			c.So(v3.String(), c.ShouldEqual, `!`)

		})

		c.Convey("getter and setter", func() {

			p := NewPool(1, func() *strings.Builder {
				return new(strings.Builder)
			}, func(v *strings.Builder) *strings.Builder {
				v.Reset()
				return v
			}, func(v *strings.Builder) *strings.Builder {
				v.WriteString("!")
				return v
			})

			c.So(p.Scale(), c.ShouldEqual, 1)
			c.So(p.Ready(), c.ShouldEqual, 1)

			v0 := p.Get()
			c.So(v0, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1) // seeded
			c.So(v0.String(), c.ShouldEqual, "")
			v0.WriteString(`life at the speed of thought`)

			p.Put(v0)
			c.So(p.Ready(), c.ShouldEqual, 2)
			c.So(v0.String(), c.ShouldEqual, `life at the speed of thought!`)

			v1 := p.Get()
			c.So(v1, c.ShouldNotBeNil)
			c.So(p.Ready(), c.ShouldEqual, 1)
			c.So(v1.String(), c.ShouldEqual, ``)

		})
	})
}
