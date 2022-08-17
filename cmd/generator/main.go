package main

import (
	_ "embed"
	"flag"
	"os"

	"github.com/o8x/dsg/internal/downloader"
	"github.com/o8x/dsg/internal/generate"
)

var (
	filtersUrl *string
)

func init() {
	filtersUrl = flag.String("url", "",
		`go run github.com/o8x/dsg/cmd/generator -url https-link`,
	)
	flag.Parse()
}

func main() {
	reader, err := downloader.DownAsReader(*filtersUrl)
	if err != nil {
		panic(err)
	}

	_ = os.Mkdir("dsg", 0755)
	if err = os.WriteFile("./dsg/dsg.go", []byte(generate.Template(reader)), 0755); err != nil {
		panic(err)
	}
}
