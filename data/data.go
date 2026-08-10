// Package data embeds the emoji dataset as a JSON blob.
//
// The dataset is derived from the Unicode CLDR / unicode.org emoji
// specification. See the go-emoji README for provenance and update notes.
package data

import _ "embed"

//go:embed emoji.json
var RawEmojiJSON []byte
