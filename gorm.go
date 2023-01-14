// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package iter

import (
	"errors"
	"sync/atomic"

	"gorm.io/gorm"
)

// errStop is the error returned by the database iterator when
// the iteration is stopped.
var errStop = errors.New("iter: stop the iteration")

// GormIter is an gorm.DB compatible database iterator.
// To use this iterator, for example:
//
//	it := NewGormIter[T](tx, batchSize)
//	for batch, ok := it.Next(); ok; batch, ok = it.Next() {
//		// Process the batch.
//		...
//		// Stop the iteration if necessary.
//		if ... {
//			it.Stop()
//			break
//		}
//	}
//	if err := it.Err(); err != nil {
//		// Handle the error.
//	}
//
// This iterator is not safe to use after Err() returns a non-nil error.
// This iterator is not safe to use after Next() returns false.
// This iterator is not safe to use after Stop() is called.
// This iterator is not safe to use after the underlying database
// connection is closed.
type GormIter[T any] struct {
	tx *gorm.DB

	batchSize int
	next      chan chan []T
	stop      chan struct{}
	finished  atomic.Bool
	err       chan error
}

// NewBatchFromGorm creates a new database iterator with the given batch size.
func NewBatchFromGorm[T any](tx *gorm.DB, batchSize int) *GormIter[T] {
	it := &GormIter[T]{
		tx: tx,

		batchSize: batchSize,
		next:      make(chan chan []T),
		stop:      make(chan struct{}),
		err:       make(chan error, 1),
	}
	go it.batchFinder()
	return it
}

func (it *GormIter[T]) batchFinder() {
	var current []T
	err := it.tx.FindInBatches(&current, it.batchSize, func(tx *gorm.DB, _ int) error {
		rows := make([]T, len(current))
		copy(rows, current)

		select {
		case <-it.stop:
			return errStop
		case ch := <-it.next:
			ch <- rows
			close(ch)
		}
		return nil
	}).Error
	it.err <- err
	close(it.err)
	it.Stop()
}

// Next implements StopErrIter[T].
func (it *GormIter[T]) Next() ([]T, bool) {
	// Proceed to the next batch.
	done := make(chan []T)
	select {
	case <-it.stop:
		return nil, false
	case it.next <- done:
	}

	// Wait for the next batch.
	rows := <-done
	return rows, true
}

// Stop implements StopErrIter[T].
func (it *GormIter[T]) Stop() {
	for {
		finished := it.finished.Load()
		if finished {
			return
		}

		// If the swap success, then we can close the channel and return.
		if it.finished.CompareAndSwap(finished, true) {
			close(it.stop)
			return
		}
	}
}

// Err implements StopErrIter[T].
func (it *GormIter[T]) Err() error {
	err := <-it.err
	if err != nil && err != errStop {
		return err
	}
	return nil
}
