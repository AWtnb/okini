# okini

ローカルのファイルパスをブックマークするCLIツール

## インストール


```
go install github.com/AWtnb/okini@latest
```

もしくはこのリポジトリをクローンしてから下記を実行。

```
go install .
```


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
# パスで削除
okini --remove /path/to/file

# 名前で削除
okini --remove myfile
```

`--remove`コマンドはパスまたは名前でブックマークを削除できます：
- 引数がパスの場合、パスで削除
- それ以外の場合、名前で削除

指定したパスまたは名前に一致するブックマークをすべて削除します。

削除後、名前の衝突がなくなったブックマークの注釈は自動的に簡略化されます。例：
- 削除前: `aa <= C:/Users/username/Desktop/aa` と `aa <= C:/Users/username/Desktop/bb/aa`
- 片方削除後: `aa` （衝突がなくなったので注釈が除去される）

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

**基本的な使い方:**
```bash
# ブックマークの追加
okini --add /path/to/file
okini --add /path/to/file myfile

# パスまたは名前で削除
okini --remove /path/to/file       # パスで削除
okini --remove myfile              # 名前で削除

# 一覧表示と検索
okini --list
okini --search myfile
```

**fzfと組み合わせて:**

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