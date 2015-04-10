// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package backdial

import (
	"encoding/gob"
	"errors"
	"github.com/gocircuit/runtime/sys"
)

type Welcome struct {
	Back sys.Addr
}

func init() {
	gob.Register(&Welcome{})
}

type conn struct {
	back sys.Addr
	sys.Conn
}

func (p *peer) handshake(u sys.Conn) (*conn, error) {
	if err := u.Send(&Welcome{p.Addr()}); err != nil {
		u.Close()
		return nil, err
	}
	w, err := u.Receive()
	if err != nil {
		u.Close()
		return nil, err
	}
	welcome, ok := w.(*Welcome)
	if !ok {
		u.Close()
		return nil, errors.New("protocol")
	}
	return &conn{welcome.Back, u}, nil
}

func (c *conn) Addr() sys.Addr {
	return c.back
}
