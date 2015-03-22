package dependency

import (
	"github.com/juju/errors"

	"github.com/juju/juju/worker"
)

// GetResourceFunc returns an indication of whether the named manifold can satisfy
// a dependency. In particular:
//  * if the named manifold does not exist, it returns false
//  * if the named manifold exists and out is nil, it returns true
//  * if the named manifold exists and out is non-nil, it returns whether
//    the named manifold's worker was able to assign a suitable value to the
//    out pointer. Appropriate types for the out pointer depend upon the
//    resources/services exposed by the worker in question.
type GetResourceFunc func(name string, out interface{}) bool

// StartFunc returns a worker or an error. All the worker's dependencies should
// be taken from the supplied Registry; if no worker can be started, a StartFunc
// should return an error rather than waiting for its dependencies to become
// available.
type StartFunc func(getResource GetResourceFunc) (worker.Worker, error)

// OutputFunc is a type coercion function for a worker generated by a StartFunc.
// When passed an out pointer to a type it recognises, it will assign a suitable
// value and return true on success.
type OutputFunc func(in worker.Worker, out interface{}) bool

// Manifold defines the behaviour of a node in an Engine's dependency graph.
type Manifold struct {

	// Inputs lists the names of the manifolds which this manifold might use.
	Inputs []string

	// Start is used to create a worker for the manifold. It must not be nil.
	Start StartFunc

	// Output is used to implement a GetResourceFunc for manifolds that declare
	// a dependency on this one; it can be nil if your manifold is a leaf node,
	// or if it exposes no services to its dependents.
	Output OutputFunc
}

// Engine runs the worker for every installed manifold, restarting them when they
// error out, and as their inputs are started or stopped. Workers that exit with
// no error will not be restarted until their dependencies change.
type Engine interface {

	// Engine is just another Worker.
	worker.Worker

	// Install causes the Engine to accept responsibility for maintaining a
	// worker corresponding to the supplied manifold, restarting it when it
	// fails and when its inputs' workers change, until the Engine shuts down.
	Install(name string, manifold Manifold) error
}

// IsFatalFunc is used to configure an Engine such that, if any worker returns
// an error that satisfies the engine's IsFatalFunc, the engine will stop all
// its workers, shut itself down, and return the original fatal error via Wait().
type IsFatalFunc func(err error) bool

// ErrUnmetDependencies can be returned by a StartFunc or a worker to indicate to
// the engine that it can't be usefully restarted until at least one of its
// dependencies changes. There's no way to specify *which* dependency you need,
// because that's a lot of implementation hassle for little practicall gain.
var ErrUnmetDependencies = errors.New("cannot run with available dependencies")
