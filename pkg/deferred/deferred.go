// Package deferred has machinery for providing JavaScript with
// deferred values. These arise because we want to have library
// functions for JavaScript to request things that will only later be
// supplied -- either because it will take some time to get, and it's
// better if JavaScript doesn't have to block (e.g., an HTTP request),
// or because we're waiting for something elsewhere to happen (e.g.,
// watching a resource for changes).
//
// In JavaScript, deferred values will often end up being represented
// by promises; but the protocol allows other representations, and in
// particular, some requests may result in an sequence of values
// rather than just a single value. The protocol also allows for
// cancelling a deferred value.
package deferred

import (
	"sync"
)

var (
	globalDeferreds = newDeferreds()
)

func newDeferreds() *deferreds {
	d := &deferreds{
		serial: 1,
	}
	d.serialCond = sync.NewCond(&d.serialMu)
	return d
}

// Register schedules an action to be performed later, with the result
// sent to `resolver`, using the global deferred scheduler.
func Register(p performFunc, r resolver) Serial {
	return globalDeferreds.Register(p, r)
}

// Wait blocks until all outstanding deferred values in the global
// scheduler are fulfilled.
func Wait() {
	globalDeferreds.Wait()
}

// Serial is a serial number used to identify deferreds between Go and
// JavaScript.
type Serial uint64

// To enforce determinism, we resolve deferred in the same order they are
// created. This is done through resolvedSerial that stores what was the last
// deferred resolved and we use a sync.Cond to handle synchronization between
// goroutines servicing the deferred.
type deferreds struct {
	serialMu       sync.Mutex
	serial         Serial
	serialCond     *sync.Cond
	resolvedSerial Serial

	outstanding sync.WaitGroup
}

func (d *deferreds) waitForSerial(s Serial) {
	d.serialMu.Lock()
	defer d.serialMu.Unlock()

	for {
		if d.resolvedSerial == s {
			return
		}
		d.serialCond.Wait()
	}
}

func (d *deferreds) serialResolved(s Serial) {
	d.serialMu.Lock()
	d.resolvedSerial = s
	d.serialMu.Unlock()
	d.serialCond.Broadcast()
}

// responder is the interface for a deferred request to use to send
// its response.
type resolver interface {
	Error(Serial, error)
	Data(Serial, []byte)
	End(Serial)
}

type performFunc func() ([]byte, error)

// Register adds a request to those being tracked, and returns the
// serial number to give back to the runtime.
func (d *deferreds) Register(perform performFunc, r resolver) Serial {
	d.serialMu.Lock()
	s := d.serial
	d.serial++
	d.serialMu.Unlock()
	d.outstanding.Add(1)
	go func(s Serial) {
		defer func() {
			d.serialResolved(s)
			d.outstanding.Done()
		}()

		b, err := perform()

		// Wait for the serial-1 goroutine to be resolved.
		d.waitForSerial(s - 1)

		if err != nil {
			r.Error(s, err)
			return
		}
		r.Data(s, b)
	}(s)
	return s
}

// Wait blocks until all outstanding deferred requests are fulfilled.
func (d *deferreds) Wait() {
	d.outstanding.Wait()
}
