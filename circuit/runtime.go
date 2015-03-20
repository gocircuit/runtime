// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package circuit

import (
	"github.com/gocircuit/alef/ns"
)

type runtime interface {

	// Low-level
	ServerAddr() ns.Addr
	SetBoot(interface{})

	// Cross-services
	Dial(ns.Addr, string) PermX
	DialSelf(string) interface{}
	TryDial(ns.Addr, string) (PermX, error)
	Listen(string, interface{})

	// Persistence of PermX values
	Export(...interface{}) interface{}
	Import(interface{}) ([]interface{}, string, error)

	// Cross-interfaces
	Ref(interface{}) X
	PermRef(interface{}) PermX

	// Type system
	RegisterValue(interface{})

	// Utility
	Hang()
}
