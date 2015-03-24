// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

// Package tcp implements a peer transport over TCP.
package tcp

import (
	"bufio"
	"encoding/binary"
	"github.com/gocircuit/core/sys"
	"net"
)

// conn implements sys.Conn
type conn struct {
	tcp *net.TCPConn
	r   *bufio.Reader
}

func newConn(c *net.TCPConn) *conn {
	if err := c.SetKeepAlive(true); err != nil {
		panic(err)
	}
	return &conn{c, bufio.NewReader(c)}
}

func (c *conn) Addr() sys.Addr {
	return nil // cannot determine peer address of remote peer
}

func (c *conn) Read() (chunk interface{}, err error) {
	k, err := binary.ReadUvarint(c.r)
	if err != nil {
		return nil, err
	}
	var q = make([]byte, k)
	var n, m int
	for m < len(q) && err == nil {
		n, err = c.r.Read(q[m:])
		m += n
	}
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (c *conn) Write(v interface{}) (err error) {
	chunk := v.([]byte)
	q := make([]byte, len(chunk)+8)
	n := binary.PutUvarint(q, uint64(len(chunk)))
	m := copy(q[n:], chunk)
	_, err = c.tcp.Write(q[:n+m])
	return err
}

func (c *conn) Close() (err error) {
	return c.tcp.Close()
}
