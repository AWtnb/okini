package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	// Find existing bookmark and overwrite, or append new one
	found := false
	for i, bm := range bookmarks {
		if bm.Name == name {
			bookmarks[i].Path = absPath
			found = true
			break
		}
	}

	if !found {
		bookmarks = append(bookmarks, Bookmark{
			Name: name,
			Path: absPath,
		})
	}

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

func run() int {
	addCmd := flag.String("add", "", "Add a bookmark for the file path")
	listCmd := flag.Bool("list", false, "List all bookmark names")
	searchCmd := flag.String("search", "", "Search path by name")

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
	fmt.Println("  okini --add <file_path> [name]  Add a bookmark")
	fmt.Println("  okini --list                    List all bookmark names")
	fmt.Println("  okini --search <name>           Search path by name")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  okini --add /path/to/file")
	fmt.Println("  okini --add /path/to/file myfile")
	fmt.Println("  okini --list | fzf | xargs okini --search")
	return 0
}

func main() {
	os.Exit(run())
}
