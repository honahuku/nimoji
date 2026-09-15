# nimoji

社員名簿のCSVから、日本語入力辞書ファイルを生成するCLIツールです。

名前・社員番号・メールアドレスを読みから変換できます。

## インストール

### Homebrew (macOS/Linux)

```bash
brew install hirosassa/tap/nimoji
```

### GitHub Releases からバイナリをダウンロード

[Releases ページ](https://github.com/hirosassa/nimoji/releases)からお使いのOS・アーキテクチャに合ったバイナリをダウンロードできます。

```bash
# 例: macOS (Apple Silicon)
curl -sL https://github.com/hirosassa/nimoji/releases/latest/download/nimoji_Darwin_arm64.tar.gz | tar xz
sudo mv nimoji /usr/local/bin/
```

### Go

```bash
go install github.com/hirosassa/nimoji@latest
```

## 使い方

```bash
nimoji -format <google|mac|msime> < employees.csv
```

### 入力CSV形式

```
社員番号,姓,名,メールアドレス,姓よみ,名よみ[,備考]
```

備考列はオプションです。所属部署などを記載できます。Google日本語入力への登録の際に、備考の項目を変換候補の補足事項として表示できます。例えば同姓同名の変換候補の区別などに便利です。

#### 例

```csv
001,田中,太郎,tanaka@example.com,たなか,たろう,営業部
002,佐藤,花子,sato@example.com,さとう,はなこ
```

### 登録されるエントリ

1人あたり以下の3エントリが登録されます。

| 読み | 変換先 | 例 |
|---|---|---|
| 姓よみ | 姓 + 名 | `たなか` → `田中太郎` |
| `ばんごう` + 姓よみ | 社員番号 | `ばんごうたなか` → `001` |
| `めーる` + 姓よみ | メールアドレス | `めーるたなか` → `tanaka@example.com` |

### 出力形式

#### Google日本語入力 (`-format google`)

TSV形式で出力されます。4列目にコメント（フルネーム / 備考）が付きます。

```bash
nimoji -format google < employees.csv > dictionary.txt
```

```
たなか	田中太郎	固有名詞	田中太郎 / 営業部
ばんごうたなか	001	固有名詞	田中太郎 / 営業部
めーるたなか	tanaka@example.com	固有名詞	田中太郎 / 営業部
```

**初回インポート手順:**

1. Google日本語入力の「辞書ツール」を開く
2. 「管理」→「新規辞書にインポート」を選択
3. 生成した `dictionary.txt` を指定し、辞書名を `nimoji` などわかりやすい名前にする

**更新時の手順（重複登録を防ぐため）:**

1. 「辞書ツール」で nimoji 用の辞書を選択
2. 全エントリを選択（Ctrl+A / Cmd+A）して削除
3. 「管理」→「選択した辞書にインポート」で新しい `dictionary.txt` を読み込む

nimoji 専用の辞書を分けておくことで、ほかの辞書に影響を与えずに安全に差し替えできます。

#### macユーザー辞書 (`-format mac`)

plist XML形式で出力されます。

```bash
nimoji -format mac < employees.csv > dictionary.plist
```

「システム設定」→「キーボード」→「ユーザー辞書」にドラッグ&ドロップで読み込めます。

**注意:** macユーザー辞書には辞書を分ける機能がないため、再インポート時に重複登録される可能性があります。更新する場合は、既存の nimoji 由来のエントリを手動で削除してから再度インポートしてください。

#### Microsoft IME ユーザー辞書 (`-format msime`)

Microsoft IME のユーザー辞書ツールが「テキストファイルからの登録」で取り込める
WORDLIST 形式（UTF-16LE, BOM付き, CRLF改行）で出力されます。

```bash
nimoji -format msime < employees.csv > dictionary.txt
```

「ユーザー辞書ツール」→「ツール」→「テキストファイルからの登録」で読み込めます。

**注意:** Microsoft IME 固有の制約に対応するため、`-format google` と異なりコメントの切り詰め・空エントリのスキップ・読みのひらがな変換を行います。

## ライセンス

MIT
