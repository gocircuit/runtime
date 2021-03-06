// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"errors"
	"github.com/gocircuit/runtime/sys"
	"io"
	"sync"
)

var ErrClash = errors.New("clash")
var ErrGone = errors.New("gone")

// Conn implements sys.Conn.
type Conn struct {
	user chan *pipe // Send newly received pipes to user (via Read method)
	sign int        // +1 or -1, names the side we are on on this connection
	x    struct {   // Index of known open pipes
		sync.Mutex
		n    PipeId           // number of pipes created on this end
		open map[PipeId]*pipe // pipes created on this end of the connection
	}
	w struct {
		sync.Mutex // Linearize write ops on sub
		u          sys.Conn
	}
	r sys.Conn // underlying connection for reading
}

func newConn(under sys.Conn, sign int) *Conn {
	s := &Conn{
		sign: sign,
		user: make(chan *pipe, 1),
	}
	s.r, s.w.u = under, under
	s.x.open = make(map[PipeId]*pipe)

	go s.recvloop()
	return s
}

// Count returns the number of pipes that are still open.
func (s *Conn) Count() (npipe int) {
	s.x.Lock()
	defer s.x.Unlock()
	return len(s.x.open)
}

// Addr implements sys.Conn.Addr.
func (s *Conn) Addr() sys.Addr {
	return s.r.Addr()
}

func (s *Conn) hijack() (u sys.Conn) {
	s.w.Lock()
	defer s.w.Unlock()
	u, s.w.u = s.w.u, nil
	return
}

// Close will abort any unclosed pipes.
func (s *Conn) Close() (err error) {
	u := s.hijack()
	if u == nil {
		return io.ErrClosedPipe
	}
	// The closure of u is sensed in the recvloop, which in turn
	// triggers the teardown sequence for this connection (notifying
	// all outstanding pipes that they have been broken).
	return u.Close()
}

func (s *Conn) teardown() {
	// Notify Read (the reading user) that the connection is broken.
	if s.user != nil {
		close(s.user)
	}

	s.x.Lock()
	// The substrate connection does not allow Write after Close.
	// To prevent writes from Conns hitting the substrate before the Conns have been notified:
	// we first remove the substrate from its field to prevents writes from Conn going through to it,
	// and then we close the substrate.
	if u := s.hijack(); u != nil {
		u.Close()
	}
	// Notify open pipes that they are now broken
	for id, p := range s.x.open {
		p.userClose()
		delete(s.x.open, id)
	}
	s.x.Unlock()
}

func (s *Conn) recvloop() {
	defer s.teardown()
	for {
		if err := s.receive(); err != nil {
			return
		}
	}
}

func (s *Conn) receive() error {
	t, err := s.r.Receive()
	if err != nil {
		return err
	}
	msg, ok := t.(*Msg)
	if !ok {
		return ErrClash
	}

	switch t := msg.Msg.(type) {
	case nil: // Introduce a new pipe
		if s.get(msg.PipeId) != nil {
			return ErrClash // Collision of pipe ids
		}
		p := newPipe(msg.PipeId, s)
		s.set(msg.PipeId, p)
		s.user <- p // Send new pipe to user
		return nil

	case *PayloadMsg:
		p := s.get(msg.PipeId)
		if p == nil { // Dead pipe
			s.writeAbort(msg.PipeId, ErrGone)
			return nil
		}
		p.userSend(t.Payload, nil)
		return nil

	case *AbortMsg:
		p := s.get(msg.PipeId)
		if p == nil {
			// Discard closures for non-existent pipes
			// Do not respond with an abort message. This would cause an avalanche of abort messages.
			return nil
		}
		s.scrub(msg.PipeId)
		p.userSend(nil, t.Reason)
		return nil
	}

	// Unexpected remote behavior
	return ErrClash
}

func (s *Conn) count() int {
	s.x.Lock()
	defer s.x.Unlock()
	return len(s.x.open)
}

func (s *Conn) get(id PipeId) *pipe {
	s.x.Lock()
	defer s.x.Unlock()
	return s.x.open[id]
}

func (s *Conn) set(id PipeId, p *pipe) {
	s.x.Lock()
	defer s.x.Unlock()
	if _, present := s.x.open[id]; present {
		panic("collision")
	}
	s.x.open[id] = p
}

func (s *Conn) scrub(id PipeId) {
	s.x.Lock()
	defer s.x.Unlock()
	delete(s.x.open, id)
}

func (s *Conn) write(msg interface{}) error {
	s.w.Lock()
	defer s.w.Unlock()
	if s.w.u == nil {
		return io.ErrUnexpectedEOF
	}
	return s.w.u.Send(msg)
}

func (s *Conn) writeOpen(id PipeId) error {
	msg := &Msg{
		PipeId: id,
	}
	return s.write(msg)
}

func (s *Conn) writePayload(id PipeId, paymsg *PayloadMsg) error {
	msg := &Msg{
		PipeId: id,
		Msg:    paymsg,
	}
	return s.write(msg)
}

func (s *Conn) writeAbort(id PipeId, reason error) error {
	msg := &Msg{
		PipeId: id,
		Msg: &AbortMsg{
			Reason: reason,
		},
	}
	return s.write(msg)
}
