package converter

import (
	"encoding/xml"
	"fmt"
	"io"
)

type plistDict struct {
	Phrase   string
	Shortcut string
}

func FormatMac(w io.Writer, employees []Employee) error {
	var entries []plistDict
	for _, e := range employees {
		entries = append(
			entries,
			plistDict{Phrase: e.LastName + e.FirstName, Shortcut: e.LastNameKana},
			plistDict{Phrase: e.ID, Shortcut: "ばんごう" + e.LastNameKana},
			plistDict{Phrase: e.Email, Shortcut: "めーる" + e.LastNameKana},
		)
	}

	fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(w, `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`)
	fmt.Fprintln(w, `<plist version="1.0">`)
	fmt.Fprintln(w, `<array>`)

	for _, entry := range entries {
		fmt.Fprintln(w, "\t<dict>")
		fmt.Fprintln(w, "\t\t<key>phrase</key>")
		fmt.Fprintf(w, "\t\t<string>%s</string>\n", xmlEscape(entry.Phrase))
		fmt.Fprintln(w, "\t\t<key>shortcut</key>")
		fmt.Fprintf(w, "\t\t<string>%s</string>\n", xmlEscape(entry.Shortcut))
		fmt.Fprintln(w, "\t</dict>")
	}

	fmt.Fprintln(w, `</array>`)
	fmt.Fprintln(w, `</plist>`)
	return nil
}

func xmlEscape(s string) string {
	var buf []byte
	_ = xml.EscapeText((*xmlWriter)(&buf), []byte(s))
	return string(buf)
}

type xmlWriter []byte

func (w *xmlWriter) Write(p []byte) (int, error) {
	*w = append(*w, p...)
	return len(p), nil
}
