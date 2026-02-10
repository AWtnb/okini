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
okini --remove /path/to/file
```

This will remove all bookmarks with the matching path, even if they have different names.

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