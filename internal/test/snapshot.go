package test

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Snapshot either writes or compares with a file in ./testdata
func Snapshot(t *testing.T, data any) {
	t.Helper()

	var actual []byte
	switch d := data.(type) {
	case string:
		actual = []byte(d)
	case []byte:
		actual = d
	default:
		a, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		actual = a
	}

	golden := filepath.Join(".", "testdata", t.Name())
	if _, err := os.Stat(golden); errors.Is(err, fs.ErrNotExist) {
		// generate new snapshot
		_ = os.MkdirAll(filepath.Dir(golden), 0o750)
		err := os.WriteFile(golden, actual, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	expected, err := os.ReadFile(filepath.Clean(golden))
	if err != nil {
		t.Fatal(err)
	}

	a, b := compareLines(expected, actual)
	if len(a) > 0 || len(b) > 0 {
		t.Logf("expected: %s", a)
		t.Logf("  actual: %s", b)
		t.Fatal("snapshot does not match")
	}
}

func compareLines(a, b []byte) ([]string, []string) {
	aLines := strings.Split(string(a), "\n")
	bLines := strings.Split(string(b), "\n")

	aLineMap := make(map[string]bool)
	bLineMap := make(map[string]bool)

	for _, line := range aLines {
		aLineMap[line] = true
	}
	for _, line := range bLines {
		bLineMap[line] = true
	}

	var uniqueA, uniqueB []string

	for line := range aLineMap {
		if !bLineMap[line] {
			uniqueA = append(uniqueA, line)
		}
	}

	for line := range bLineMap {
		if !aLineMap[line] {
			uniqueB = append(uniqueB, line)
		}
	}

	return uniqueA, uniqueB
}
