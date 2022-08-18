package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/o8x/dsg/internal/downloader"
	"github.com/o8x/dsg/internal/generate"
)

var (
	filtersUrl *string
	usage      = `Usage:
  generate -url https://xxx

Use go run:
  go get -d github.com/o8x/dsg
  go run github.com/o8x/dsg/cmd/generator -url https://xxx
`
)

func init() {
	filtersUrl = flag.String("url", "", usage)
	flag.Parse()
}

func main() {
	if *filtersUrl == "" {
		fmt.Print(usage)
		return
	}

	reader, err := downloader.DownAsReader(*filtersUrl)
	if err != nil {
		panic(err)
	}

	_ = os.Mkdir("dsg", 0755)
	if err = os.WriteFile("./dsg/dsg.go", []byte(generate.Template(reader)), 0755); err != nil {
		panic(err)
	}
}
