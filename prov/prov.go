package prov

import (
	"github.com/gocircuit/runtime/circuit"
)

// Provisioner is a system that provisions new hosts and runs circuit workers on them.
type Provisioner interface {
	// Provision provisions n new hosts and launches workers on each of them.
	Provision(profile interface{}, n int) []Worker // Returns a list of circuit addresses

	// Reset stops all workers spawned with this provisioner.
	Reset()
}

// Worker is a computing resource that executes jobs.
type Worker interface {
	// Addr returns the circuit address of the worker.
	Addr() string

	// Spawn executes…
	Spawn(prog interface{}) circuit.X

	// Stop kills this worker and releases the underlying host.
	Stop()
}
