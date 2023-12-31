// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goexperiment.rangefunc

package iter

import "gorm.io/gorm"

// GormBatch returns a function that corporate with range-over-func syntax
// to iterate over a slice in batches.
//
// Example:
//
//	var users []user
//	for i, batch := range iter.GormBatch[user](db, 1<<10) {
//		users = append(users, batch...)
//	}
func GormBatch[E any](tx *gorm.DB, batchSize int) func(func(int, []E) bool) {
	it := NewBatchFromGorm[E](tx, batchSize)
	return func(yield func(int, []E) bool) {
		defer it.Stop()

		n := 0
		for batch, ok := it.Next(); ok; batch, ok = it.Next() {
			if !yield(n, batch) {
				return
			}
			n += len(batch)
		}
	}
}
