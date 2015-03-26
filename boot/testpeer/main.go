// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"flag"
	"github.com/gocircuit/core/boot"
	"github.com/gocircuit/core/circuit"
	"github.com/gocircuit/core/sys/tcp"
	"os"
)

var flagAddr = flag.String("addr", "", "Our address")
var flagDial = flag.String("dial", "", "Their address")

func main() {
	flag.Parse()
	if err := boot.BootTCP(*flagAddr); err != nil {
		println("boot error:", err.Error())
		os.Exit(1)
	}
	circuit.Listen("greet", &helloService{})
	if *flagDial != "" {
		a, err := tcp.ResolveAddr(*flagDial)
		if err != nil {
			println("dial resolve error:", err.Error())
			os.Exit(1)
		}
		x, err := circuit.TryDial(a, "greet")
		if err != nil {
			println("circuit dial error:", err.Error())
			os.Exit(1)
		}
		x.Call("Hello")
	}
	select {}
}

type helloService struct{}

func (s *helloService) Hello() {
	println("hello")
}

func init() {
	circuit.RegisterValue(&helloService{})
}
