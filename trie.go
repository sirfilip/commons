package commons

import (
	"strings"
)

var trans = map[string]string{
	" ": "|",
	"|": " ",
}

type Trie interface {
	Add(word string)
	Search(word string) bool
	Autocomplete(prefix string, maxResults int) []string
}

type trieNode struct {
	memory map[string]*trieNode
}

func newTrieNode() *trieNode {
	return &trieNode{memory: make(map[string]*trieNode)}
}

type trie struct {
	root *trieNode
}

func NewTrie() *trie {
	return &trie{root: newTrieNode()}
}

func (t *trie) Add(word string) {
	letters := strings.Split(strings.TrimSpace(word), "")
	curr := t.root
	for i, letter := range letters {
		if letter == "" {
			continue
		}
		isLastLetter := i == len(letters)-1
		curr = t.addLetter(curr, letter, isLastLetter)
	}
}

func (t *trie) Search(word string) bool {
	var found bool
	letters := strings.Split(strings.TrimSpace(word), "")
	curr := t.root
	for i, letter := range letters {
		isLastLetter := i == len(letters)-1
		curr, found = t.search(curr, letter, isLastLetter)
		if !found {
			return false
		}
	}
	return true
}

func (t *trie) Autocomplete(prefix string, maxResults int) []string {
	var suggestions []string
	var found bool

	if maxResults < 1 {
		return suggestions
	}

	curr := t.root
	letters := strings.Split(prefix, "")
	for _, letter := range letters {
		curr, found = t.search(curr, letter, false)
		if !found {
			return suggestions
		}
	}

	return t.autocomplete(curr, prefix, maxResults)
}

func (t *trie) autocomplete(node *trieNode, prefix string, maxResults int) []string {
	var suggestions []string
	stack := []struct {
		node   *trieNode
		prefix string
	}{
		{
			node:   node,
			prefix: prefix,
		},
	}
SuggestionsSearch:
	for len(stack) > 0 {
		curr := stack[0]
		stack = stack[1:]
		for letter, node := range curr.node.memory {
			if letter == "*" {
				suggestions = append(suggestions, curr.prefix)
				if len(suggestions) == maxResults {
					break SuggestionsSearch
				}
				continue
			}
			if translation, ok := trans[letter]; ok {
				letter = translation
			}
			stack = append(stack, struct {
				node   *trieNode
				prefix string
			}{
				node:   node,
				prefix: curr.prefix + letter,
			})
		}
	}
	return suggestions
}

func (t *trie) addLetter(node *trieNode, letter string, isLastLetter bool) *trieNode {
	if translation, ok := trans[letter]; ok {
		letter = translation
	}
	next, ok := node.memory[letter]
	if !ok {
		next = newTrieNode()
		node.memory[letter] = next
	}
	if isLastLetter {
		if _, ok := next.memory["*"]; !ok {
			next.memory["*"] = nil
		}
	}
	return next
}

func (t *trie) search(node *trieNode, letter string, isLastLetter bool) (*trieNode, bool) {
	next, found := node.memory[letter]
	if !found {
		return nil, false
	}
	if isLastLetter {
		_, found = next.memory["*"]
	}
	return next, found
}
