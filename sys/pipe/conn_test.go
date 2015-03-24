// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"github.com/gocircuit/core/sys"
	"github.com/gocircuit/core/sys/sandbox"
	"testing"
)

func TestConn(t *testing.T) {
	var c [2]sys.Conn
	var d [2]*Conn
	c[0], c[1] = sandbox.NewPair()
	d[0], d[1] = newConn(c[0], 1), newConn(c[1], -1)

	go func() {
		p, err := d[0].Dial()
		if err != nil {
			t.Fatalf("dial %v", err)
		}
		if err = p.Send(1); err != nil {
			t.Fatalf("write %v", err)
		}
		p.Close()
	}()

	p, err := d[1].Accept()
	if err != nil {
		t.Fatalf("accept %v", err)
	}
	x, err := p.Receive()
	if err != nil {
		t.Fatalf("read %v", err)
	}
	println(x.(int))
	p.Close()
}
