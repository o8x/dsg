package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/o8x/dsg"
	"github.com/o8x/dsg/internal/downloader"
	"github.com/o8x/dsg/internal/generate"
)

func main() {
	var cmd = &cobra.Command{
		Use: "dsg",
	}

	cmd.AddCommand(makeGenerateCommand())
	cmd.AddCommand(makeTestCommand())
	_ = cmd.Execute()
}

func makeTestCommand() *cobra.Command {
	from := ""
	gen := &cobra.Command{
		Use:   "match --url https://url example.com",
		Short: "pattern test",
		Args:  checkURL,
		Run: func(cmd *cobra.Command, domains []string) {
			if err := dsg.Load(from); err != nil {
				log.Fatalf("load url %s: %v", from, err)
			}

			nd := dsg.Get()
			for _, d := range domains {
				if nd.Exist(d) {
					log.Printf("index hit: %s\n", d)
					continue
				}

				if p, ok := nd.Match(d); ok {
					log.Printf("match to rule: '%s'\n", p.Origin)
					continue
				}

				log.Printf("unable to match: %s\n", d)
			}
		},
		Example: "\tmatch --url https://url google.com",
	}

	addURLVar(gen, &from)
	return gen
}

func makeGenerateCommand() *cobra.Command {
	from := ""
	dest := ""
	gen := &cobra.Command{
		Use:   "generate --url [--dest]",
		Short: "generate go code from url",
		Args: func(cmd *cobra.Command, args []string) error {
			if _, err := url.Parse(cmd.Flag("url").Value.String()); err != nil {
				return fmt.Errorf("incorrect url format")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, values []string) {
			reader, err := downloader.DownAsReader(from)
			if err != nil {
				log.Fatalf("download %s: %v", from, err)
			}

			if err = os.Mkdir(path.Dir(dest), 0755); err != nil {
				log.Fatalf("make dir %s: %v", dest, err)
			}

			if err = os.WriteFile(dest, []byte(generate.Template(reader)), 0755); err != nil {
				log.Fatalf("generate %s: %v", dest, err)
			}
		},
		Example: "\tgenerate --url xxx [--dest]",
	}

	addURLVar(gen, &from)
	gen.Flags().StringVar(&dest, "dest", "dsg/dst.go", "remote rules url")
	return gen
}

func addURLVar(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVar(value, "url", "", "remote rules url")
	_ = cmd.MarkFlagRequired("url")
}

func checkURL(cmd *cobra.Command, args []string) error {
	if _, err := url.Parse(cmd.Flag("url").Value.String()); err != nil {
		return fmt.Errorf("incorrect url format")
	}

	return nil
}
