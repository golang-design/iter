// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package iter_test

import (
	"testing"

	"golang.design/x/iter"
)

func TestNewBatchIterFromSlice(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		batchSize int
	}{
		{
			name:      "empty slice",
			slice:     []int{},
			batchSize: 1,
		},
		{
			name: "non-empty slice with batch size 1",
			slice: []int{
				1, 2, 3, 4, 5,
			},
			batchSize: 1,
		},
		{
			name: "non-empty slice with batch size 2",
			slice: []int{
				1, 2, 3, 4, 5,
			},
			batchSize: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := iter.NewBatchFromSlice(tt.slice, tt.batchSize)
			var got []int
			for batch, ok := it.Next(); ok; batch, ok = it.Next() {
				got = append(got, batch...)
			}
			if !sliceEqual(got, tt.slice) {
				t.Errorf("NewBatchIterFromSlice() = %v, want %v", got, tt.slice)
			}
		})
	}
}

func sliceEqual[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
