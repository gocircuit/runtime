// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package codec

import (
	"github.com/gocircuit/core/sys"
)

type conn struct {
	enc Encoder
	dec Decoder
	u   sys.Conn
}

func newConn(codec Codec, u sys.Conn) *conn {
	return &conn{
		enc: codec.NewEncoder(),
		dec: codec.NewDecoder(),
		u:   u,
	}
}

func (c *conn) Addr() sys.Addr {
	return c.u.Addr()
}

type frame struct {
	Value interface{}
}

func (c *conn) Send(v interface{}) (err error) {
	chunk, err := c.enc.Encode(frame{v})
	if err != nil {
		return err
	}
	return c.u.Send(chunk)
}

func (c *conn) Receive() (v interface{}, err error) {
	chunk, err := c.u.Receive()
	if err != nil {
		return nil, err
	}
	var f frame
	if err = c.dec.Decode(chunk.([]byte), &f); err != nil {
		return nil, err
	}
	return f.Value, nil
}

func (c *conn) Close() (err error) {
	return c.u.Close()
}
