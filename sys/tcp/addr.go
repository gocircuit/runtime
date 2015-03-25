// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package tcp

import (
	"github.com/gocircuit/core/sys"
	"net"
)

// Addr implements sys.Addr
type Addr net.TCPAddr

func ResolveAddr(s string) (*Addr, error) {
	a, err := net.ResolveTCPAddr("tcp", s)
	if err != nil {
		return nil, err
	}
	return (*Addr)(a), nil
}

func NewAddr(u *net.TCPAddr) sys.Addr {
	return (*Addr)(u)
}

func (a *Addr) TCP() *net.TCPAddr {
	return (*net.TCPAddr)(a)
}

func (a *Addr) String() string {
	return a.TCP().String()
}

func (a *Addr) Id() sys.Id {
	return sys.Id(a.String())
}
