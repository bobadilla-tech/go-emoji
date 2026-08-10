package emoji

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Random_ReturnsValidEmoji(t *testing.T) {
	t.Parallel()
	svc := NewService()

	for i := 0; i < 100; i++ {
		e := svc.Random()
		assert.NotEmpty(t, e.Emoji, "expected non-empty emoji from Random()")
		assert.NotEmpty(t, e.Name, "expected non-empty name from Random()")
		assert.NotEmpty(t, e.Category, "expected non-empty category from Random()")
		assert.NotEmpty(t, e.Unicode, "expected non-empty unicode from Random()")
	}
}

func TestService_GetByName(t *testing.T) {
	t.Parallel()
	svc := NewService()

	tests := []struct {
		name      string
		input     string
		wantFound bool
	}{
		{name: "valid name", input: "grinning_face", wantFound: true},
		{name: "valid name uppercase", input: "GRINNING_FACE", wantFound: true},
		{name: "unknown name", input: "not_a_real_emoji", wantFound: false},
		{name: "empty string", input: "", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, found := svc.GetByName(tt.input)
			assert.Equal(t, tt.wantFound, found, "GetByName(%q): got found=%v, want %v", tt.input, found, tt.wantFound)
		})
	}
}

func TestService_GetByName_ReturnsExpectedEmoji(t *testing.T) {
	t.Parallel()
	svc := NewService()

	e, ok := svc.GetByName("grinning_face")
	require.True(t, ok, "expected grinning_face to be found")

	assert.Equal(t, "😀", e.Emoji)
	assert.Equal(t, "grinning_face", e.Name)
	assert.Equal(t, "Smileys & Emotion", e.Category)
	assert.Equal(t, "U+1F600", e.Unicode)
}

func TestService_Search(t *testing.T) {
	t.Parallel()
	svc := NewService()

	tests := []struct {
		name      string
		query     string
		wantEmpty bool
	}{
		{name: "smile query returns results", query: "smile", wantEmpty: false},
		{name: "heart query returns results", query: "heart", wantEmpty: false},
		{name: "nonsense returns empty", query: "zzzyyyxxx", wantEmpty: true},
		{name: "category search", query: "food", wantEmpty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := svc.Search(tt.query)
			isEmpty := result.Total == 0
			assert.Equal(t, tt.wantEmpty, isEmpty, "Search(%q): got empty=%v, want empty=%v (total=%d)", tt.query, isEmpty, tt.wantEmpty, result.Total)
			assert.Equal(t, result.Total, len(result.Items), "Search(%q): Total=%d does not match len(Items)=%d", tt.query, result.Total, len(result.Items))
		})
	}
}

func TestService_Search_CaseInsensitive(t *testing.T) {
	t.Parallel()
	svc := NewService()

	upper := svc.Search("SMILE")
	lower := svc.Search("smile")
	assert.Equal(t, lower.Total, upper.Total, "Search(SMILE)=%d != Search(smile)=%d", upper.Total, lower.Total)
}

func TestService_Search_ReturnsCompleteEmojis(t *testing.T) {
	t.Parallel()
	svc := NewService()

	result := svc.Search("smile")
	require.NotZero(t, result.Total, "expected at least one result for 'smile'")

	for _, e := range result.Items {
		assert.True(t, e.Emoji != "" && e.Name != "" && e.Category != "" && e.Unicode != "", "expected all fields to be non-empty for emoji: %+v", e)
	}
}

func TestService_Search_EmptyResults(t *testing.T) {
	t.Parallel()
	svc := NewService()

	result := svc.Search("zzzyyyxxx")
	assert.Equal(t, 0, result.Total, "expected 0 results for nonsense query, got %d", result.Total)
	assert.NotNil(t, result.Items, "expected non-nil items slice for empty results")
}

func TestService_ByCategory(t *testing.T) {
	t.Parallel()
	svc := NewService()

	category := "Animals & Nature"
	result := svc.ByCategory(category)
	assert.NotZero(t, result.Total, "expected at least one emoji in %q", category)
	assert.Equal(t, result.Total, len(result.Items), "Total=%d does not match len(Items)=%d", result.Total, len(result.Items))

	for _, e := range result.Items {
		assert.Equal(t, category, e.Category, "expected emoji in category %q, got %q", category, e.Category)
	}
}

func TestService_ByCategory_CaseInsensitive(t *testing.T) {
	t.Parallel()
	svc := NewService()

	upper := svc.ByCategory("ANIMALS & NATURE")
	title := svc.ByCategory("Animals & Nature")
	assert.Equal(t, title.Total, upper.Total, "ByCategory(uppercase)=%d != ByCategory(title)=%d", upper.Total, title.Total)
}

func TestService_ByCategory_NoMatch(t *testing.T) {
	t.Parallel()
	svc := NewService()

	result := svc.ByCategory("not_a_category")
	assert.Equal(t, 0, result.Total, "expected 0 results for unknown category, got %d", result.Total)
	assert.NotNil(t, result.Items, "expected non-nil items slice for empty results")
}

func TestService_Categories(t *testing.T) {
	t.Parallel()
	svc := NewService()

	categories := svc.Categories()
	assert.Equal(t, 8, len(categories), "expected 8 categories, got %d: %v", len(categories), categories)

	seen := make(map[string]bool, len(categories))
	for _, c := range categories {
		assert.NotEmpty(t, c, "expected non-empty category")
		assert.False(t, seen[c], "duplicate category %q", c)
		seen[c] = true
	}

	assert.Contains(t, categories, "Symbols")
}

func TestDatasetIntegrity(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 808, len(emojis), "embedded dataset count")

	names := make(map[string]bool, len(emojis))
	for _, e := range emojis {
		assert.True(t, e.Emoji != "" && e.Name != "" && e.Category != "" && e.Unicode != "", "expected all fields to be non-empty for emoji: %+v", e)
		assert.False(t, names[e.Name], "duplicate name %q", e.Name)
		names[e.Name] = true
	}
}
