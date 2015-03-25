// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package pool

import (
	"github.com/gocircuit/core/sys"
	"github.com/gocircuit/core/sys/pipe"
	"sync"
)

type pool struct {
	sync.Mutex
	x map[sys.Id]cache
}

// cache is a set of physical connections
type cache map[*pipe.Conn]struct{}

func newPool() *pool {
	return &pool{x: make(map[sys.Id]cache)}
}

func (p *pool) Add(addr sys.Addr, conn *pipe.Conn) {
	p.Lock()
	defer p.Unlock()
	h, ok := p.x[addr.Id()]
	if !ok {
		h = make(map[*pipe.Conn]struct{})
		p.x[addr.Id()] = h
	}
	h[conn] = struct{}{}
}

func (p *pool) Get(addr sys.Addr) *pipe.Conn {
	p.Lock()
	defer p.Unlock()
	h, ok := p.x[addr.Id()]
	if !ok {
		return nil
	}
	if len(h) == 0 {
		delete(p.x, addr.Id())
		return nil
	}
	for c, _ := range h {
		return c
	}
	panic(1)
}

func (p *pool) Remove(addr sys.Addr, conn *pipe.Conn) {
	p.Lock()
	defer p.Unlock()
	p.remove(addr, conn)
}

func (p *pool) remove(addr sys.Addr, conn *pipe.Conn) {
	h, ok := p.x[addr.Id()]
	if !ok {
		return
	}
	delete(h, conn)
	if len(h) == 0 {
		delete(p.x, addr.Id())
	}
}

func (p *pool) RemoveIfUnused(addr sys.Addr, conn *pipe.Conn) bool {
	p.Lock()
	defer p.Unlock()
	if conn.Count() > 0 {
		return false
	}
	p.remove(addr, conn)
	return true
}
