package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeApp(t *testing.T, dir, marker string) string {
	t.Helper()
	app := filepath.Join(dir, "Discord.app")
	if err := os.MkdirAll(app, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, "marker"), []byte(marker), 0o644); err != nil {
		t.Fatal(err)
	}
	return app
}

func readMarker(t *testing.T, app string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(app, "marker"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestSwapBundleReplacesLiveAndRemovesOld(t *testing.T) {
	root := t.TempDir()
	live := writeApp(t, filepath.Join(root, "live"), "old")
	staged := writeApp(t, filepath.Join(root, "staged"), "new")
	if err := swapBundle(staged, live); err != nil {
		t.Fatal(err)
	}
	if got := readMarker(t, live); got != "new" {
		t.Fatalf("live marker = %q, want new", got)
	}
	if _, err := os.Stat(live + ".old"); !os.IsNotExist(err) {
		t.Fatalf("old bundle still present: %v", err)
	}
}

func TestSwapBundleRestoresLiveWhenStagedMissing(t *testing.T) {
	root := t.TempDir()
	live := writeApp(t, filepath.Join(root, "live"), "old")
	if err := swapBundle(filepath.Join(root, "missing.app"), live); err == nil {
		t.Fatal("expected an error")
	}
	if got := readMarker(t, live); got != "old" {
		t.Fatalf("live marker = %q, want old restored", got)
	}
}
