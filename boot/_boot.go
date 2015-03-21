// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package boot

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/gocircuit/circuit/kit/debug/kill"
	"github.com/gocircuit/circuit/kit/lockfile"
	"github.com/gocircuit/circuit/sys/lang"
	_ "github.com/gocircuit/circuit/sys/tele"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

// BootTCP loads the runtime over unencrypted TCP transport.
func BootTCP(addr *net.TCPAddr) n.Addr {
	//debug.InstallCtrlCPanic()

	// Randomize execution
	rand.Seed(time.Now().UnixNano())

	// Generate worker ID
	id := n.ChooseId()
	vardir = strings.Replace(vardir, "%W", id.String(), 1)

	// Initialize networking
	if len(key) > 0 {
		log.Println("Using symmetric HMAC authentication and RC4 encryption.")
	}
	t := n.NewTransport(id, addr, key)
	fmt.Println(t.Addr().String())

	// Initialize language runtime
	circuit.Bind(lang.New(t))
	return t.Addr()
}
