package main

import "testing"

func TestEcho(t *testing.T) {
	var tests = []struct {
		newline bool
		sep string
		args []string
		want string
	} {
		{true, "", []string{}, "\n"},
		{false, "", []string{}, ""},
		{true, "\t", []string{"one", "two", "three"}, "one\ttwo\tthree\n"},
		{true, ",", []string{"a", "b", "c"}, "a, b, c\n"},
		{false, ":", []string{"1", "2", "3"}, "1:2:3"},
	}
}
