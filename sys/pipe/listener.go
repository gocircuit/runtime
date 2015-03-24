// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"net"

	"github.com/gocircuit/alef/sys/tele/codec"
	"github.com/gocircuit/alef/sys/tele/trace"
)

type Listener struct {
	frame trace.Frame
	sub   *codec.Listener
}

func NewListener(frame trace.Frame, sub *codec.Listener) *Listener {
	l := &Listener{frame: frame, sub: sub}
	frame.Bind(l)
	return l
}

func (l *Listener) AcceptSession() *AcceptSession {
	sub := l.sub.Accept()
	if sub == nil {
		panic("accepted nil conn")
	}
	return newAcceptSession(l.frame.Refine("accept"), sub)
}

func (l *Listener) Addr() net.Addr {
	return l.sub.Addr()
}
