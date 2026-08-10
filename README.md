# go-emoji

[![CI](https://github.com/bobadilla-tech/go-emoji/actions/workflows/ci.yml/badge.svg)](https://github.com/bobadilla-tech/go-emoji/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/bobadilla-tech/go-emoji/graph/badge.svg)](https://codecov.io/gh/bobadilla-tech/go-emoji)
[![Go Reference](https://pkg.go.dev/badge/github.com/bobadilla-tech/go-emoji.svg)](https://pkg.go.dev/github.com/bobadilla-tech/go-emoji)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Zero-dependency Go library for emoji lookup and search, backed by an embedded
Unicode CLDR dataset. Library only — no CLI. The only external module is
`testify`, used exclusively in tests.

The dataset lives in `data/emoji.json` (embedded via `//go:embed`) so emoji
updates are data changes, not code changes. Structure inspired by
[`diamondburned/go-emoji`](https://github.com/diamondburned/go-emoji).

## Install

```sh
go get github.com/bobadilla-tech/go-emoji
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/bobadilla-tech/go-emoji"
)

func main() {
	svc := emoji.NewService()

	fmt.Println(svc.GetByName("grinning_face")) // {😀 grinning_face Smileys & Emotion U+1F600} true
	fmt.Println(svc.Random().Emoji)

	results := svc.Search("heart")
	fmt.Printf("%d matches\n", results.Total)

	animals := svc.ByCategory("Animals & Nature")
	fmt.Printf("%d animals\n", animals.Total)

	fmt.Println(svc.Categories())
}
```

### API

- `Service.Random() Emoji` — random emoji (crypto/rand).
- `Service.GetByName(name string) (Emoji, bool)` — exact name match, case-insensitive.
- `Service.Search(query string) List` — case-insensitive substring match on name or category.
- `Service.ByCategory(category string) List` — case-insensitive exact category match.
- `Service.Categories() []string` — the distinct categories, in order of first appearance.
- `data.RawEmojiJSON []byte` — the embedded raw JSON.

## Data

`data/emoji.json` holds 808 records, each with `emoji` (glyph), `name` (CLDR
snake_case short name), `category` (Unicode category), and `unicode` (U+XXXX
notation). Derived from the Unicode CLDR / unicode.org emoji specification.

## License

[MIT](LICENSE). Data derived from the Unicode CLDR emoji specification, which
is distributed under the Unicode Terms of Use.
