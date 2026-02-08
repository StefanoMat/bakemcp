package main

import (
	"flag"
	"fmt"
	"os"

	"your-mcp/internal/cli"
)

func main() {
	var (
		output = flag.String("o", "", "output directory (default: current directory)")
		force  = flag.Bool("f", false, "overwrite non-empty output directory")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bakemcp [options] <openapi-input>\n")
		fmt.Fprintf(os.Stderr, "  openapi-input  path to OpenAPI 3.x file (JSON or YAML)\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := cli.Config{
		InputPath: args[0],
		OutputDir: *output,
		Force:     *force,
	}
	code, err := cli.Run(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}
