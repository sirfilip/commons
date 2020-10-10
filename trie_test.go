package commons

import "testing"

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
			words: []string{"bad", "cat", "catter"},
			expectation: map[string]bool{
				"bad":        true,
				"cat":        true,
				"catter":     true,
				"ca":         false,
				"   cat    ": true,
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
