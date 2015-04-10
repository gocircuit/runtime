// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"encoding/gob"
	"io"
	"net"
	"sync"

	"github.com/gocircuit/runtime/sys"
)

type sandbox struct {
	lk sync.Mutex
	l  map[sys.Id]*listener
}

var s = &sandbox{l: make(map[sys.Id]*listener)}

// NewSandbox creates a new transport instance, part of a sandbox network in memory
func NewSandbox() sys.Peer {
	s.lk.Lock()
	defer s.lk.Unlock()

	l := &listener{
		id: sys.ChooseId(),
		ch: make(chan *halfconn),
	}
	l.a = &addr{ID: l.id, l: l}
	s.l[l.id] = l
	return l
}

func (l *listener) Listen(net.Addr) sys.Listener {
	return l
}

func dial(remote sys.Addr) (sys.Conn, error) {
	pr, pw := io.Pipe()
	qr, qw := io.Pipe()
	srvhalf := &halfconn{PipeWriter: qw, PipeReader: pr}
	clihalf := &halfconn{PipeWriter: pw, PipeReader: qr}
	s.lk.Lock()
	l := s.l[remote.(*addr).Id()]
	s.lk.Unlock()
	if l == nil {
		panic("unknown listener id")
	}
	go func() {
		l.ch <- srvhalf
	}()
	return ReadWriterConn(l.Addr(), clihalf), nil
}

// addr implements Addr
type addr struct {
	ID sys.Id
	l  *listener
}

func (a *addr) Id() sys.Id {
	return a.ID
}

func (a *addr) Network() string {
	return "sandbox"
}

func (a *addr) String() string {
	return a.ID.String()
}

func init() {
	gob.Register(&addr{})
}

// listener implements Listener
type listener struct {
	id sys.Id
	a  *addr
	ch chan *halfconn
}

func (l *listener) Addr() sys.Addr {
	return l.a
}

func (l *listener) Accept() (sys.Conn, error) {
	return ReadWriterConn(l.Addr(), <-l.ch), nil
}

func (l *listener) Close() {
	s.lk.Lock()
	defer s.lk.Unlock()
	delete(s.l, l.id)
}

func (l *listener) Dial(remote sys.Addr) (sys.Conn, error) {
	return dial(remote)
}

// halfconn is one end of a byte-level connection
type halfconn struct {
	*io.PipeReader
	*io.PipeWriter
}

func (h *halfconn) Close() error {
	return h.PipeWriter.Close()
}
