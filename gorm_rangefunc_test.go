// Copyright 2023 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build goexperiment.rangefunc

package iter_test

import (
	"os"
	"strconv"
	"testing"

	"golang.design/x/iter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormBatch(t *testing.T) {
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

	var users []user
	for i, batch := range iter.GormBatch[user](db, 1<<10) {
		users = append(users, batch...)
		t.Log(i)
	}
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
