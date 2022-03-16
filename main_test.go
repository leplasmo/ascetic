package main_test

import "testing"

func TestMain(t *testing.T) {
	got := 1
	want := 1
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
