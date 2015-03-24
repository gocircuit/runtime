// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gocircuit/core/sys/tele/codec"
)

// Read reads the next pipe from the connection.
// If successful, it returns a sys.Conn object for the received pipe and a nil error.
// Otherwise, the connection has been closed and a non-nil error is returned.
func (s *Conn) Read() (interface{}, error) {
	p, ok := <-s.user
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return p, nil
}

func (s *Conn) NewPipe() sys.Conn {
	s.x.Lock()
	defer s.x.Unlock()
	id := s.sign * s.x.n
	if s.x.open[id] != nil {
		panic("collision")
	}
	p := newPipe(id, s)
	s.x.n++
	return p
}

func (s *Conn) Write(p interface{}) (err error) {
	q, ok := p.(*pipe)
	if !ok {
		panic("can only write pipes to this connection")
	}
	if err = s.writePayload(q.pipeId, nil); err != nil {
		return
	}
	s.set(q.pipeId, q)
	return nil
}
