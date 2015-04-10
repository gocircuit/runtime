// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"github.com/gocircuit/runtime/sys"
	"io"
)

// Read reads the next pipe from the connection.
// If successful, it returns a sys.Conn object for the received pipe and a nil error.
// Otherwise, the connection has been closed and a non-nil error is returned.
func (s *Conn) Accept() (sys.Conn, error) {
	p, ok := <-s.user
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return p, nil
}

func (s *Conn) Dial() (sys.Conn, error) {
	s.x.Lock()
	s.x.n++ // skip zero, because we would get collision on it
	id := PipeId(s.sign) * s.x.n
	if s.x.open[id] != nil {
		panic("collision")
	}
	p := newPipe(id, s)
	s.x.Unlock()

	if err := s.writeOpen(p.pipeId); err != nil {
		return nil, err
	}
	s.set(p.pipeId, p)
	return p, nil
}
