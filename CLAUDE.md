# CLAUDE.md

このファイルは、リポジトリのコードを扱う際に Claude Code (claude.ai/code) へ提供するガイダンスです。

## コマンド

```bash
# ビルド
go build ./...

# 実行（posts/ の Markdown を public/ に HTML 変換）
go run . -src posts -out public

# テスト実行
go test ./...

# 特定パッケージのテストのみ実行
go test ./internal/parser/...
```

## アーキテクチャ

Markdown ファイルを HTML に変換する静的サイトジェネレーター (SSG) です。

**データフロー:**
1. `main.go` が CLI フラグ (`-src`, `-out`) を読み取り、ソースディレクトリを走査する
2. 各 `.md` ファイルを `internal/parser.ParseFile(slug, raw)` に渡す
   - `---` 区切りの YAML フロントマターを `gopkg.in/yaml.v3` で分割
   - Markdown 本文を `github.com/yuin/goldmark` で HTML にレンダリング
   - `Post{Title, Date, Content, Slug}` を返す
3. 生成された `Post` を `template/post.html` に渡してレンダリングする想定（Go の `html/template` 構文: `{{ .Title }}`、`{{ .Date }}`、`{{ .Content }}`）

**主要な型:**
- `Post` — パーサーからテンプレートへ渡す中心的なデータ構造体
- `frontMatter` — YAML アンマーシャル用の内部構造体（非公開）

**現在の状態:** `main.go` での `template/post.html` へのレンダリング処理はまだ実装されていない。現状はパース結果を標準出力に表示するのみ。
