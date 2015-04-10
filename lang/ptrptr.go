// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	"strings"

	"github.com/gocircuit/runtime/circuit"
	"github.com/gocircuit/runtime/sys"
)

func (r *Runtime) callGetPtr(srcID circuit.HandleID, exporter sys.Addr) (circuit.X, error) {
	conn, err := r.t.Dial(exporter)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rvmsg, err := writeReturn(conn, &getPtrMsg{ID: srcID})
	if err != nil {
		return nil, err
	}

	return r.importEitherPtr(rvmsg, exporter)
}

func (r *Runtime) serveGetPtr(req *getPtrMsg, conn sys.Conn) {
	defer conn.Close()

	h := r.exp.Lookup(req.ID)
	if h == nil {
		if err := conn.Send(&returnMsg{Err: NewError("getPtr: no exp handle")}); err != nil {
			// See comment in serveCall.
			if strings.HasPrefix(err.Error(), "gob") {
				panic(err)
			}
		}
		return
	}
	expReply, _ := r.exportValues([]interface{}{r.Ref(h.Value.Interface())}, conn.Addr())
	conn.Send(&returnMsg{Out: expReply})
}

func (r *Runtime) readGotPtrPtr(ptrPtr []*ptrPtrMsg, conn sys.Conn) error {
	p := make(map[circuit.HandleID]struct{})
	for _, pp := range ptrPtr {
		p[pp.ID] = struct{}{}
	}
	for len(p) > 0 {
		m_, err := conn.Receive()
		if err != nil {
			return err
		}
		m, ok := m_.(*gotPtrMsg)
		if !ok {
			return NewError("gotPtrMsg expected")
		}
		_, present := p[m.ID]
		if !present {
			return NewError("ack'ing unsent ptrPtrMsg")
		}
		delete(p, m.ID)
	}
	return nil
}
