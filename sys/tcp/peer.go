// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package tcp

import (
	"github.com/gocircuit/runtime/sys"
	"net"
)

func New(addr *Addr) (_ sys.Peer, err error) {
	p := &peer{}
	if p.l, err = net.ListenTCP("tcp", addr.TCP()); err != nil {
		return nil, err
	}
	return p, nil
}

// peer implements sys.Peer
type peer struct {
	l *net.TCPListener
}

func (p *peer) Accept() (sys.Conn, error) {
	c, err := p.l.Accept()
	if err != nil {
		return nil, err
	}
	return newConn(c.(*net.TCPConn)), nil
}

func (p *peer) Addr() sys.Addr {
	return NewAddr(p.l.Addr().(*net.TCPAddr))
}

func (p *peer) Dial(addr sys.Addr) (sys.Conn, error) {
	c, err := net.DialTCP("tcp", nil, addr.(*Addr).TCP())
	if err != nil {
		return nil, err
	}
	return newConn(c), nil
}
