// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

package iter_test

import (
	"os"
	"strconv"
	"testing"

	"golang.design/x/iter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type user struct {
	gorm.Model
	Name string
	Age  int64
}

func TestGormIter(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	t.Cleanup(func() {
		db.Migrator().DropTable(&user{})
		os.Remove("test.db")
	})

	db.AutoMigrate(&user{})

	var want []user
	for i := int64(0); i < 2000; i++ {
		u := user{Name: strconv.FormatInt(i, 10), Age: i}
		want = append(want, u)
		db.Create(&u)
	}

	it := iter.NewBatchFromGorm[user](db, 1<<10)
	users := iter.BatchToSlice[user](it)
	if len(users) != len(want) {
		t.Fatalf("got %d users, want %d", len(users), len(want))
	}

	for i := range users {
		if users[i].Name != want[i].Name {
			t.Fatalf("got %s, want %s", users[i].Name, want[i].Name)
		}
		if users[i].Age != want[i].Age {
			t.Fatalf("got %d, want %d", users[i].Age, want[i].Age)
		}
	}
}

func BenchmarkGormIter(b *testing.B) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	b.Cleanup(func() {
		db.Migrator().DropTable(&user{})
		os.Remove("test.db")
	})

	db.AutoMigrate(&user{})
	for i := int64(0); i < 2000; i++ {
		u := user{Name: strconv.FormatInt(i, 10), Age: i}
		db.Create(&u)
	}

	b.Run("GormIter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			it := iter.NewBatchFromGorm[user](db, 1<<10)
			for batch, ok := it.Next(); ok; batch, ok = it.Next() {
				_ = batch
			}
		}
	})
}
