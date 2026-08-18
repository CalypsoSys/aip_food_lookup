package main

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/CalypsoSys/aip_food_lookup/internal/foodcatalog"
)

func TestCollectLoggedSearchesReadsPlainAndGzipLogs(t *testing.T) {
	directory := t.TempDir()
	plain := "x - - [18/Aug/2026:01:29:51 +0000] \"GET /search?key=aga HTTP/1.1\" 200 49 \"-\" \"Dart\" 1ms\n"
	plain += "x - - [18/Aug/2026:01:29:52 +0000] \"GET /robots.txt HTTP/1.1\" 404 19 \"-\" \"-\" 0ms\n"
	if err := os.WriteFile(filepath.Join(directory, "access.log"), []byte(plain), 0644); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(filepath.Join(directory, "access.log.1.gz"))
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	_, _ = writer.Write([]byte("x - - [18/Aug/2026:01:29:53 +0000] \"GET /search?key=agave HTTP/1.1\" 200 49 \"-\" \"Dart\" 1ms\n"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	searches, total, err := collectLoggedSearches(directory)
	if err != nil || total != 2 || searches["aga"].Count != 1 || searches["agave"].Count != 1 {
		t.Fatalf("unexpected result: total=%d searches=%v err=%v", total, searches, err)
	}
}

func TestCatalogStoreUsesPrefixAndSoundMatching(t *testing.T) {
	directory := t.TempDir()
	if err := os.Mkdir(filepath.Join(directory, "allowed"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "allowed", "fruit.dat"), []byte("Agave\n"), 0644); err != nil {
		t.Fatal(err)
	}
	foods, err := foodcatalog.Load(directory)
	if err != nil {
		t.Fatal(err)
	}
	if !foodcatalog.Covered(foods, "aga") || !foodcatalog.Covered(foods, "agave") {
		t.Fatal("expected catalog prefix matches")
	}
}
