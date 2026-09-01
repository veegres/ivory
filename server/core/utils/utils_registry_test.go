package utils

import (
	"errors"
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry[string, int]()
	r.Register("a", 1)

	got, err := r.Get("a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewRegistry[string, int]()

	_, err := r.Get("missing")
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("got %v, want ErrNotImplemented", err)
	}
}

func TestRegistry_All_ReturnsACopy(t *testing.T) {
	r := NewRegistry[string, int]()
	r.Register("a", 1)

	all := r.All()
	all["a"] = 99
	all["b"] = 2

	got, err := r.Get("a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != 1 {
		t.Errorf("mutating All()'s result changed the registry: got %d, want 1", got)
	}
	if _, err := r.Get("b"); err == nil {
		t.Error("adding to All()'s result added an entry to the registry")
	}
}
