package data

import (
	"os"
	"path/filepath"
	"testing"
)

func testFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}
	return path
}

func TestLoadCSV(t *testing.T) {
	path := testFile(t, "data.csv", "name,value\nAlice,10\nBob,20\nCharlie,30\n")

	ds, err := Load(path, LoadConfig{HasHeader: true})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(ds.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(ds.Rows))
	}
	if ds.Rows[0].Label != "Alice" || ds.Rows[0].Value != 10 {
		t.Errorf("first row: %+v", ds.Rows[0])
	}
}

func TestLoadCSVNoHeader(t *testing.T) {
	path := testFile(t, "data.csv", "10\n20\n30\n")

	ds, err := Load(path, LoadConfig{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(ds.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(ds.Rows))
	}
}

func TestLoadJSONMap(t *testing.T) {
	path := testFile(t, "data.json", `{"a": 1, "b": 2, "c": 3}`)

	ds, err := Load(path, LoadConfig{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(ds.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(ds.Rows))
	}
}

func TestLoadJSONArray(t *testing.T) {
	path := testFile(t, "data.json", `[{"label":"a","value":1},{"label":"b","value":2}]`)

	ds, err := Load(path, LoadConfig{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(ds.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(ds.Rows))
	}
	if ds.Rows[0].Label != "a" || ds.Rows[0].Value != 1 {
		t.Errorf("first row: %+v", ds.Rows[0])
	}
}

func TestLoadEmptyFile(t *testing.T) {
	path := testFile(t, "empty.csv", "")

	_, err := Load(path, LoadConfig{})
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	_, err := Load("/nonexistent/file.csv", LoadConfig{})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadJSONSimpleArray(t *testing.T) {
	path := testFile(t, "data.json", `[10, 20, 30]`)

	ds, err := Load(path, LoadConfig{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(ds.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(ds.Rows))
	}
	if ds.Rows[0].Value != 10 {
		t.Errorf("first value: %v", ds.Rows[0].Value)
	}
}

func TestReadInfo(t *testing.T) {
	path := testFile(t, "data.csv", "name,value\nAlice,10\nBob,20\n")

	info, err := ReadInfo(path)
	if err != nil {
		t.Fatalf("ReadInfo: %v", err)
	}
	if info.Format != "csv" {
		t.Errorf("expected csv, got %s", info.Format)
	}
	if info.Rows != 3 {
		t.Errorf("expected 3 rows, got %d", info.Rows)
	}
	if info.Header[0] != "name" {
		t.Errorf("expected header 'name', got '%s'", info.Header[0])
	}
}

func TestLoadCSVWithLabels(t *testing.T) {
	path := testFile(t, "data.csv", "category,sales\nElectronics,500\nClothing,300\nBooks,100\n")

	ds, err := Load(path, LoadConfig{HasHeader: true, LabelCol: "category", ValueCol: "sales"})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(ds.Rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(ds.Rows))
	}
	if ds.Rows[0].Label != "Electronics" || ds.Rows[0].Value != 500 {
		t.Errorf("first row: %+v", ds.Rows[0])
	}
}
