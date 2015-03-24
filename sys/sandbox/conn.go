// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package sandbox

import (
	"github.com/gocircuit/core/sys"
	"io"
)

type conn struct {
	send chan<- interface{}
	recv <-chan interface{}
}

func NewPair() (c, d sys.Conn) {
	x, y := make(chan interface{}, 5), make(chan interface{}, 5)
	return &conn{x, y}, &conn{y, x}
}

func (c *conn) Receive() (interface{}, error) {
	v, ok := <-c.recv
	if !ok {
		return nil, io.ErrUnexpectedEOF
	}
	return v, nil
}

func (c *conn) Send(v interface{}) error {
	c.send <- v
	return nil
}

func (c *conn) Close() error {
	close(c.send)
	return nil
}

func (c *conn) Addr() sys.Addr {
	return nil
}
