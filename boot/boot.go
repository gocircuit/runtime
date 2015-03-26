// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package boot

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gocircuit/core/circuit"
	_ "github.com/gocircuit/core/kit/debug/kill"
	"github.com/gocircuit/core/lang"
	sysuse "github.com/gocircuit/core/sys/use"
)

// BootTCP loads the runtime over unencrypted TCP transport.
func BootTCP(addr string) error {

	// Randomize execution
	rand.Seed(time.Now().UnixNano())

	// Load peer networking
	t, err := sysuse.NewClearTCP(addr)
	if err != nil {
		return err
	}
	fmt.Println(t.Addr().String())

	// Initialize language runtime
	circuit.Bind(lang.New(t))

	return nil
}
