package converter

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestFormatMSIME(t *testing.T) {
	employees := []Employee{
		{
			ID:            "001",
			LastName:      "田中",
			FirstName:     "太郎",
			Email:         "tanaka@example.com",
			LastNameKana:  "たなか",
			FirstNameKana: "たろう",
		},
	}

	var buf bytes.Buffer
	if err := FormatMSIME(&buf, employees); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, buf.Bytes())
	lines := strings.Split(got, "\r\n")

	wantHeader := []string{
		"\uFEFF!Microsoft IME Dictionary Tool",
		"!Version:",
		"!Format:WORDLIST",
	}
	for i, want := range wantHeader {
		if lines[i] != want {
			t.Errorf("line %d: expected %q, got %q", i, want, lines[i])
		}
	}
	if lines[6] != "" {
		t.Errorf("line 6: expected blank separator line, got %q", lines[6])
	}

	wantEntries := []string{
		"たなか\t田中太郎\t固有名詞\t田中太郎",
		"ばんごうたなか\t001\t固有名詞\t田中太郎",
		"めーるたなか\ttanaka@example.com\t固有名詞\t田中太郎",
	}
	for i, want := range wantEntries {
		if lines[7+i] != want {
			t.Errorf("entry line %d: expected %q, got %q", i, want, lines[7+i])
		}
	}

	if !strings.HasSuffix(got, "\r\n") {
		t.Errorf("expected output to end with CRLF, got %q", got)
	}
}

func TestFormatMSIME_ProducesUTF16LEWithBOM(t *testing.T) {
	employees := []Employee{
		{ID: "001", LastName: "田中", FirstName: "太郎", Email: "tanaka@example.com", LastNameKana: "たなか"},
	}

	var buf bytes.Buffer
	if err := FormatMSIME(&buf, employees); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data := buf.Bytes()
	if len(data) < 2 || data[0] != 0xff || data[1] != 0xfe {
		t.Fatalf("expected UTF-16LE BOM (0xff 0xfe), got % x", data[:2])
	}
}

func TestFormatMSIME_ConvertsKatakanaReadingToHiragana(t *testing.T) {
	employees := []Employee{
		{ID: "001", LastName: "寺田", FirstName: "親弘", Email: "terada@example.com", LastNameKana: "テラダ"},
	}

	var buf bytes.Buffer
	if err := FormatMSIME(&buf, employees); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, buf.Bytes())
	if !strings.Contains(got, "てらだ\t寺田親弘\t固有名詞\t寺田親弘") {
		t.Errorf("expected hiragana reading, got:\n%s", got)
	}
	if strings.Contains(got, "テラダ\t") {
		t.Errorf("katakana reading should have been converted, got:\n%s", got)
	}
}

func TestFormatMSIME_KeepsChoonmarkInReading(t *testing.T) {
	employees := []Employee{
		{ID: "001", LastName: "ロマン", FirstName: "ポール", Email: "roman@example.com", LastNameKana: "ローマン"},
	}

	var buf bytes.Buffer
	if err := FormatMSIME(&buf, employees); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, buf.Bytes())
	if !strings.Contains(got, "ろーまん\t") {
		t.Errorf("expected choonmark 'ー' to be kept as-is, got:\n%s", got)
	}
}

func TestFormatMSIME_TruncatesCommentExceedingLimit(t *testing.T) {
	longNote := strings.Repeat("あ", 200)
	employees := []Employee{
		{ID: "001", LastName: "田中", FirstName: "太郎", Email: "tanaka@example.com", LastNameKana: "たなか", Note: longNote},
	}

	var buf bytes.Buffer
	if err := FormatMSIME(&buf, employees); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, buf.Bytes())
	for _, line := range strings.Split(got, "\r\n") {
		fields := strings.Split(line, "\t")
		if len(fields) < 4 {
			continue
		}
		comment := []rune(fields[3])
		if len(comment) > commentMaxLength {
			t.Errorf("comment exceeds %d runes (%d): %q", commentMaxLength, len(comment), fields[3])
		}
	}
}

func TestFormatMSIME_SkipsEntryWithEmptyWord(t *testing.T) {
	employees := []Employee{
		{ID: "001", LastName: "友近", FirstName: "玲也", Email: "", LastNameKana: "トモチカ"},
	}

	var buf bytes.Buffer
	if err := FormatMSIME(&buf, employees); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, buf.Bytes())
	if strings.Contains(got, "めーる") {
		t.Errorf("expected entry with empty word to be skipped, got:\n%s", got)
	}
	// 姓名・社員番号のエントリは残る
	if !strings.Contains(got, "ともちか\t友近玲也\t固有名詞\t友近玲也") {
		t.Errorf("expected name entry to remain, got:\n%s", got)
	}
}

// decodeUTF16LE は UTF-16LE (BOM 付き) バイト列を UTF-8 文字列にデコードするテスト用ヘルパー。
func decodeUTF16LE(t *testing.T, data []byte) string {
	t.Helper()
	if len(data) < 2 {
		t.Fatalf("data too short to contain BOM: % x", data)
	}
	if len(data)%2 != 0 {
		t.Fatalf("odd-length UTF-16LE data: %d bytes", len(data))
	}
	u16 := make([]uint16, 0, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		u16 = append(u16, uint16(data[i])|uint16(data[i+1])<<8)
	}
	// BOM (U+FEFF) はそのまま文字列先頭に残す（テストの期待値側もそれを含めて比較する）。
	return string(utf16.Decode(u16))
}
