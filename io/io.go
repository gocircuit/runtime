// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package io facilitates sharing io reader/writer/closers across workers.
package io

import (
	"runtime"
	"time"

	"github.com/gocircuit/core/circuit"
)

func init() {
	// If we are passing I/O objects cross-worker, we want to ensure that the GC
	// is activated regularly so that reclaimed I/O objects will close their
	// underlying resources in a timely manner.
	go func() {
		for {
			time.Sleep(5 * time.Second)
			runtime.GC()
		}
	}()

	//circuit.RegisterValue(XReader{})
	//circuit.RegisterValue(XWriter{})
	//circuit.RegisterValue(&XCloser{})
	//
	circuit.RegisterValue(&XReadCloser{})
	circuit.RegisterValue(&XWriteCloser{})
	circuit.RegisterValue(&XReadWriteCloser{})
	//
	circuit.RegisterValue(&XReadWriter{})
}

func unpackBytes(x interface{}) []byte {
	if x == nil {
		return nil
	}
	return x.([]byte)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
