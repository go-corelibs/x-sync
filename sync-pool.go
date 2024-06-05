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
	"sync"
)

var _ Pool[bool] = (*cPool[bool])(nil) // compile-time type checking

// PoolHookFn is the function signature used with NewPool
type PoolHookFn[V interface{}] func(v V) V

// Pool is a convenience interface for working with the standard sync.Pool.
//
// Pool instances can have up to two hooks provided, the first is always the
// "getter" PoolHookFn and the second is always the "setter" PoolHookFn.
//
// The "getter" function is used to modify the instance before returning it
// during Get calls
//
// The "setter" function is used to modify the instance before putting it into
// the sync.Pool
type Pool[V interface{}] interface {
	// Scale returns the amount of instances to Seed this sync.Pool when drained
	Scale() int
	// Ready returns the best-guess number of instances still in the sync.Pool
	Ready() int
	// Seed adds a count of new instances to this sync.Pool
	Seed(count int)
	// Get retrieves a typed instance from this sync.Pool, and if the pool is
	// drained, uses Seed with Scale to expand this sync.Pool
	Get() V
	// Put recycles an existing instance, ignores nil instances
	Put(v V)
}

type cPool[V interface{}] struct {
	ready  int // ready is an estimate
	scale  int
	maker  func() V
	getter PoolHookFn[V]
	setter PoolHookFn[V]
	pool   sync.Pool // pool is the underlying storage
	m      *sync.RWMutex
}

// NewPool constructs a new Pool instance with the given scale, maker function
// and the two optional PoolHookFn functions
//
// Passing nil for the hooks is valid, for example to create a Pool without a
// getter hook but with a setter, pass a nil for the first hooks argument
//
// Scale values less than or equal to zero are clamped to a minimum scale of 1
func NewPool[V interface{}](scale int, maker func() V, hooks ...PoolHookFn[V]) Pool[V] {
	if scale <= 0 {
		scale = 1
	}
	p := &cPool[V]{
		ready: 0,
		scale: scale,
		maker: maker,
		m:     &sync.RWMutex{},
	}
	switch len(hooks) {
	case 1:
		p.getter = hooks[0]
	case 2:
		p.getter = hooks[0]
		p.setter = hooks[1]
	}
	p.pool = sync.Pool{New: p.new}
	p.Seed(scale)
	return p
}

func (p *cPool[V]) new() interface{} {
	// sync.Pool has been drained
	return p.maker()
}

func (p *cPool[V]) Scale() int {
	return p.scale
}

func (p *cPool[V]) Seed(count int) {
	if count <= 0 {
		count = p.scale
	}
	for i := 0; i < count; i++ {
		p.pool.Put(p.maker())
	}
	p.m.Lock()
	p.ready += count
	p.m.Unlock()
}

func (p *cPool[V]) Ready() int {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.ready
}

func (p *cPool[V]) Get() V {
	ready := p.Ready()
	if ready -= 1; ready <= 0 {
		p.m.Lock()
		p.ready = 0
		p.m.Unlock()
		p.Seed(p.scale)
	} else {
		p.m.Lock()
		p.ready -= 1
		p.m.Unlock()
	}
	if p.getter != nil {
		return p.getter(p.pool.Get().(V))
	}
	return p.pool.Get().(V)
}

func (p *cPool[V]) Put(v V) {
	var vn V
	if interface{}(v) != interface{}(vn) && p.setter != nil {
		// let the updater update things
		v = p.setter(v)
	}
	if interface{}(v) != interface{}(vn) {
		// no setter or setter didn't return nil
		p.m.Lock()
		p.ready += 1
		p.m.Unlock()
		p.pool.Put(v)
	}
}
