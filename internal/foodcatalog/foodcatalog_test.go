package foodcatalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CalypsoSys/godoublemetaphone/pkg/godoublemetaphone"
)

func TestParseEntry(t *testing.T) {
	name, aliases, ok := ParseEntry("Chobani Yogurt - All\tchobani yogurt\tchobani vanilla yogurt")
	if !ok {
		t.Fatal("expected catalog entry to parse")
	}
	if name != "Chobani Yogurt - All" {
		t.Fatalf("name = %q", name)
	}
	if len(aliases) != 2 || aliases[1] != "chobani vanilla yogurt" {
		t.Fatalf("aliases = %#v", aliases)
	}
}

func TestLoadEntriesYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dairy.yaml")
	contents := []byte("- name: Chobani Yogurt - All\n  aliases:\n    - chobani vanilla yogurt\n")
	if err := os.WriteFile(path, contents, 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := LoadEntries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || len(entries[0].Aliases) != 1 {
		t.Fatalf("entries = %#v", entries)
	}
}

func TestMatchAliasReturnsCanonicalName(t *testing.T) {
	name := "Chobani Yogurt - All"
	metaphone := godoublemetaphone.NewShortDoubleMetaphone(name)
	foods := []Food{{
		Allowed:                 false,
		Name:                    name,
		Aliases:                 []string{"chobani vanilla yogurt"},
		PrimaryShortMetaphone:   metaphone.PrimaryShortKey(),
		AlternateShortMetaphone: metaphone.AlternateShortKey(),
	}}

	result := Match(foods, "chobani vanilla yogurt", "searchbytext")
	if len(result.NotAllowed) != 1 || result.NotAllowed[0] != name {
		t.Fatalf("not allowed result = %#v", result.NotAllowed)
	}
}

func TestMatchMultiWordSoundRequiresEachWord(t *testing.T) {
	foods := []Food{
		foodForTest("Coconut Milk"),
		foodForTest("Coconut Oil"),
	}

	result := Match(foods, "coconut milk", "searchbysound")
	if len(result.Allowed) != 1 || result.Allowed[0] != "Coconut Milk" {
		t.Fatalf("allowed result = %#v", result.Allowed)
	}
}

func foodForTest(name string) Food {
	metaphone := godoublemetaphone.NewShortDoubleMetaphone(name)
	return Food{
		Allowed:                 true,
		Name:                    name,
		PrimaryShortMetaphone:   metaphone.PrimaryShortKey(),
		AlternateShortMetaphone: metaphone.AlternateShortKey(),
	}
}
