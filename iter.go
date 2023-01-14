// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

// Package iter provides a generic iterator interface and utils.
//
// This is an implementation of the Go language iterator proposal.
// See https://go.dev/issue/54245 for more details.
package iter

// Iter supports iterating over a sequence of values of type `E`.
type Iter[E any] interface {
	// Next returns the next value in the iteration if there is one,
	// and reports whether the returned value is valid.
	// Once Next returns ok==false, the iteration is over,
	// and all subsequent calls will return ok==false.
	Next() (E, bool)
}

// StopIter is an optional interface for Iter that allows the caller
// to stop iteration early.
type StopIter[E any] interface {
	Iter[E]

	// Stop indicates that the iterator will no longer be used.
	// After a call to Stop, future calls to Next may panic.
	// Stop may be called multiple times;
	// all calls after the first will have no effect.
	Stop()
}

// StopErrIter is an optional interface for Iter that allows the caller
// to check for errors encountered during iteration.
type StopErrIter[E any] interface {
	StopIter[E]

	// Err returns the error, if any, that was encountered during iteration.
	Err() error
}
