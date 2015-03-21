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

	"github.com/gocircuit/alef/peer"
)

type sandbox struct {
	lk sync.Mutex
	l  map[peer.Id]*listener
}

var s = &sandbox{l: make(map[peer.Id]*listener)}

// NewSandbox creates a new transport instance, part of a sandbox network in memory
func NewSandbox() peer.Peer {
	s.lk.Lock()
	defer s.lk.Unlock()

	l := &listener{
		id: peer.ChooseId(),
		ch: make(chan *halfconn),
	}
	l.a = &addr{ID: l.id, l: l}
	s.l[l.id] = l
	return l
}

func (l *listener) Listen(net.Addr) peer.Listener {
	return l
}

func dial(remote peer.Addr) (peer.Conn, error) {
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
	ID peer.Id
	l  *listener
}

func (a *addr) Id() peer.Id {
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
	id peer.Id
	a  *addr
	ch chan *halfconn
}

func (l *listener) Addr() peer.Addr {
	return l.a
}

func (l *listener) Accept() peer.Conn {
	return ReadWriterConn(l.Addr(), <-l.ch)
}

func (l *listener) Close() {
	s.lk.Lock()
	defer s.lk.Unlock()
	delete(s.l, l.id)
}

func (l *listener) Dial(remote peer.Addr) (peer.Conn, error) {
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
