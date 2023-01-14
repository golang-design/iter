// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package iter

// NewBatchFromSlice returns a batch iterator over a slice.
// Note that the slice is not copied and the caller must not modify it.
// To use the returned iterator, here is an example:
//
//	it := iter.NewBatchFromSlice[T](s, 42)
//	for batch, ok := it.Next(); ok; batch, ok = it.Next() {
//		// do something with batch
//	}
func NewBatchFromSlice[T any](s []T, batchSize int) Iter[[]T] {
	return &sliceIter[T]{s: s, batchSize: batchSize}
}

type sliceIter[T any] struct {
	s         []T
	batchSize int
	i         int
}

// Next implements Iter[T] and returns a batch of elements.
func (si *sliceIter[T]) Next() (ss []T, ok bool) {
	if si.i >= len(si.s) {
		return []T{}, false
	}
	end := si.i + si.batchSize
	if end > len(si.s) {
		end = len(si.s)
	}
	ss = si.s[si.i:end]
	si.i = end
	ok = true
	return
}

// ToSlice returns a slice containing all the elements in an iterator.
func BatchToSlice[E any](it Iter[[]E]) []E {
	var r []E
	for batch, ok := it.Next(); ok; batch, ok = it.Next() {
		for i := range batch {
			r = append(r, batch[i])
		}
	}
	return r
}
