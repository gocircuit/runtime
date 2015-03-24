// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"log"
	"sync"

	"github.com/gocircuit/core/lang/acid"
	"github.com/gocircuit/core/lang/prof"
	"github.com/gocircuit/core/lang/types"
	"github.com/gocircuit/core/sys"
)

// Runtime represents that state of the circuit program at the present moment.
// This state can change in two ways: by a 'linguistic' action ...
type Runtime struct {
	t    sys.Peer
	exp  *expTabl
	imp  *impTabl
	srv  srvTabl
	blk  sync.Mutex
	boot interface{}
	lk   sync.Mutex
	live map[sys.Id]struct{} // Set of peers we monitor for liveness
	prof *prof.Profile
	dwg  sync.WaitGroup
}

func New(t sys.Peer) *Runtime {
	r := &Runtime{
		t:    t,
		exp:  makeExpTabl(types.ValueTabl),
		imp:  makeImpTabl(types.ValueTabl),
		live: make(map[sys.Id]struct{}),
		prof: prof.New(),
	}
	r.srv.Init()
	go func() {
		for {
			if err := r.accept(t); err != nil {
				panic(err)
			}
		}
	}()
	r.Listen("acid", acid.New())
	return r
}

func (r *Runtime) ServerAddr() sys.Addr {
	return r.t.Addr()
}

func (r *Runtime) SetBoot(v interface{}) {
	r.blk.Lock()
	defer r.blk.Unlock()
	if v != nil {
		types.RegisterValue(v)
	}
	r.boot = v
}

func (r *Runtime) accept(l sys.Listener) error {
	conn, err := l.Accept()
	if err != nil {
		return err
	}
	// The transport layer assumes that the user is always blocked on
	// transport.Accept and conn.Read for all accepted connections.
	// This is achieved by forking the goroutine below.
	go func() {
		req, err := conn.Read()
		if err != nil {
			log.Println("unexpected eof conn", err.Error())
			return
		}
		// Importing reptr variables involves waiting on other runtimes,
		// we fork request handling to dedicated go routines.
		// No rate-limiting/throttling is performed in the circuit.
		// It is the responsibility of Listener and/or the user app logic to
		// keep the runtime from contending.
		switch q := req.(type) {
		case *dialMsg:
			r.serveDial(q, conn)
		case *callMsg:
			r.serveCall(q, conn)
		case *dropPtrMsg:
			r.serveDropPtr(q, conn)
		case *getPtrMsg:
			r.serveGetPtr(q, conn)
		case *dontReplyMsg:
			// Don't reply. Intentionally don't close the conn.
			// It will close when the process dies.
		default:
			log.Printf("unknown request %v", req)
		}
	}()
	return nil
}

func (r *Runtime) Hang() {
	<-(chan struct{})(nil)
}

func (r *Runtime) RegisterValue(v interface{}) {
	types.RegisterValue(v)
}
