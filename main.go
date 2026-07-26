package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	delete := flag.Bool("d", false, "delete orphan .RAF files (default: dry-run)")
	flag.Parse()

	folder := "."
	if flag.NArg() > 0 {
		folder = flag.Arg(0)
	}
	abs, err := filepath.Abs(folder)
	if err != nil {
		fail(err)
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		fail(fmt.Errorf("not a directory: %s", abs))
	}

	orphans, kept, err := findOrphans(abs)
	if err != nil {
		fail(err)
	}

	for _, path := range orphans {
		name := filepath.Base(path)
		if !*delete {
			fmt.Println(name)
			continue
		}
		if err := os.Remove(path); err != nil {
			fail(err)
		}
		fmt.Println("deleted", name)
	}

	action := "would delete"
	if *delete {
		action = "deleted"
	}
	fmt.Printf("%s %d orphan .RAF (%d kept) in %s\n", action, len(orphans), kept, abs)
}

func findOrphans(folder string) (orphans []string, kept int, err error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, 0, err
	}

	jpg := map[string]struct{}{}
	var rafs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		stem := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
		switch ext {
		case ".jpg", ".jpeg":
			jpg[stem] = struct{}{}
		case ".raf":
			rafs = append(rafs, filepath.Join(folder, name))
		}
	}

	for _, path := range rafs {
		stem := strings.ToLower(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		if _, ok := jpg[stem]; ok {
			kept++
		} else {
			orphans = append(orphans, path)
		}
	}
	return orphans, kept, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
