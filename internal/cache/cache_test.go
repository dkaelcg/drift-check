package cache

import (
	"os"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir, 5*time.Minute)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	entry := Entry{
		ResourceID: "aws_instance.web",
		Attributes: map[string]string{"instance_type": "t3.micro"},
	}
	if err := c.Set(entry); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := c.Get("aws_instance.web")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("expected entry, got nil")
	}
	if got.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("attribute mismatch: got %q", got.Attributes["instance_type"])
	}
}

func TestCache_Miss(t *testing.T) {
	dir := t.TempDir()
	c, _ := New(dir, 5*time.Minute)

	got, err := c.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestCache_Expired(t *testing.T) {
	dir := t.TempDir()
	c, _ := New(dir, 1*time.Millisecond)

	_ = c.Set(Entry{ResourceID: "res1", Attributes: map[string]string{"k": "v"}})
	time.Sleep(5 * time.Millisecond)

	got, err := c.Get("res1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Error("expected nil for expired entry")
	}
}

func TestCache_Invalidate(t *testing.T) {
	dir := t.TempDir()
	c, _ := New(dir, 5*time.Minute)

	_ = c.Set(Entry{ResourceID: "res2", Attributes: map[string]string{}})
	if err := c.Invalidate("res2"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}

	got, _ := c.Get("res2")
	if got != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestCache_InvalidateMissing(t *testing.T) {
	dir := t.TempDir()
	c, _ := New(dir, 5*time.Minute)

	if err := c.Invalidate("does-not-exist"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNew_BadDir(t *testing.T) {
	// Use a file path as a directory to force failure.
	f, _ := os.CreateTemp("", "cache-test")
	defer os.Remove(f.Name())
	f.Close()

	_, err := New(f.Name()+"/subdir", time.Minute)
	if err == nil {
		t.Error("expected error for bad dir, got nil")
	}
}
