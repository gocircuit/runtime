// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"github.com/gocircuit/core/sys"
)

func New(u sys.Peer) *Peer {
	return &Peer{u}
}

type Peer struct {
	u sys.Peer
}

func (p *Peer) Dial(addr sys.Addr) (*Conn, error) {
	c, err := p.u.Dial(addr)
	if err != nil {
		return nil, err
	}
	return newConn(c, 1), nil
}

func (p *Peer) Accept() (*Conn, error) {
	c, err := p.u.Accept()
	if err != nil {
		return nil, err
	}
	return newConn(c, -1), nil
}

func (p *Peer) Addr() sys.Addr {
	return p.u.Addr()
}
