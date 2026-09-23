// renamer: rename every file inside a directory to a single base name.
//
// Because two files can't share a name, files are numbered:
//   photo.jpg, notes.txt, song.mp3  ->  vacation_1.jpg, vacation_2.txt, vacation_3.mp3
// File extensions are kept. Sub-folders are left alone.
//
// Usage:
//   renamer <directory> <new-name>
//   renamer               (asks you for both)
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var dir, name string
	if len(os.Args) >= 3 {
		dir, name = os.Args[1], os.Args[2]
	} else {
		dir = ask(in, "Directory to rename files in (blank = current folder): ")
		if dir == "" {
			dir = "."
		}
		name = ask(in, "New name for the files: ")
	}

	if err := run(in, dir, name); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		pause(in)
		os.Exit(1)
	}
	pause(in)
}

func run(in *bufio.Reader, dir, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("new name can't be empty")
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return fmt.Errorf(`new name can't contain any of / \ : * ? " < > |`)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Don't rename this program if it lives in the target folder.
	self, _ := os.Executable()
	self, _ = filepath.EvalSymlinks(self)

	var files []string
	for _, e := range entries {
		if !e.Type().IsRegular() {
			continue // skip folders, symlinks, etc.
		}
		full, _ := filepath.Abs(filepath.Join(dir, e.Name()))
		if real, err := filepath.EvalSymlinks(full); err == nil && real == self {
			continue
		}
		files = append(files, e.Name())
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found in %q", dir)
	}
	sort.Strings(files)

	// Build the new names, padding numbers so they sort nicely (01, 02, ... 10).
	width := len(fmt.Sprint(len(files)))
	newNames := make([]string, len(files))
	for i, f := range files {
		newNames[i] = fmt.Sprintf("%s_%0*d%s", name, width, i+1, filepath.Ext(f))
	}

	fmt.Printf("\n%d file(s) in %s will be renamed:\n", len(files), dir)
	for i, f := range files {
		fmt.Printf("  %s  ->  %s\n", f, newNames[i])
	}
	if a := strings.ToLower(ask(in, "\nGo ahead? [y/N]: ")); a != "y" && a != "yes" {
		fmt.Println("Cancelled, nothing was changed.")
		return nil
	}

	// Two passes via temporary names, so a new name never clobbers
	// a file that hasn't been renamed yet.
	temps := make([]string, len(files))
	for i, f := range files {
		temps[i] = fmt.Sprintf(".renamer_tmp_%d_%d", os.Getpid(), i)
		if err := os.Rename(filepath.Join(dir, f), filepath.Join(dir, temps[i])); err != nil {
			rollback(dir, files[:i], temps[:i])
			return fmt.Errorf("renaming %s: %w", f, err)
		}
	}
	for i, t := range temps {
		dst := filepath.Join(dir, newNames[i])
		if _, err := os.Lstat(dst); err == nil {
			rollback(dir, files, temps)
			return fmt.Errorf("%s already exists, nothing was changed", newNames[i])
		}
		if err := os.Rename(filepath.Join(dir, t), dst); err != nil {
			return fmt.Errorf("renaming to %s: %w", newNames[i], err)
		}
	}

	fmt.Printf("Done! Renamed %d file(s).\n", len(files))
	return nil
}

func rollback(dir string, originals, temps []string) {
	for i, t := range temps {
		os.Rename(filepath.Join(dir, t), filepath.Join(dir, originals[i]))
	}
}

func ask(in *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := in.ReadString('\n')
	return strings.TrimSpace(line)
}

// pause keeps the window open when the program was double-clicked.
func pause(in *bufio.Reader) {
	if len(os.Args) < 3 {
		ask(in, "Press Enter to exit...")
	}
}
