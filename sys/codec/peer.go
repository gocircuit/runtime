// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"github.com/gocircuit/alef/sys"
)

func New(codec Codec, chunk sys.Peer) sys.Peer {
	return &peer{codec, chunk}
}

type peer struct {
	codec Codec
	chunk sys.Peer
}

func (p *peer) Accept() sys.Conn {
	return newConn(p.codec, p.chunk.Accept())
}

func (p *peer) Addr() sys.Addr {
	return p.chunk.Addr()
}

func (p *peer) Dial(addr sys.Addr) (conn sys.Conn, err error) {
	if conn, err = p.chunk.Dial(addr); err != nil {
		return nil, err
	}
	return newConn(p.codec, conn), nil
}
