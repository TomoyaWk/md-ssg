package generator

import (
	"fmt"
	"html/template"
	"io"
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
	tmpl   *template.Template
	cssSrc string
}

func New(templatePath string) (*Generator, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf(" %w", err)
	}
	cssSrc := filepath.Join(filepath.Dir(templatePath), "style.css")
	return &Generator{tmpl: tmpl, cssSrc: cssSrc}, nil
}

// CopyStaticAssets はテンプレートと同じディレクトリの style.css を出力先へコピーする
func (g *Generator) CopyStaticAssets(outDir string) error {
	src, err := os.Open(g.cssSrc)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("style.cssの読み込みに失敗しました: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(outDir, "style.css"))
	if err != nil {
		return fmt.Errorf("style.cssの出力に失敗しました: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("style.cssのコピーに失敗しました: %w", err)
	}
	return nil
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
