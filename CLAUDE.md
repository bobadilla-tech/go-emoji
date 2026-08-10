# go-emoji

Minimal Go library for emoji lookup and search, backed by an embedded Unicode
CLDR dataset in `data/emoji.json`. Zero runtime dependencies; `testify` is a
test-only dependency. Public API is `emoji.NewService()` with `Random`,
`GetByName`, `Search`, `ByCategory`, and `Categories`.

## Layout

- Root module (`github.com/bobadilla-tech/go-emoji`) — the library, zero
  runtime dependencies.
- `data/` — the embedded dataset: `data/emoji.json` plus a one-line
  `//go:embed` in `data/data.go`. `data` exposes the raw bytes; the root
  `emoji` package parses them once at init.

## Commands

```sh
go test ./...                # run tests (root module)
go test ./... -covermode=atomic -coverprofile=coverage.out && go tool cover -func=coverage.out
go vet ./...
gofmt -l .
```

Adding or bumping dependencies: `go get` then `go mod tidy` — keep `go.sum`
checked in.

No Makefile, no lint config beyond `go vet` + `gofmt` — matches the rest of the
bobadilla-tech Go packages.

## Updating the dataset

Editing emoji entries is a data PR: modify `data/emoji.json`, run
`go test ./...` (the dataset-integrity test enforces the record count and
non-empty fields), and update the count asserted in `emoji_test.go` if entries
were added or removed.
