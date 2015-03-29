// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package circuit exposes the core functionalities provided by the circuit programming environment
package circuit

import (
	"github.com/gocircuit/core/circuit/module"
	"github.com/gocircuit/core/lang/types"
	"github.com/gocircuit/core/sys"
)

var mod = module.Slot{Name: "language"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v runtime) {
	mod.Set(v)
}

func get() runtime {
	return mod.Get().(runtime)
}

// Operators

// RegisterValue registers the type of v with the circuit runtime type system.
// As a result this program becomes able to send and receive cross-interfaces pointing to objects of this type.
// By convention, RegisterValue should be invoked from a dedicated init
// function within of the package that defines the type of v.
func RegisterValue(v interface{}) {
	types.RegisterValue(v)
}

// Ref returns a cross-interface to the local value v.
func Ref(v interface{}) X {
	return get().Ref(v)
}

// PermRef returns a permanent cross-interface to the local value v.
func PermRef(v interface{}) PermX {
	return get().PermRef(v)
}

// ServerAddr returns the address of this worker.
func ServerAddr() sys.Addr {
	return get().ServerAddr()
}

func setBoot(v interface{}) {
	get().SetBoot(v)
}

// Dial contacts the worker specified by addr and requests a cross-worker
// interface to the named service.
// If service is not being listened to at this worker, nil is returned.
// Failures to contact the worker for external/physical reasons result in a
// panic.
func Dial(addr sys.Addr, service string) PermX {
	return get().Dial(addr, service)
}

// DialSelf works similarly to Dial, except it dials into the calling worker
// itself and instead of returning a cross-interface to the service receiver,
// it returns a native Go interface. DialSelf never fails.
func DialSelf(service string) interface{} {
	return get().DialSelf(service)
}

// Listen registers the receiver object as a receiver for the named service.
// Subsequent calls to Dial from other works, addressing this worker and the
// same service name, will return a cross-interface to receiver.
func Listen(service string, receiver interface{}) {
	get().Listen(service, receiver)
}

// TryDial behaves like Dial, with the difference that instead of panicking in
// the event of external/physical issues, an error is returned instead.
func TryDial(addr sys.Addr, service string) (PermX, error) {
	return get().TryDial(addr, service)
}

// Export recursively rewrites the values val into a Go type that can be
// serialiazed with package encoding/gob. The values val can contain permanent
// cross-interfaces (but no non-permanent ones).
func Export(val ...interface{}) interface{} {
	return get().Export(val...)
}

// Import converts the exported value, that was produced as a result of Export,
// back into its original form.
func Import(exported interface{}) ([]interface{}, string, error) {
	return get().Import(exported)
}

// Hang never returns
func Hang() {
	get().Hang()
}
