// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goexperiment.rangefunc

package iter

// SliceBatch returns a function that corporate with range-over-func syntax
// to iterate over a slice in batches.
//
// Example:
//
//	for _, batch := range iter.SliceBatch(s, 1<<10) {
//		// ...
//	}
func SliceBatch[E any](s []E, batchSize int) func(func(int, []E) bool) {
	return func(yield func(int, []E) bool) {
		for i := 0; i < len(s); i += batchSize {
			end := i + batchSize
			if end > len(s) {
				end = len(s)
			}
			if !yield(i, s[i:end]) {
				return
			}
		}
	}
}
