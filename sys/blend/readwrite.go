// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package blend

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gocircuit/alef/sys/tele/codec"
)

// AcceptSession
type AcceptSession struct {
	Session
}

func newAcceptSession(sub *codec.Conn) *AcceptSession {
	ab := &AcceptSession{}
	ab.init(make(chan *conn, 1), sub, nil)
	return ab
}

func (as *AcceptSession) Accept() *conn {
	return <-as.Session.ach
}

// DialSession
type DialSession struct {
	Session
}

func newDialSession(sub *codec.Conn, scrub func()) *DialSession {
	db := &DialSession{}
	db.init(nil, sub, scrub)
	return db
}

func (ds *DialSession) Dial() *conn {
	ds.Session.o_.Lock()
	defer ds.Session.o_.Unlock()
	if ds.Session.o_open[ds.o_ndial] != nil {
		panic("u")
	}
	conn := newConn(ds.Session.o_ndial, &ds.Session)
	ds.Session.o_open[ds.Session.o_ndial] = conn
	ds.Session.o_ndial++
	return conn
}

// session
type session struct {
	scrub func()
	ach   chan *conn
	//
	x struct {
		lk    sync.Mutex // Sync o_open conns structure
		ndial ConnId
		open  map[ConnId]*conn
		use   time.Time
	}
	r struct {
		u sys.Conn // underlying connection for reading
	}
	w struct {
		lk sync.Mutex // Linearize write ops on sub
		u  sys.Conn
	}
}

func (s *session) init(acceptChan chan *conn, under sys.Conn, scrub func()) {
	s.scrub = scrub
	s.ach = acceptChan
	s.r.u, s.w.u = under, under
	s.x.open = make(map[ConnId]*conn)
	s.x.use = time.Now()

	go s.readloop()
}

func (s *session) NumConn() (numconn int, lastuse time.Time) {
	s.x.lk.Lock()
	defer s.x.lk.Unlock()
	return len(s.x.open), s.x.use
}

func (s *session) Addr() sys.Addr {
	return s.r.Addr()
}

func (s *session) hijack() (u sys.Conn) {
	s.w.lk.Lock()
	defer s.w.lk.Unlock()
	u, s.w.u = s.w.u, nil
	return
}

func (s *session) Close() (err error) {
	u := s.hijack()
	if u == nil {
		return io.ErrClosedPipe
	}
	return u.Close()
}

func (s *session) teardown() {
	// Notify accepters, if an accept session
	if s.ach != nil {
		close(s.ach)
	}

	s.x.lk.Lock()
	// The substrate connection does not allow Write after Close.
	// To prevent writes from Conns hitting the substrate before the Conns have been notified:
	// we first remove the substrate from its field to prevents writes from Conn going through to it,
	// and then we close the substrate.
	if u := s.hijack(); u != nil {
		u.Close()
	}
	// Notify open connections
	for id, conn := range s.x.open {
		conn.promptClose()
		delete(s.x.open, id)
	}
	s.x.lk.Unlock()

	if s.scrub != nil {
		s.scrub()
	}
}

func (s *session) readloop() {
	defer s.teardown()
	for {
		if err := s.read(); err != nil {
			return
		}
	}
}

func (s *session) read() error {
	msg := &Msg{}
	if err := s.r.Read(msg); err != nil {
		return err
	}

	switch t := msg.Msg.(type) {
	case *PayloadMsg:
		conn := s.get(msg.ConnId)
		if conn != nil {
			// Existing connection
			conn.prompt(t.Payload, nil)
			return nil
		}
		// Dead connection
		if t.SeqNo > 0 {
			s.writeAbort(msg.ConnId, ErrGone)
			return nil
		}
		// New connection
		if s.ach != nil {
			conn = newConn(msg.ConnId, s)
			s.set(msg.ConnId, conn)
			conn.prompt(t.Payload, nil)
			s.ach <- conn // Send new connection to user
			return nil
		} else {
			s.writeAbort(msg.ConnId, ErrOff)
			return nil
		}

	case *AbortMsg:
		conn := s.get(msg.ConnId)
		if conn == nil {
			// Discard CLOSE for non-existent connections
			// Do not respond with a CLOSE packet. It would cause an avalanche of CLOSEs.
			return nil
		}
		s.scrub(msg.ConnId)
		conn.prompt(nil, t.Err)
		return nil
	}

	// Unexpected remote behavior
	return ErrClash
}

func (s *session) count() int {
	s.x.lk.Lock()
	defer s.x.lk.Unlock()
	return len(s.x.open)
}

func (s *session) get(id ConnId) *conn {
	s.x.lk.Lock()
	defer s.x.lk.Unlock()
	s.x.use = time.Now()
	return s.x.open[id]
}

func (s *session) set(id ConnId, conn *conn) {
	s.x.lk.Lock()
	defer s.x.lk.Unlock()
	s.x.open[id] = conn
}

func (s *session) scrub(id ConnId) {
	s.x.lk.Lock()
	defer s.x.lk.Unlock()
	delete(s.x.open, id)
}

func (s *session) write(msg *Msg) error {
	s.w.lk.Lock()
	defer s.w.lk.Unlock()
	if s.w.u == nil {
		return io.ErrUnexpectedEOF
	}
	return s.w.u.Write(msg)
}

func (s *session) writeAbort(id ConnId, reason error) error {
	msg := &Msg{
		ConnId: id,
		Msg: &AbortMsg{
			Err: reason,
		},
	}
	return s.write(msg)
}
