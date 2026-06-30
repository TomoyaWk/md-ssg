package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Title   string
	Date    string
	Content string // レンダリング済みHTML
	Slug    string // ファイル名から拡張子を除いたもの
}

type frontMatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

// ParseFile はMarkdownファイルの内容を受け取りPostを返す
func ParseFile(slug, src string) (Post, error) {
	fm, body, err := splitFrontMatter(src)
	if err != nil {
		return Post{}, fmt.Errorf("parse front matter: %w", err)
	}

	html, err := renderMarkdown(body)
	if err != nil {
		return Post{}, fmt.Errorf("render markdown: %w", err)
	}

	return Post{
		Title:   fm.Title,
		Date:    fm.Date,
		Content: html,
		Slug:    slug,
	}, nil
}

// splitFrontMatter は --- で囲まれたYAMLブロックと本文を分離する
func splitFrontMatter(src string) (frontMatter, string, error) {
	var fm frontMatter

	// Front Matterがない場合はそのまま返す
	if !strings.HasPrefix(src, "---") {
		return fm, src, nil
	}

	// 2つ目の --- を探す
	rest := src[3:] // 最初の --- をスキップ
	idx := strings.Index(rest, "---")
	if idx == -1 {
		return fm, src, fmt.Errorf("unclosed front matter block")
	}

	yamlBlock := strings.TrimSpace(rest[:idx])
	body := strings.TrimSpace(rest[idx+3:])

	if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err != nil {
		return fm, body, fmt.Errorf("yaml unmarshal: %w", err)
	}

	return fm, body, nil
}

func renderMarkdown(src string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
