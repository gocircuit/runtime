// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package pool

import (
	"github.com/gocircuit/core/sys"
)

type conn struct {
	addr sys.Addr
	sys.Conn
}

func newConn(under sys.Conn, addr sys.Addr) *conn {
	return &conn{addr, under}
}

func (c *conn) Addr() sys.Addr {
	return c.addr
}
