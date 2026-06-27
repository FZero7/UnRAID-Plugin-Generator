package nut_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"smsups-nut-bridge/internal/mapper"
	"smsups-nut-bridge/internal/nut"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.dev")

	vars := mapper.NUTVars{
		"battery.charge": "100",
		"ups.status":     "OL",
		"input.voltage":  "233.4",
	}

	if err := nut.Write(path, vars); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range vars {
		line := k + ": " + v
		if !strings.Contains(string(content), line) {
			t.Errorf("expected %q in file, got:\n%s", line, content)
		}
	}
}

func TestWrite_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.dev")

	if err := nut.Write(path, mapper.NUTVars{"ups.status": "OB"}); err != nil {
		t.Fatal(err)
	}
	if err := nut.Write(path, mapper.NUTVars{"ups.status": "OL"}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "ups.status: OB") {
		t.Error("old content should be gone after overwrite")
	}
	if !strings.Contains(string(content), "ups.status: OL") {
		t.Error("new content not found")
	}
}

func TestWrite_AtomicOnSameDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.dev")

	err := nut.Write(path, mapper.NUTVars{"ups.status": "OL"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("file not created")
	}
}
