package utils_test

import (
	"testing"

	"github.com/vitor-mariano/regex-tui/pkg/utils"
)

func TestNewSet(t *testing.T) {
	s := utils.NewSet("a", "b", "a")
	if s.Size() != 2 {
		t.Fatalf("Size() = %d, want 2", s.Size())
	}
	if !s.Contains("a") || !s.Contains("b") {
		t.Fatal("NewSet should contain initial items")
	}
	if s.Contains("c") {
		t.Fatal("NewSet should not contain items that were not added")
	}
}

func TestSetAddContainsRemove(t *testing.T) {
	s := utils.NewSet[string]()
	if s.Size() != 0 {
		t.Fatalf("empty Size() = %d, want 0", s.Size())
	}

	s.Add("x", "y")
	if !s.Contains("x") || !s.Contains("y") {
		t.Fatal("Add should insert items")
	}
	if s.Size() != 2 {
		t.Fatalf("Size() = %d, want 2", s.Size())
	}

	s.Add("x")
	if s.Size() != 2 {
		t.Fatalf("duplicate Add changed Size() to %d", s.Size())
	}

	s.Remove("x")
	if s.Contains("x") {
		t.Fatal("Remove should drop the item")
	}
	if !s.Contains("y") {
		t.Fatal("Remove should leave other items")
	}
	if s.Size() != 1 {
		t.Fatalf("Size() = %d, want 1", s.Size())
	}

	s.Remove("missing")
	if s.Size() != 1 {
		t.Fatalf("Remove of missing item changed Size() to %d", s.Size())
	}
}
