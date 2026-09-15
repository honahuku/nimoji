package converter

import (
	"strings"
	"testing"
)

func TestParseGoogleTSV(t *testing.T) {
	input := "テラダ\t寺田親弘\t固有名詞\t寺田親弘 / 開発部\n" +
		"ばんごうテラダ\t001\t固有名詞\t寺田親弘 / 開発部\n"

	entries, err := ParseGoogleTSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	expected := DictionaryEntry{
		Reading:      "テラダ",
		Word:         "寺田親弘",
		PartOfSpeech: "固有名詞",
		Comment:      "寺田親弘 / 開発部",
	}
	if entries[0] != expected {
		t.Errorf("expected %#v, got %#v", expected, entries[0])
	}
}

func TestParseGoogleTSV_AllowsMissingComment(t *testing.T) {
	entries, err := ParseGoogleTSV(strings.NewReader("たなか\t田中太郎\t固有名詞\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Comment != "" {
		t.Errorf("expected empty comment, got %q", entries[0].Comment)
	}
}

func TestParseGoogleTSV_RejectsInvalidColumns(t *testing.T) {
	_, err := ParseGoogleTSV(strings.NewReader("たなか\t田中太郎\n"))
	if err == nil {
		t.Fatal("expected an error for invalid columns")
	}
}

func TestFormatMSIMEEntries(t *testing.T) {
	entries := []DictionaryEntry{
		{
			Reading:      "テラダ",
			Word:         "寺田親弘",
			PartOfSpeech: "固有名詞",
			Comment:      "寺田親弘 / 開発部",
		},
	}

	var buf strings.Builder
	if err := FormatMSIMEEntries(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := decodeUTF16LE(t, []byte(buf.String()))
	want := "てらだ\t寺田親弘\t固有名詞\t寺田親弘 / 開発部"
	if !strings.Contains(got, want) {
		t.Errorf("expected converted entry %q, got:\n%s", want, got)
	}
}
