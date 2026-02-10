# okini

ローカルのファイルパスをブックマークするCLIツール

## インストール

**Unix/Linux/macOS:**
```bash
go build -o ~/go/bin/okini
```

**Windows (PowerShell):**
```powershell
go build -o "$env:USERPROFILE\go\bin\okini.exe"
```

**Windows (コマンドプロンプト):**
```cmd
go build -o %USERPROFILE%\go\bin\okini.exe
```

`~/go/bin` (Unix) または `%USERPROFILE%\go\bin` (Windows) がPATHに含まれていることを確認してください。

## 使い方

### ブックマークの登録

```bash
# パスのみ指定（名前は最終要素になる）
okini --add /path/to/file

# 名前を指定
okini --add /path/to/file myfile
```

同じ名前のブックマークが既に存在する場合、衝突を避けるために既存のブックマークと新しいブックマークの両方に自動的にフルパスの注釈が付けられます。例：
- 既存: `cc` → `cc <= /aaa/bb/cc` になる
- 新規: `cc` → `cc <= /dd/ff/cc` になる

### ブックマークの削除

```bash
okini --remove /path/to/file
```

指定したパスに一致するブックマークをすべて削除します（名前が異なっていても削除されます）。

### ブックマーク名の一覧取得

```bash
okini --list
```

fzfと組み合わせて使用する例：

```bash
okini --list | fzf
```

### 名前からパスの取得

```bash
okini --search myfile
```

### コマンドライン例

**Unix/Linux/macOS:**
```bash
okini --list | fzf | xargs okini --search
```

**Windows (PowerShell):**
```powershell
okini --search (okini --list | fzf)
```

### fzfと組み合わせた使用例

ブックマークから選択してパスを取得：

**Unix/Linux/macOS:**
```bash
okini --list | fzf | xargs -I {} okini --search {}
```

**Windows (PowerShell):**
```powershell
okini --search (okini --list | fzf)
```

選択したパスに移動：

**Unix/Linux/macOS:**
```bash
cd $(okini --list | fzf | xargs -I {} okini --search {})
```

**Windows (PowerShell):**
```powershell
cd (okini --search (okini --list | fzf))
```

エイリアスとして登録すると便利です：

**Bash/Zsh (.bashrc や .zshrc に追加):**
```bash
alias cdoki='cd $(okini --list | fzf | xargs -I {} okini --search {})'
```

**PowerShell ($PROFILE に追加):**
```powershell
function cdoki { cd (okini --search (okini --list | fzf)) }
```

## データの保存場所

ブックマークデータは以下のディレクトリに保存されます：

- Linux: `~/.config/okini/bookmarks.json`
- macOS: `~/Library/Application Support/okini/bookmarks.json`
- Windows: `%AppData%\okini\bookmarks.json`

## データ形式

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