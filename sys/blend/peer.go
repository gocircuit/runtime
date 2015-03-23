// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"github.com/gocircuit/alef/sys"
)

type peer struct {
	u sys.Peer
}

func New(u sys.Peer) sys.Peer {
	return &Transport{u: u}
}

func (p *peer) Dial(addr sys.Addr) (sys.Conn, error) {
}

func (p *peer) Dial(addr sys.Addr, scrb func()) (*DialSession, error) {
	sub, err := d.sub.Dial(addr)
	if err != nil {
		return nil, err
	}
	return newDialSession(d.frame.Refine("dial"), sub, scrb), nil // codec.Dial always returns instantaneously
}

func (p *peer) AcceptSession() *AcceptSession {
	sub := l.sub.Accept()
	if sub == nil {
		panic("accepted nil conn")
	}
	return newAcceptSession(l.frame.Refine("accept"), sub)
}

func (p *peer) Addr() sys.Addr {
	return p.u.Addr()
}
