package api

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type CLIOptions struct {
	Command string
	File    string
	Batch   string
	Query   string
}

func ParseCLI(args []string) CLIOptions {
	flags := flag.NewFlagSet("telemetryd", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	options := CLIOptions{}
	flags.StringVar(&options.Command, "command", "serve", "command")
	flags.StringVar(&options.File, "file", "", "import file")
	flags.StringVar(&options.Batch, "batch", "cli", "batch id")
	flags.StringVar(&options.Query, "query", "", "search query")
	_ = flags.Parse(args)
	return options
}

func ValidateCLI(options CLIOptions) error {
	switch strings.ToLower(options.Command) {
	case "serve", "import", "report", "search":
		return nil
	default:
		return fmt.Errorf("unsupported command %q", options.Command)
	}
}

func Usage() string {
	return "telemetryd -command serve|import|report|search [-file path] [-batch id] [-query text]"
}
