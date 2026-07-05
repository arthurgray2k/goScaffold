package core

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Alice")
	want := "Hello, Alice!"
	if got != want {
		t.Errorf("Hello(\"Alice\") = %q; want %q", got, want)
	}
}
