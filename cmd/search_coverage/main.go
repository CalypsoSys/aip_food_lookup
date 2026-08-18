package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/CalypsoSys/aip_food_lookup/internal/foodcatalog"
)

var accessSearchLogPattern = regexp.MustCompile(`\[([^]]+)\] "GET ([^ ]+) HTTP/[^ ]+" ([0-9]+) `)

type loggedSearch struct {
	Count     int
	Statuses  map[int]int
	FirstSeen string
	LastSeen  string
}

func main() {
	if err := runCLI(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCLI(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: search_coverage <extract|check> [options]")
	}
	switch args[0] {
	case "extract":
		return extractSearches(args[1:])
	case "check":
		return checkSearches(args[1:])
	default:
		return fmt.Errorf("unknown mode %q; use extract or check", args[0])
	}
}

func extractSearches(args []string) error {
	flags := flag.NewFlagSet("extract", flag.ContinueOnError)
	logs := flags.String("logs", "/srv/logs/aip-food-lookup/api", "directory containing access.log files")
	output := flags.String("output", "searches.tsv", "output TSV path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	searches, total, err := collectLoggedSearches(*logs)
	if err != nil {
		return err
	}
	file, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("create output %s: %w", *output, err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, "# AIP Food Lookup search coverage export")
	fmt.Fprintln(writer, "# key\tcount\tfirst_seen\tlast_seen\tstatuses")
	for _, key := range sortedKeys(searches) {
		search := searches[key]
		fmt.Fprintf(writer, "%s\t%d\t%s\t%s\t%s\n", key, search.Count, search.FirstSeen, search.LastSeen, formatStatuses(search.Statuses))
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("write output %s: %w", *output, err)
	}
	fmt.Printf("Wrote %d unique search keys (%d requests) to %s\n", len(searches), total, *output)
	return nil
}

func checkSearches(args []string) error {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	input := flags.String("input", "searches.tsv", "TSV exported by extract mode")
	catalog := flags.String("catalog", "../../data", "local repository data directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	searches, err := readSearchExport(*input)
	if err != nil {
		return err
	}
	foods, err := foodcatalog.Load(*catalog)
	if err != nil {
		return fmt.Errorf("load local catalog: %w", err)
	}
	covered := 0
	for _, key := range sortedKeys(searches) {
		if foodcatalog.Covered(foods, key) {
			covered++
		}
	}
	fmt.Printf("Local catalog: %s\nSearch keys:   %d\nCovered:       %d\nUncovered:     %d\n\n", *catalog, len(searches), covered, len(searches)-covered)
	fmt.Println("UNCOVERED SEARCHES")
	fmt.Println("count | key | statuses | first seen | last seen")
	fmt.Println("------|-----|----------|------------|----------")
	for _, key := range sortedKeys(searches) {
		if !foodcatalog.Covered(foods, key) {
			search := searches[key]
			fmt.Printf("%5d | %s | %s | %s | %s\n", search.Count, key, formatStatuses(search.Statuses), search.FirstSeen, search.LastSeen)
		}
	}
	fmt.Println()
	fmt.Println("COVERED SEARCHES")
	fmt.Println("count | key")
	fmt.Println("------|-----")
	for _, key := range sortedKeys(searches) {
		if foodcatalog.Covered(foods, key) {
			fmt.Printf("%5d | %s\n", searches[key].Count, key)
		}
	}
	return nil
}

func collectLoggedSearches(directory string) (map[string]loggedSearch, int, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "access.log*"))
	if err != nil {
		return nil, 0, err
	}
	sort.Strings(paths)
	searches := make(map[string]loggedSearch)
	total := 0
	for _, path := range paths {
		if err := scanAccessLog(path, searches, &total); err != nil {
			return nil, total, err
		}
	}
	return searches, total, nil
}

func scanAccessLog(path string, searches map[string]loggedSearch, total *int) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open access log %s: %w", path, err)
	}
	defer file.Close()
	var reader io.Reader = file
	if strings.HasSuffix(path, ".gz") {
		compressed, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("open gzip access log %s: %w", path, err)
		}
		defer compressed.Close()
		reader = compressed
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		match := accessSearchLogPattern.FindStringSubmatch(scanner.Text())
		if len(match) == 0 {
			continue
		}
		target, err := url.ParseRequestURI(match[2])
		if err != nil || target.Path != "/search" {
			continue
		}
		key := strings.TrimSpace(target.Query().Get("key"))
		status, err := strconv.Atoi(match[3])
		if key == "" || err != nil {
			continue
		}
		search := searches[key]
		if search.Statuses == nil {
			search.Statuses = make(map[int]int)
		}
		search.Count++
		search.Statuses[status]++
		if search.FirstSeen == "" || match[1] < search.FirstSeen {
			search.FirstSeen = match[1]
		}
		if match[1] > search.LastSeen {
			search.LastSeen = match[1]
		}
		searches[key] = search
		*total++
	}
	return scanner.Err()
}

func readSearchExport(path string) (map[string]loggedSearch, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open search export %s: %w", path, err)
	}
	defer file.Close()
	searches := make(map[string]loggedSearch)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			return nil, fmt.Errorf("invalid search export line: %q", line)
		}
		count, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("invalid count in search export: %q", line)
		}
		statuses := make(map[int]int)
		for _, entry := range strings.Split(fields[4], ",") {
			parts := strings.SplitN(entry, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid statuses in search export: %q", line)
			}
			status, statusErr := strconv.Atoi(parts[0])
			statusCount, countErr := strconv.Atoi(parts[1])
			if statusErr != nil || countErr != nil {
				return nil, fmt.Errorf("invalid statuses in search export: %q", line)
			}
			statuses[status] = statusCount
		}
		searches[fields[0]] = loggedSearch{Count: count, FirstSeen: fields[2], LastSeen: fields[3], Statuses: statuses}
	}
	return searches, scanner.Err()
}

func sortedKeys(searches map[string]loggedSearch) []string {
	keys := make([]string, 0, len(searches))
	for key := range searches {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func formatStatuses(statuses map[int]int) string {
	keys := make([]int, 0, len(statuses))
	for key := range statuses {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%d:%d", key, statuses[key]))
	}
	return strings.Join(parts, ",")
}
