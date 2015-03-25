// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package backdial

import (
	"github.com/gocircuit/core/sys"
)

type peer struct {
	u sys.Peer
}

func New(under sys.Peer) sys.Peer {
	return &peer{under}
}

func (p *peer) Accept() (sys.Conn, error) {
	c, err := p.u.Accept()
	if err != nil {
		return nil, err
	}
	return p.handshake(c)
}

func (p *peer) Addr() sys.Addr {
	return p.u.Addr()
}

func (p *peer) Dial(addr sys.Addr) (conn sys.Conn, err error) {
	c, err := p.u.Dial(addr)
	if err != nil {
		return nil, err
	}
	return p.handshake(c)
}
