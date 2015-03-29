// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package io

import (
	"io"
	"runtime"

	"github.com/gocircuit/core/circuit"
	"github.com/gocircuit/core/errors"
)

// Server-side types

// XReader is a cross-worker exportable object that exposes an underlying local io.Reader.
type XReader struct {
	io.Reader
}

func (x XReader) Read(n int) ([]byte, error) {
	p := make([]byte, n)
	m, err := x.Reader.Read(p)
	return p[:m], errors.Pack(err)
}

// XWriter is a cross-worker exportable object that exposes an underlying local io.Writer.
type XWriter struct {
	io.Writer
}

func (x XWriter) Write(p []byte) (int, error) {
	n, err := x.Writer.Write(p)
	return n, errors.Pack(err)
}

// XCloser is a cross-worker exportable object that exposes an underlying local io.Writer.
type XCloser struct {
	io.Closer
}

// NewXCloser attaches a finalizer to the object which calls Close.
// In cases when a cross-interface to this object is lost because of a failed remote worker,
// the attached finalizer will ensure that before we forget this object the channel it
// encloses will be closed.
func NewXCloser(u io.Closer) *XCloser {
	x := &XCloser{u}
	runtime.SetFinalizer(x, func(x *XCloser) {
		x.Closer.Close()
	})
	return x
}

func (x XCloser) Close() error {
	return errors.Pack(x.Closer.Close())
}

// XReadWriteCloser
type XReadWriteCloser struct {
	XReader
	XWriter
	*XCloser
}

func NewXReadWriteCloser(u io.ReadWriteCloser) circuit.X {
	return circuit.Ref(&XReadWriteCloser{XReader{u}, XWriter{u}, NewXCloser(u)})
}

// XReadCloser
type XReadCloser struct {
	XReader
	*XCloser
}

func NewXReader(u io.Reader) circuit.X {
	return circuit.Ref(XReader{u})
}

func NewXReadCloser(u io.ReadCloser) circuit.X {
	return circuit.Ref(&XReadCloser{XReader{u}, NewXCloser(u)})
}

// XWriteCloser
type XWriteCloser struct {
	XWriter
	*XCloser
}

func NewXWriteCloser(u io.WriteCloser) circuit.X {
	return circuit.Ref(&XWriteCloser{XWriter{u}, NewXCloser(u)})
}

// XReadWriter
type XReadWriter struct {
	XReader
	XWriter
}

func NewXReadWriter(u io.ReadWriter) circuit.X {
	return circuit.Ref(&XReadWriter{XReader{u}, XWriter{u}})
}
