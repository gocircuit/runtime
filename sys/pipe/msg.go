// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package pipe

import (
	"encoding/gob"
)

type (
	PipeId int
	SeqNo  uint32
)

type PayloadMsg struct {
	SeqNo   SeqNo
	Payload interface{} // User-supplied type; must be registered with encoding/gob
}

type AbortMsg struct {
	Reason error
}

type Msg struct {
	PipeId PipeId
	Msg    interface{} // PayloadMsg or AbortMsg or nil (introduces a new pipe id)
}

func init() {
	gob.Register(&Msg{})
	gob.Register(&PayloadMsg{})
	gob.Register(&AbortMsg{})
}
