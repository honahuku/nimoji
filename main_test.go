package main

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestRunGoogleTSVToMSIME(t *testing.T) {
	input := "テラダ\t寺田親弘\t固有名詞\t寺田親弘 / 開発部\n"
	var output bytes.Buffer
	var errOutput bytes.Buffer

	err := run(
		[]string{"-from", "google", "-format", "msime"},
		strings.NewReader(input),
		&output,
		&errOutput,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, output.Bytes())
	want := "てらだ\t寺田親弘\t固有名詞\t寺田親弘 / 開発部"
	if !strings.Contains(got, want) {
		t.Errorf("expected converted entry %q, got:\n%s", want, got)
	}
}

func TestRunRejectsGoogleTSVWithNonMSIMEOutput(t *testing.T) {
	err := run(
		[]string{"-from", "google", "-format", "google"},
		strings.NewReader("たなか\t田中太郎\t固有名詞\n"),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)
	if err == nil {
		t.Fatal("expected an error for unsupported format combination")
	}
}

func decodeUTF16LE(t *testing.T, data []byte) string {
	t.Helper()
	if len(data) < 2 || len(data)%2 != 0 {
		t.Fatalf("invalid UTF-16LE data length: %d", len(data))
	}

	codeUnits := make([]uint16, 0, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		codeUnits = append(codeUnits, uint16(data[i])|uint16(data[i+1])<<8)
	}
	return string(utf16.Decode(codeUnits))
}
