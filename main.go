package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hirosassa/nimoji/converter"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, in io.Reader, out, errOut io.Writer) error {
	flags := flag.NewFlagSet("nimoji", flag.ContinueOnError)
	flags.SetOutput(errOut)
	from := flags.String("from", "csv", "input format: csv or google")
	format := flags.String("format", "google", "output format: google, mac, or msime")
	if err := flags.Parse(args); err != nil {
		return err
	}

	switch *from {
	case "csv":
		employees, err := converter.ParseCSV(in)
		if err != nil {
			return err
		}
		switch *format {
		case "google":
			return converter.FormatGoogle(out, employees)
		case "mac":
			return converter.FormatMac(out, employees)
		case "msime":
			return converter.FormatMSIME(out, employees)
		default:
			return fmt.Errorf("unknown format %q (use 'google', 'mac', or 'msime')", *format)
		}
	case "google":
		if *format != "msime" {
			return fmt.Errorf("input format %q only supports output format 'msime'", *from)
		}
		entries, err := converter.ParseGoogleTSV(in)
		if err != nil {
			return err
		}
		return converter.FormatMSIMEEntries(out, entries)
	default:
		return fmt.Errorf("unknown input format %q (use 'csv' or 'google')", *from)
	}
}
