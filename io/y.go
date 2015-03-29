// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package io facilitates sharing io reader/writer/closers across workers.
package io

import (
	"fmt"
	"io"
	"os"

	"github.com/gocircuit/core/circuit"
	"github.com/gocircuit/core/errors"
)

// Client-side types

// YReader
type YReader struct {
	circuit.X
}

func (y YReader) Read(p []byte) (n int, err error) {
	defer func() {
		// 	println(fmt.Sprintf("yread n=%d err=%v r=%v", n, err, recover()))
		if r := recover(); r != nil {
			println(fmt.Sprintf("r=%v", r))
			os.Exit(1)
		}
	}()
	r := y.Call("Read", len(p))
	q, err := unpackBytes(r[0]), errors.Unpack(r[1])
	if len(q) > len(p) {
		panic("corrupt i/o server")
	}
	copy(p, q)
	if err != nil && err.Error() == "EOF" {
		err = io.EOF
	}
	return len(q), err
}

// YWriter
type YWriter struct {
	circuit.X
}

func (y YWriter) Write(p []byte) (n int, err error) {
	r := y.Call("Write", p)
	return r[0].(int), errors.Unpack(r[1])
}

// YCloser
type YCloser struct {
	circuit.X
}

func (y YCloser) Close() error {
	return errors.Unpack(y.Call("Close")[0])
}

// YReadCloser
type YReadCloser struct {
	YReader
	YCloser
}

func NewYReader(u interface{}) YReader {
	return YReader{u.(circuit.X)}
}

func NewYReadCloser(u interface{}) *YReadCloser {
	return &YReadCloser{YReader{u.(circuit.X)}, YCloser{u.(circuit.X)}}
}

// YWriteCloser
type YWriteCloser struct {
	YWriter
	YCloser
}

func NewYWriteCloser(u interface{}) *YWriteCloser {
	return &YWriteCloser{YWriter{u.(circuit.X)}, YCloser{u.(circuit.X)}}
}

// YReadWriteCloser
type YReadWriteCloser struct {
	YReader
	YWriter
	YCloser
}

func NewYReadWriteCloser(u interface{}) *YReadWriteCloser {
	return &YReadWriteCloser{YReader{u.(circuit.X)}, YWriter{u.(circuit.X)}, YCloser{u.(circuit.X)}}
}

// YReadWriter
type YReadWriter struct {
	YReader
	YWriter
}

func NewYReadWriter(u interface{}) *YReadWriter {
	return &YReadWriter{YReader{u.(circuit.X)}, YWriter{u.(circuit.X)}}
}
