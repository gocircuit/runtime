// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package net

import (
	"github.com/gocircuit/core/sys"
	"io"
)

type Peer interface {
	Dial(addr sys.Addr) (Conn, error)
	Accept() (Conn, error)
	Addr() sys.Addr
}

type Conn interface {
	io.ReadWriteCloser
	Addr() sys.Addr
}
