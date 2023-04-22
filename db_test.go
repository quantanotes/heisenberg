package main

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	path, _ := GetDir("/tests/db.db")
	db, err := NewDB(path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	p1 := pair{
		K: "a",
		V: value{Vec: []float32{1, 2, 3}, Meta: []byte("hello")},
	}

	p2 := pair{
		K: "b",
		V: value{Vec: []float32{2, 3, 4}, Meta: []byte("world")},
	}

	err = db.NewCollection("test")
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = db.Put(p1, "test")
	if err != nil {
		panic(err)
	}

	err = db.Put(p2, "test")
	if err != nil {
		panic(err)
	}

	r1, err := db.Get("a", "test")
	if err != nil {
		panic(err)
	}

	r2, err := db.Get("b", "test")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(r1.V.Meta.([]byte)))
	fmt.Println(string(r2.V.Meta.([]byte)))
}
