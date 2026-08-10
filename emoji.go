// Package emoji provides emoji lookup, search, and category queries backed by
// the embedded Unicode CLDR emoji dataset.
package emoji

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/bobadilla-tech/go-emoji/data"
)

// Emoji represents a single emoji with its metadata.
type Emoji struct {
	Emoji    string `json:"emoji"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Unicode  string `json:"unicode"`
}

// List represents a collection of emoji search results.
type List struct {
	Items []Emoji `json:"items"`
	Total int     `json:"total"`
}

// emojis is the full dataset, loaded once from the embedded JSON at init.
var emojis = mustLoad(data.RawEmojiJSON)

func mustLoad(raw []byte) []Emoji {
	var es []Emoji
	if err := json.Unmarshal(raw, &es); err != nil {
		panic(fmt.Sprintf("go-emoji: parse embedded emoji.json: %v", err))
	}
	return es
}

// Service provides emoji lookup and search operations.
type Service struct{}

// NewService returns a new Service.
func NewService() *Service { return &Service{} }

// Random returns a randomly selected emoji using a cryptographically secure
// random number generator.
func (s *Service) Random() Emoji {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(emojis))))
	if err != nil {
		// Fallback to first emoji on the extremely unlikely crypto/rand failure.
		return emojis[0]
	}
	return emojis[n.Int64()]
}

// GetByName returns the emoji matching the given name (snake_case).
// Returns the emoji and true if found, or a zero value and false if not.
func (s *Service) GetByName(name string) (Emoji, bool) {
	name = strings.ToLower(name)
	for _, e := range emojis {
		if e.Name == name {
			return e, true
		}
	}
	return Emoji{}, false
}

// Search returns all emojis whose name or category contains the query string
// (case-insensitive). Returns a List with matching results.
func (s *Service) Search(query string) List {
	query = strings.ToLower(query)
	matches := make([]Emoji, 0)
	for _, e := range emojis {
		if strings.Contains(e.Name, query) || strings.Contains(strings.ToLower(e.Category), query) {
			matches = append(matches, e)
		}
	}
	return List{Items: matches, Total: len(matches)}
}

// ByCategory returns all emojis in the given category (case-insensitive).
// Returns a List with matching results.
func (s *Service) ByCategory(category string) List {
	category = strings.ToLower(category)
	matches := make([]Emoji, 0)
	for _, e := range emojis {
		if strings.ToLower(e.Category) == category {
			matches = append(matches, e)
		}
	}
	return List{Items: matches, Total: len(matches)}
}

// Categories returns the distinct emoji categories in the dataset, in order
// of first appearance.
func (s *Service) Categories() []string {
	seen := make(map[string]struct{}, 8)
	cat := make([]string, 0, 8)
	for _, e := range emojis {
		if _, ok := seen[e.Category]; ok {
			continue
		}
		seen[e.Category] = struct{}{}
		cat = append(cat, e.Category)
	}
	return cat
}
