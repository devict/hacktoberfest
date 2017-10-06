package main

import (
	"reflect"
	"testing"
)

func TestTop(t *testing.T) {
	data := map[string]int{
		"apple":      1234,
		"banana":     234,
		"cherry":     999,
		"date":       423,
		"elderberry": 1234,
	}

	want := []string{"apple", "elderberry", "cherry"}
	wantAlt := []string{"elderberry", "apple", "cherry"}

	got := top(3, data)

	// Don't panic if we have fewer values than 
	_ = top(33, data)

	if !reflect.DeepEqual(got, want) && !reflect.DeepEqual(got, wantAlt) {
		t.Errorf("top(%d, %v) failed", 3, data)
		t.Errorf("got %v", got)
		t.Errorf("want %v", want)
		t.Errorf("or want %v", wantAlt)
	}
}