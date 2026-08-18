package foodcatalog

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CalypsoSys/godoublemetaphone/pkg/godoublemetaphone"
)

type Food struct {
	Allowed                 bool
	Name                    string
	PrimaryShortMetaphone   uint16
	AlternateShortMetaphone uint16
}

type Result struct {
	Allowed    []string
	NotAllowed []string
}

func Load(directory string) ([]Food, error) {
	var foods []Food
	err := filepath.Walk(directory, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || filepath.Ext(path) != ".dat" {
			return nil
		}
		folder := filepath.Base(filepath.Dir(path))
		if folder != "allowed" && folder != "not_allowed" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			name := strings.TrimSpace(scanner.Text())
			if name == "" {
				continue
			}
			metaphone := godoublemetaphone.NewShortDoubleMetaphone(name)
			foods = append(foods, Food{Allowed: folder == "allowed", Name: name, PrimaryShortMetaphone: metaphone.PrimaryShortKey(), AlternateShortMetaphone: metaphone.AlternateShortKey()})
		}
		return scanner.Err()
	})
	if err != nil {
		return nil, err
	}
	return foods, nil
}

func Match(foods []Food, query string, typeSearch string) Result {
	query = strings.ToLower(strings.TrimSpace(query))
	sdm := godoublemetaphone.NewShortDoubleMetaphone(query)
	textSearch, soundSearch := true, true
	if typeSearch == "searchbytext" {
		soundSearch = false
	}
	if typeSearch == "searchbysound" {
		textSearch = false
	}
	var allowed, notAllowed []string
	for _, food := range foods {
		if textSearch && strings.HasPrefix(strings.ToLower(food.Name), query) {
			if food.Allowed {
				allowed = append(allowed, food.Name)
			} else {
				notAllowed = append(notAllowed, food.Name)
			}
		}
		if soundSearch && fuzzySoundMatch(query, sdm, food) {
			if food.Allowed {
				allowed = append(allowed, food.Name)
			} else {
				notAllowed = append(notAllowed, food.Name)
			}
		}
	}
	return Result{Allowed: sortedUnique(allowed), NotAllowed: sortedUnique(notAllowed)}
}

func Covered(foods []Food, query string) bool {
	result := Match(foods, query, "searchbytextandsound")
	return len(result.Allowed) > 0 || len(result.NotAllowed) > 0
}

// SpellingDistanceAllowed exposes the catalog's spelling threshold for focused tests.
func SpellingDistanceAllowed(query, candidate string) bool {
	return spellingDistanceAllowed(query, candidate)
}

func fuzzySoundMatch(query string, queryMetaphone godoublemetaphone.ShortDoubleMetaphone, food Food) bool {
	if spellingDistanceAllowed(query, food.Name) {
		return true
	}
	for _, token := range searchableTokens(food.Name) {
		if spellingDistanceAllowed(query, token) {
			return true
		}
	}
	if !metaphoneKeysMatch(queryMetaphone, food) {
		return false
	}
	limit := spellingDistanceLimit(query)
	if len(query) > 4 {
		limit++
	}
	if levenshteinDistance(query, strings.ToLower(food.Name)) <= limit {
		return true
	}
	for _, token := range searchableTokens(food.Name) {
		if levenshteinDistance(query, token) <= limit {
			return true
		}
	}
	return false
}

func metaphoneKeysMatch(query godoublemetaphone.ShortDoubleMetaphone, food Food) bool {
	queryKeys := validMetaphoneKeys(query.PrimaryShortKey(), query.AlternateShortKey())
	foodKeys := validMetaphoneKeys(food.PrimaryShortMetaphone, food.AlternateShortMetaphone)
	for key := range queryKeys {
		if foodKeys[key] {
			return true
		}
	}
	return false
}

func validMetaphoneKeys(keys ...uint16) map[uint16]bool {
	valid := make(map[uint16]bool)
	for _, key := range keys {
		if key != godoublemetaphone.METAPHONE_INVALID_KEY {
			valid[key] = true
		}
	}
	return valid
}

func spellingDistanceAllowed(query, candidate string) bool {
	query, candidate = strings.ToLower(strings.TrimSpace(query)), strings.ToLower(strings.TrimSpace(candidate))
	if query == "" || candidate == "" {
		return false
	}
	if query == candidate {
		return true
	}
	if query[0] != candidate[0] {
		return false
	}
	return levenshteinDistance(query, candidate) <= spellingDistanceLimit(query)
}

func searchableTokens(candidate string) []string {
	fields := strings.FieldsFunc(strings.ToLower(candidate), func(r rune) bool { return r < 'a' || r > 'z' })
	var result []string
	for _, field := range fields {
		if len(field) >= 3 {
			result = append(result, field)
		}
	}
	return result
}

func spellingDistanceLimit(query string) int {
	if len(query) <= 4 {
		return 1
	}
	return 3
}

func levenshteinDistance(a, b string) int {
	previous, current := make([]int, len(b)+1), make([]int, len(b)+1)
	for index := range previous {
		previous[index] = index
	}
	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			current[j] = min(current[j-1]+1, previous[j]+1, previous[j-1]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func min(values ...int) int {
	result := values[0]
	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}
	return result
}

func sortedUnique(items []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(items))
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	sort.Strings(result)
	return result
}
