// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// Package tele implements the circuit/use/n networking module using Teleport Transport
package tele

import (
	"net"
	"os"

	"github.com/gocircuit/alef/kit/tele"
	"github.com/gocircuit/alef/kit/tele/blend"
	"github.com/gocircuit/alef/peer"
)

// workerID is the ID for this transport endpoint.
// addr is the networking address to listen to.
func NewTransport(workerID peer.Id, addr net.Addr, key []byte) peer.Peer {
	var u *blend.Transport
	if len(key) == 0 {
		u = tele.NewStructOverTCP()
	} else {
		u = tele.NewStructOverTCPWithHMAC(key)
	}
	l := newListener(workerID, os.Getpid(), u.Listen(addr))
	return &Transport{
		Id: workerID,
		Dialer:   newDialer(l.Addr(), u),
		Listener: l,
	}
}

// Transport cumulatively represents the ability to listen for connections and dial into remote endpoints.
type Transport struct {
	peer.Id
	*Dialer
	*Listener
}
