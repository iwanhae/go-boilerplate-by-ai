package store

import "testing"

func TestMemorySetGet(t *testing.T) {
	s := NewMemory()
	if err := s.Set("a", 1); err != nil {
		t.Fatalf("set: %v", err)
	}
	v, err := s.Get("a")
	if err != nil || v.(int) != 1 {
		t.Fatalf("get: %v %v", v, err)
	}
	if _, err := s.Get("b"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
}

func TestMemoryListAndDelete(t *testing.T) {
	s := NewMemory()
	s.Set("p1", 1)
	s.Set("p2", 2)
	vals, _ := s.List("p")
	if len(vals) != 2 {
		t.Fatalf("list: %v", vals)
	}
	if err := s.Delete("p1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := s.Delete("p1"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	vals, _ = s.List("")
	if len(vals) != 1 {
		t.Fatalf("list after delete: %v", vals)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}
