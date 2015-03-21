// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package tcp

import (
	"github.com/gocircuit/alef/peer"
	"net"
)

func New(addr peer.Addr) (_ peer.Peer, err error) {
	p := &peer{}
	if p.l, err = net.ListenTCP("tcp", t); err != nil {
		return nil, err
	}
	return p, nil
}

// peer implements peer.Peer
type peer struct {
	l *net.TCPListener
}

func (p *peer) Accept() peer.Conn {
	c, err := p.l.Accept()
	if err != nil {
		panic(err)
	}
	return newConn(c)
}

func (p *peer) Addr() peer.Addr {
	return NewAddr(p.l.Addr())
}

func (p *peer) Dial(addr peer.Addr) (peer.Conn, error) {
	c, err := net.DialTCP("tcp", nil, addr.(*Addr).TCP())
	if err != nil {
		return nil, err
	}
	return newConn(c), nil
}
