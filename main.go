package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	SrcDir string
	OutDir string
}

// 引数パース
func parseArgs() Config {
	src := flag.String("src", "posts", "md格納先フォルダ指定")
	out := flag.String("out", "public", "HTML出力先フォルダ指定")
	flag.Parse()

	return Config{
		SrcDir: *src,
		OutDir: *out,
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
	for _, md := range mds {
		if md.IsDir() {
			continue
		}
		fmt.Printf("Found: %s\n", md.Name())
	}
	return nil
}

func main() {
	cfg := parseArgs()

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
