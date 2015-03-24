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

func New(u sys.Peer) sys.Peer {
	return &peer{u}
}

type peer struct {
	u sys.Peer
}

func (p *peer) Dial(addr sys.Addr) (sys.Conn, error) {
	c, err := p.u.Dial(addr)
	if err != nil {
		return nil, err
	}
	return newConn(c, 1), nil
}

func (p *peer) Accept() sys.Conn {
	return newConn(p.u.Accept(), -1)
}

func (p *peer) Addr() sys.Addr {
	return p.u.Addr()
}
