package main

import (
	"flag"
	"fmt"
	"md-ssg/internal/generator"
	"md-ssg/internal/parser"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SrcDir       string
	OutDir       string
	TemplatePath string
}

// 引数パース
func parseArgs() Config {
	src := flag.String("src", "posts", "md格納先フォルダ指定")
	out := flag.String("out", "public", "HTML出力先フォルダ指定")
	tmpl := flag.String("template", "templates/post.html", "パーステンプレートファイル場所指定")
	flag.Parse()

	return Config{
		SrcDir:       *src,
		OutDir:       *out,
		TemplatePath: *tmpl,
	}
}

func run(cfg Config) error {
	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return fmt.Errorf("HTML出力先フォルダの作成に失敗: %w", err)
	}

	mds, err := os.ReadDir(cfg.SrcDir)
	if err != nil {
		return fmt.Errorf("md格納先フォルダの読み込みに失敗: %w", err)
	}

	gen, err := generator.New(cfg.TemplatePath)
	if err != nil {
		return fmt.Errorf("テンプレートファイルの読み込みに失敗: %w", err)
	}

	if err := gen.CopyStaticAssets(cfg.OutDir); err != nil {
		return fmt.Errorf("静的アセットのコピーに失敗: %w", err)
	}

	for _, md := range mds {
		if md.IsDir() || filepath.Ext(md.Name()) != ".md" {
			continue
		}
		srcPath := filepath.Join(cfg.SrcDir, md.Name())
		raw, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("ファイルの読み込みに失敗: %s: %w", md.Name(), err)
		}
		slug := strings.TrimSuffix(md.Name(), ".md")
		post, err := parser.ParseFile(slug, string(raw))
		if err != nil {
			return fmt.Errorf("ファイルの解析に失敗: %s: %w", md.Name(), err)
		}

		if err := gen.Render(post, cfg.OutDir); err != nil {
			return fmt.Errorf("render %s: %w", md.Name(), err)
		}
		fmt.Printf("generated: %s.html\n", post.Slug)
	}
	return nil
}

func main() {
	cfg := parseArgs()

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print("test")
	os.Exit(0)
}
