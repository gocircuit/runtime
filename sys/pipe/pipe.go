// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"github.com/gocircuit/runtime/sys"
	"io"
	"sync"
)

// pipe implements sys.Conn.
type pipe struct {
	pipeId PipeId
	conn   connPipe
	user   struct {
		sync.Mutex // write-to-user channel
		ch         chan *readReturn
		eof        bool
	}
	write struct {
		sync.Mutex
		n   SeqNo // Number of writes
		eof bool  // Write-side closure
	}
}

// Interface of underlying connection that faces the pipe implementation.
type connPipe interface {
	writePayload(PipeId, *PayloadMsg) error
	scrub(PipeId)
}

type readReturn struct {
	Payload interface{}
	Err     error
}

func newPipe(pipeId PipeId, conn connPipe) *pipe {
	p := &pipe{
		pipeId: pipeId,
		conn:   conn,
	}
	p.user.ch = make(chan *readReturn, 3)
	return p
}

// Addr returns nil.
func (p *pipe) Addr() sys.Addr {
	panic(1) // pipes don't support remote addresses
	return nil
}

// Receive reads the next object.
func (p *pipe) Receive() (interface{}, error) {
	rr, ok := <-p.user.ch
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return rr.Payload, rr.Err
}

func (p *pipe) userSend(payload interface{}, err error) {
	p.user.Lock()
	defer p.user.Unlock()
	if p.user.eof {
		return
	}
	p.user.ch <- &readReturn{
		Payload: payload,
		Err:     err,
	}
	if err != nil {
		close(p.user.ch)
		p.user.eof = true
	}
}

func (p *pipe) userClose() {
	p.user.Lock()
	defer p.user.Unlock()
	if p.user.eof {
		return
	}
	close(p.user.ch)
	p.user.eof = true
}

// Write writes the chunk to the connection.
func (p *pipe) Send(v interface{}) error {
	p.write.Lock()
	defer p.write.Unlock()
	if p.write.eof {
		panic("writing after close")
	}
	p.write.n++
	msg := &PayloadMsg{
		SeqNo:   p.write.n - 1,
		Payload: v,
	}
	return p.conn.writePayload(p.pipeId, msg)
}

// Close closes the connection. It is synchronized with Write and will not interrupt a concurring write.
func (p *pipe) Close() error {
	p.conn.scrub(p.pipeId) // Scrub outside of write lock
	//
	p.write.Lock()
	if p.write.eof {
		p.write.Unlock()
		return io.ErrUnexpectedEOF
	}
	p.write.eof = true
	p.write.Unlock()
	//
	p.userClose()
	return nil
}
