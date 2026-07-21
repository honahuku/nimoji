package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hirosassa/nimoji/converter"
)

func main() {
	format := flag.String("format", "google", "output format: google, mac, or msime")
	flag.Parse()

	employees, err := converter.ParseCSV(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch *format {
	case "google":
		err = converter.FormatGoogle(os.Stdout, employees)
	case "mac":
		err = converter.FormatMac(os.Stdout, employees)
	case "msime":
		err = converter.FormatMSIME(os.Stdout, employees)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown format %q (use 'google', 'mac', or 'msime')\n", *format)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
