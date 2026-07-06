package generator

import (
	"fmt"
	"html/template"
	"md-ssg/internal/parser"
	"os"
	"path/filepath"
)

type TemplateData struct {
	Title   string
	Date    string
	Content template.HTML
}

type Generator struct {
	tmpl *template.Template
}

func New(templatePath string) (*Generator, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf(" %w", err)
	}
	return &Generator{tmpl: tmpl}, nil
}

func (g *Generator) Render(post parser.Post, OutDir string) error {
	data := TemplateData{
		Title:   post.Title,
		Date:    post.Date,
		Content: template.HTML(post.Content),
	}

	// 出力ファイルパス
	outPath := filepath.Join(OutDir, post.Slug+".html")
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("出力ファイルの作成に失敗しました: %w", err)
	}

	defer file.Close()
	if err := g.tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("ファイルの変換に失敗しました: %w", err)
	}
	return nil
}
