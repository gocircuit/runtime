// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package lang

import (
	//"fmt"
	//"runtime/debug"

	"github.com/gocircuit/runtime/circuit"
	"github.com/gocircuit/runtime/lang/types"
	"github.com/gocircuit/runtime/sys"
)

func (r *Runtime) Listen(service string, receiver interface{}) {
	if IsX(receiver) {
		panic("listen service receiver cannot be a cross-interface")
	}
	types.RegisterValue(receiver)
	r.srv.Add(service, receiver)
}

// Dial returns an ptr to the permanent xvalue of the addressed remote runtime.
// It panics if any errors get in the way.
func (r *Runtime) Dial(addr sys.Addr, service string) circuit.PermX {
	if addr == nil {
		return nil
	}
	ptr, err := r.TryDial(addr, service)
	if err != nil {
		panic(err)
	}
	return ptr
}

// TryDial returns an ptr to the permanent xvalue of the addressed remote runtime
func (r *Runtime) TryDial(addr sys.Addr, service string) (circuit.PermX, error) {
	conn, err := r.t.Dial(addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	retrn, err := writeReturn(conn, &dialMsg{Service: service})
	if err != nil {
		return nil, err
	}

	return r.importEitherPtr(retrn, addr)
}

func (r *Runtime) DialSelf(service string) interface{} {
	return r.srv.Get(service)
}

func (r *Runtime) serveDial(req *dialMsg, conn sys.Conn) {
	// Go guarantees the defer runs even if panic occurs
	defer conn.Close()

	expDial, _ := r.exportValues([]interface{}{PermRef(r.srv.Get(req.Service))}, conn.Addr())
	conn.Send(&returnMsg{Out: expDial})
	// Waiting for export acks not necessary since expDial is always a permptr.
}

// Utils

func writeReturn(conn sys.Conn, msg interface{}) ([]interface{}, error) {
	if err := conn.Send(msg); err != nil {
		return nil, err
	}
	reply, err := conn.Receive()
	if err != nil {
		return nil, err
	}
	retrn, ok := reply.(*returnMsg)
	if !ok {
		return nil, NewError("foreign return type")
	}
	if retrn.Err != nil {
		return nil, err
	}
	return retrn.Out, nil
}

func (r *Runtime) importEitherPtr(retrn []interface{}, exporter sys.Addr) (circuit.PermX, error) {
	//debug.PrintStack()
	//println(fmt.Sprintf("retrn=%v exporter=%v", retrn, exporter))
	out, err := r.importValues(retrn, nil, exporter, false, nil)
	if err != nil {
		return nil, err
	}
	if len(out) != 1 {
		return nil, NewError("unexpected return value count")
	}
	if out[0] == nil {
		return nil, nil
	}
	ptr, ok := out[0].(circuit.PermX)
	if !ok {
		return nil, NewError("value is not a permanent cross-interface")
	}
	// XXX: Shouldn't this also have a case for non-permanent crossreferences?
	return ptr, nil
}
