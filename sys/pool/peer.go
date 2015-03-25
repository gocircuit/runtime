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
	"time"
)

// Peer is a sys.Peer with pooling and caching of physical connections and
// multiplexed logical connections over physical ones.
type Peer struct {
	u    *pipe.Peer
	ach  chan *conn
	pool *pool
}

func New(u sys.Peer) *Peer {
	p := &Peer{
		u:    pipe.New(u),
		ach:  make(chan *conn, 1),
		pool: newPool(),
	}
	go p.accept()
	return p
}

func (p *Peer) Addr() sys.Addr {
	return p.u.Addr()
}

func (p *Peer) accept() {
	for {
		c, err := p.u.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer c.Close()
			for {
				q, err := c.Accept()
				if err != nil {
					return
				}
				p.ach <- newConn(q, c.Addr())
			}
		}()
	}
}

func (p *Peer) Accept() (sys.Conn, error) {
	c, ok := <-p.ach
	if !ok {
		panic(1)
	}
	return c, nil
}

func (p *Peer) Dial(addr sys.Addr) (r sys.Conn, err error) {
	c := p.pool.Get(addr)
	for i := 0; i < 2; i++ {
		if c == nil { // Create new physical connection
			if c, err = p.u.Dial(addr); err != nil {
				return nil, err
			}
			go p.expire(addr, c)
			p.pool.Add(addr, c)
		}
		if r, err = c.Dial(); err != nil { // Create logic pipe connection
			c.Close()
			p.pool.Remove(addr, c)
			continue // If the physical connection is broken, re-attempt once
		}
		return newConn(r, addr), nil
	}
	return nil, err
}

func (p *Peer) expire(addr sys.Addr, c *pipe.Conn) {
	for {
		time.Sleep(time.Second * 10)
		if p.pool.RemoveIfUnused(addr, c) {
			return
		}
	}
}
