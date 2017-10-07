package main

import (
	"reflect"
	"testing"
)

func TestLabelFilter(t *testing.T) {
	data := Labels{
		{Name: "hacktoberfest", Color: "#ffffff"},
		{Name: "help wanted", Color: "#ffffff"},
		{Name: "test", Color: "#ffffff"},
	}

	want := map[string]string{"test": "#ffffff"}
	got := labelFilter(data)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("labelFilter(%v) failed", data)
		t.Errorf("got %v", got)
		t.Errorf("want %v", want)
	}
}
