// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package sys

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"

	"github.com/gocircuit/core/errors"
)

var ErrParse = errors.NewError("parse")

// Id represents the identity of a circuit worker process.
type Id string

// String returns a canonical string representation of this worker Id.
func (r Id) String() string {
	return string(r)
}

// ChooseId returns a random worker Id.
func ChooseId() Id {
	return Int64Id(rand.Int63())
}

func Int64Id(src int64) Id {
	return Id(fmt.Sprintf("Q%016x", src))
}

func UInt64Id(src uint64) Id {
	return Id(fmt.Sprintf("Q%016x", src))
}

// ParseOrHashId tries to parse the string s as a canonical worker Id representation.
// If it fails, it treats s as an unconstrained string and hashes it to a worker Id value.
// In either case, it returns a Id value.
func ParseOrHashId(s string) Id {
	id, err := ParseId(s)
	if err != nil {
		return HashId(s)
	}
	return id
}

// ParseId parses the string s for a canonical representation of a worker
// Id and returns a corresponding Id value.
func ParseId(s string) (Id, error) {
	if len(s) != 17 || s[0] != 'Q' {
		return "", ErrParse
	}
	ui64, err := strconv.ParseUint(s[1:], 16, 64)
	if err != nil {
		return "", err
	}
	return UInt64Id(ui64), nil
}

// HashId hashes the unconstrained string s into a worker Id value.
func HashId(s string) Id {
	h := fnv.New64a()
	h.Write([]byte(s))
	return UInt64Id(h.Sum64())
}
