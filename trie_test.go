package commons

import (
	"strings"
	"testing"
)

func assertSameSuggestions(t *testing.T, expect, got []string) {
	t.Helper()
	if len(expect) != len(got) {
		t.Errorf("Expected %#v but got %#v", expect, got)
		return
	}
	set := make(map[string]struct{})
	for _, word := range expect {
		set[word] = struct{}{}
	}
	for _, word := range got {
		if _, ok := set[word]; !ok {
			t.Errorf("Expected %#v but got %#v", expect, got)
			return

		}
	}
}

func TestTrie_Search(t *testing.T) {
	for title, test := range map[string]struct {
		words       []string
		expectation map[string]bool
	}{
		"it finds the correct words": {
			words: []string{"bad", "cat", "catter", "a phrase"},
			expectation: map[string]bool{
				"bad":        true,
				"cat":        true,
				"catter":     true,
				"ca":         false,
				"   cat    ": true,
				"a phrase":   true,
			},
		},
	} {
		t.Run(title, func(t *testing.T) {
			tr := NewTrie()
			for _, word := range test.words {
				tr.Add(word)
			}
			for word, exp := range test.expectation {
				res := tr.Search(word)
				if res != exp {
					t.Errorf("Expectation failed: %v for %s", exp, word)
				}
			}
		})
	}
}

func TestTrie_Autocomplete(t *testing.T) {
	for title, test := range map[string]struct {
		memory     []string
		prefix     string
		maxResults int
		expect     []string
	}{
		"fetches all results": {
			memory:     []string{"bad", "bat", "better", "cat"},
			prefix:     "b",
			maxResults: 3,
			expect:     []string{"bad", "bat", "better"},
		},
		"returns only max items": {
			memory:     []string{"bad", "bat", "better", "cat"},
			prefix:     "b",
			maxResults: 2,
			expect:     []string{"bad", "bat"},
		},
		"no matches found": {
			memory:     []string{"bad", "bat", "better", "cat"},
			prefix:     "other",
			maxResults: 2,
			expect:     []string{},
		},
		"supports phrases": {
			memory:     []string{"nice", "nice weather", "other"},
			prefix:     "ni",
			maxResults: 100,
			expect:     []string{"nice", "nice weather"},
		},
		"works with negative max results": {
			memory:     []string{"nice", "nice weather", "other"},
			prefix:     "ni",
			maxResults: -1,
			expect:     []string{},
		},
	} {
		t.Run(title, func(t *testing.T) {
			tr := NewTrie()
			for _, word := range test.memory {
				tr.Add(word)
			}
			got := tr.Autocomplete(test.prefix, test.maxResults)
			assertSameSuggestions(t, test.expect, got)
		})
	}
}

func BenchmarkTrie_Autocomplete(b *testing.B) {
	var r []string
	tr := NewTrie()
	for i := 2; i < 100; i++ {
		tr.Add(strings.Repeat("b", i))
		tr.Add("b" + strings.Repeat("a", i))
		tr.Add("b" + strings.Repeat("c", i))
		tr.Add("b" + strings.Repeat("d", i))
		tr.Add("b" + strings.Repeat("e", i))
		tr.Add("b" + strings.Repeat("f", i))
		tr.Add("b" + strings.Repeat("g", i))
		tr.Add("b" + strings.Repeat("h", i))
		tr.Add("b" + strings.Repeat("i", i))
		tr.Add("b" + strings.Repeat("j", i))
		tr.Add("b" + strings.Repeat("k", i))
		tr.Add("b" + strings.Repeat("l", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = tr.Autocomplete("b", 10000)
		for range r {
		}
	}
}
