package converter

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf16"
)

// 129文字以上は取り込みエラーになる。
const commentMaxLength = 128

const bom = "\uFEFF"

// FormatMSIME は MS-IME ユーザー辞書ツールの WORDLIST 形式(UTF-16LE, BOM付き, CRLF)で出力する。
// エラー対応のためコメントの長さ制限・空単語のスキップ・読みのひらがな変換を行う。
func FormatMSIME(w io.Writer, employees []Employee) error {
	entries := make([]DictionaryEntry, 0, len(employees)*3)
	for _, e := range employees {
		fullName := e.LastName + e.FirstName
		comment := fullName
		if e.Note != "" {
			comment = fullName + " / " + e.Note
		}
		entries = append(
			entries,
			DictionaryEntry{Reading: e.LastNameKana, Word: fullName, PartOfSpeech: "固有名詞", Comment: comment},
			DictionaryEntry{Reading: "ばんごう" + e.LastNameKana, Word: e.ID, PartOfSpeech: "固有名詞", Comment: comment},
			DictionaryEntry{Reading: "めーる" + e.LastNameKana, Word: e.Email, PartOfSpeech: "固有名詞", Comment: comment},
		)
	}
	return FormatMSIMEEntries(w, entries)
}

func FormatMSIMEEntries(w io.Writer, entries []DictionaryEntry) error {
	var body bytes.Buffer
	for _, entry := range entries {
		if entry.Word == "" {
			continue
		}
		comment := truncateRunes(entry.Comment, commentMaxLength)
		reading := katakanaToHiragana(entry.Reading)
		if _, err := fmt.Fprintf(&body, "%s\t%s\t%s\t%s\r\n", reading, entry.Word, entry.PartOfSpeech, comment); err != nil {
			return err
		}
	}

	header := bom + "!Microsoft IME Dictionary Tool\r\n" +
		"!Version:\r\n" +
		"!Format:WORDLIST\r\n" +
		"!User Dictionary Name: \r\n" +
		"!Output File Name: \r\n" +
		"!DateTime: \r\n" +
		"\r\n"

	return encodeUTF16LE(w, header+body.String())
}

func encodeUTF16LE(w io.Writer, s string) error {
	for _, u := range utf16.Encode([]rune(s)) {
		if _, err := w.Write([]byte{byte(u & 0xff), byte(u >> 8)}); err != nil {
			return err
		}
	}
	return nil
}

func truncateRunes(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}

// 長音記号「ー」は対応するひらがながないためそのまま残す。
func katakanaToHiragana(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'ァ' && r <= 'ヶ' {
			runes[i] = r - 0x60
		}
	}
	return string(runes)
}
