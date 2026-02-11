package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Bookmark struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Bookmarks []Bookmark

const appName = "okini"

func getDataPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dataDir := filepath.Join(userConfigDir, appName)
	// Create directory with owner rwx, group rx, others rx permissions
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "bookmarks.json"), nil
}

func loadBookmarks() (Bookmarks, error) {
	dataPath, err := getDataPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Bookmarks{}, nil
		}
		return nil, err
	}

	var bookmarks Bookmarks
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func saveBookmarks(bookmarks Bookmarks) error {
	dataPath, err := getDataPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return err
	}

	// Write file with owner rw, group r, others r permissions
	return os.WriteFile(dataPath, data, 0644)
}

func annotatedName(path string) string {
	return fmt.Sprintf("%s <= %s", filepath.Base(path), filepath.ToSlash(path))
}

func getBaseName(name string) string {
	// If already annotated, extract the base name part
	if before, _, ok := strings.Cut(name, " <= "); ok {
		return before
	}
	return name
}

func simplifyAnnotations(bookmarks Bookmarks) Bookmarks {
	// Create a map to count base names
	baseNameCount := make(map[string]int)
	for _, bm := range bookmarks {
		baseName := getBaseName(bm.Name)
		baseNameCount[baseName]++
	}

	// Remove annotations for bookmarks that no longer have conflicts
	result := make(Bookmarks, len(bookmarks))
	for i, bm := range bookmarks {
		baseName := getBaseName(bm.Name)
		// If this is the only bookmark with this base name and it's annotated, simplify it
		if baseNameCount[baseName] == 1 && strings.Contains(bm.Name, " <= ") {
			result[i] = Bookmark{Name: baseName, Path: bm.Path}
		} else {
			result[i] = bm
		}
	}

	return result
}

func addBookmark(path, name string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check if path exists
	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	// Use base name if name is not specified
	if name == "" {
		name = filepath.Base(absPath)
	}

	bookmarks, err := loadBookmarks()
	if err != nil {
		return err
	}

	// If same bookmark already exists, do nothing
	for _, bm := range bookmarks {
		if bm.Name == name && bm.Path == absPath {
			return fmt.Errorf("entry duplicates: name=%s, path=%s", name, absPath)
		}
	}

	// Check if there's any bookmark with the same base name
	hasConflict := false
	for i, bm := range bookmarks {
		if getBaseName(bm.Name) == name {
			hasConflict = true
			// Annotate existing bookmark if not already annotated
			if !strings.Contains(bm.Name, " <= ") {
				bookmarks[i].Name = annotatedName(bm.Path)
			}
		}
	}

	// If there's a conflict, annotate the new bookmark too
	if hasConflict {
		name = annotatedName(absPath)
	}

	bookmarks = append(bookmarks, Bookmark{
		Name: name,
		Path: absPath,
	})

	return saveBookmarks(bookmarks)
}

func listBookmarks() error {
	bookmarks, err := loadBookmarks()
	if err != nil {
		return err
	}

	for _, bm := range bookmarks {
		fmt.Println(bm.Name)
	}
	return nil
}

func searchBookmark(name string) error {
	bookmarks, err := loadBookmarks()
	if err != nil {
		return err
	}

	for _, bm := range bookmarks {
		if bm.Name == name {
			fmt.Println(bm.Path)
			return nil
		}
	}

	return fmt.Errorf("bookmark not found: %s", name)
}

func removeBookmark(path string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	bookmarks, err := loadBookmarks()
	if err != nil {
		return err
	}

	// Filter out all bookmarks with matching path
	filtered := make(Bookmarks, 0)
	removedCount := 0
	for _, bm := range bookmarks {
		if bm.Path != absPath {
			filtered = append(filtered, bm)
		} else {
			removedCount++
		}
	}

	if removedCount == 0 {
		return fmt.Errorf("no bookmark found with path: %s", absPath)
	}

	// Simplify annotations if conflicts are resolved
	filtered = simplifyAnnotations(filtered)

	if err := saveBookmarks(filtered); err != nil {
		return err
	}

	fmt.Printf("Removed %d bookmark(s)\n", removedCount)
	return nil
}

func removeBookmarkByName(name string) error {
	bookmarks, err := loadBookmarks()
	if err != nil {
		return err
	}

	// Filter out all bookmarks with matching name
	filtered := make(Bookmarks, 0)
	removedCount := 0
	for _, bm := range bookmarks {
		if getBaseName(bm.Name) != name {
			filtered = append(filtered, bm)
		} else {
			removedCount++
		}
	}

	if removedCount == 0 {
		return fmt.Errorf("no bookmark found with name: %s", name)
	}

	// Simplify annotations if conflicts are resolved
	filtered = simplifyAnnotations(filtered)

	if err := saveBookmarks(filtered); err != nil {
		return err
	}

	fmt.Printf("Removed %d bookmark(s)\n", removedCount)
	return nil
}

func run() int {
	addCmd := flag.String("add", "", "Add a bookmark for the file path")
	listCmd := flag.Bool("list", false, "List all bookmark names")
	searchCmd := flag.String("search", "", "Search path by name")
	removeCmd := flag.String("remove", "", "Remove bookmark(s) by path or name")

	flag.Parse()

	// Add mode
	if *addCmd != "" {
		args := flag.Args()
		name := ""
		if 0 < len(args) {
			name = args[0]
		}

		if err := addBookmark(*addCmd, name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		fmt.Printf("Bookmark added: %s\n", *addCmd)
		return 0
	}

	// Remove mode
	if *removeCmd != "" {
		// If the argument is a path, remove by path
		if _, err := os.Stat(*removeCmd); err == nil {
			if err := removeBookmark(*removeCmd); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return 1
			}
			return 0
		}
		// Otherwise, remove by name
		if err := removeBookmarkByName(*removeCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		return 0
	}

	// List mode
	if *listCmd {
		if err := listBookmarks(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		return 0
	}

	// Search mode
	if *searchCmd != "" {
		if err := searchBookmark(*searchCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		return 0
	}

	// Show help
	fmt.Println("okini - File path bookmark tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  okini --add <file_path> [name]     Add a bookmark")
	fmt.Println("  okini --remove <path_or_name>      Remove bookmark(s) by path or name")
	fmt.Println("  okini --list                       List all bookmark names")
	fmt.Println("  okini --search <n>              Search path by name")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  okini --add /path/to/file")
	fmt.Println("  okini --add /path/to/file myfile")
	fmt.Println("  okini --remove /path/to/file       # Remove by path")
	fmt.Println("  okini --remove myfile              # Remove by name")
	fmt.Println("  okini --list | fzf | xargs okini --search")
	return 0
}

func main() {
	os.Exit(run())
}
