// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goexperiment.rangefunc

package iter_test

import (
	"testing"

	"golang.design/x/iter"
)

func TestBatchSlice(t *testing.T) {
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
			var got []int
			for i, batch := range iter.SliceBatch(tt.slice, tt.batchSize) {
				got = append(got, batch...)
				t.Log(i)
			}
			if !sliceEqual(got, tt.slice) {
				t.Errorf("SliceBatch() = %v, want %v", got, tt.slice)
			}
		})
	}
}
