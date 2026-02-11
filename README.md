# okini

[[Japanese](./README-ja.md)]

A CLI tool for bookmarking local file paths

## Installation

**Unix/Linux/macOS:**
```bash
go build -o ~/go/bin/okini
```

**Windows (PowerShell):**
```powershell
go build -o "$env:USERPROFILE\go\bin\okini.exe"
```

**Windows (Command Prompt):**
```cmd
go build -o %USERPROFILE%\go\bin\okini.exe
```

Make sure `~/go/bin` (Unix) or `%USERPROFILE%\go\bin` (Windows) is in your PATH.

## Usage

### Adding Bookmarks

```bash
# Specify path only (name will be the final path element)
okini --add /path/to/file

# Specify custom name
okini --add /path/to/file myfile
```

If a bookmark with the same name already exists, both the existing and new bookmarks will be automatically annotated with their full paths to avoid conflicts. For example:
- Existing: `cc` → becomes `cc <= /aaa/bb/cc`
- New: `cc` → becomes `cc <= /dd/ff/cc`

### Removing Bookmarks

```bash
# Remove by path
okini --remove /path/to/file

# Remove by name
okini --remove myfile
```

The `--remove` command can remove bookmarks by either path or name:
- If the argument is a path, it removes by path
- Otherwise, it removes by name

This will remove all bookmarks matching the specified path or name.

After removal, if a bookmark no longer has name conflicts, its annotation will be automatically simplified. For example:
- Before removal: `aa <= C:/Users/username/Desktop/aa` and `aa <= C:/Users/username/Desktop/bb/aa`
- After removing one: `aa` (annotation removed since there's no conflict anymore)

### Listing Bookmark Names

```bash
okini --list
```

This is designed to be piped to tools like fzf.

### Getting Path by Name

```bash
okini --search myfile
```

### Command-line Examples

**Basic usage:**
```bash
# Add bookmarks
okini --add /path/to/file
okini --add /path/to/file myfile

# Remove by path or name
okini --remove /path/to/file       # Remove by path
okini --remove myfile              # Remove by name

# List and search
okini --list
okini --search myfile
```

**With fzf:**

**Unix/Linux/macOS:**
```bash
okini --list | fzf | xargs okini --search
```

**Windows (PowerShell):**
```powershell
okini --search (okini --list | fzf)
```

### Combined Usage with fzf

Select a bookmark and get its path:

**Unix/Linux/macOS:**
```bash
okini --list | fzf | xargs -I {} okini --search {}
```

**Windows (PowerShell):**
```powershell
okini --search (okini --list | fzf)
```

Navigate to selected path:

**Unix/Linux/macOS:**
```bash
cd $(okini --list | fzf | xargs -I {} okini --search {})
```

**Windows (PowerShell):**
```powershell
cd (okini --search (okini --list | fzf))
```

Registering as an alias is convenient:

**Bash/Zsh (add to .bashrc or .zshrc):**
```bash
alias cdoki='cd $(okini --list | fzf | xargs -I {} okini --search {})'
```

**PowerShell (add to $PROFILE):**
```powershell
function cdoki { cd (okini --search (okini --list | fzf)) }
```

## Data Storage Location

Bookmark data is stored in the following directories:

- Linux: `~/.config/okini/bookmarks.json`
- macOS: `~/Library/Application Support/okini/bookmarks.json`
- Windows: `%AppData%\okini\bookmarks.json`

## Data Format

```json
[
  {
    "name": "myfile",
    "path": "/absolute/path/to/file"
  },
  {
    "name": "documents",
    "path": "/home/user/Documents"
  }
]
```