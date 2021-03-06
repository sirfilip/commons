package commons

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrNotFound = errors.New("Not found")

type LRU interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte)
}

type lruNode struct {
	key  string
	val  []byte
	prev *lruNode
	next *lruNode
}

type lru struct {
	capacity int
	cache    map[string]*lruNode
	head     *lruNode
	tail     *lruNode
}

func NewLRU(capacity int) *lru {
	return &lru{capacity: capacity, cache: make(map[string]*lruNode)}
}

func (l lru) Get(key string) ([]byte, error) {
	node, found := l.cache[key]
	if !found {
		return nil, ErrNotFound
	}
	return node.val, nil
}

func (l *lru) Set(key string, val []byte) {
	if l.capacity < 1 {
		return
	}

	node, found := l.cache[key]
	if found {
		l.promote(node)
	} else {
		l.add(&lruNode{key: key, val: val})
	}
}

func (l *lru) promote(node *lruNode) {
	l.remove(node)
	l.add(node)
}

func (l *lru) remove(node *lruNode) {
	delete(l.cache, node.key)
	if l.head.key == node.key {
		l.head = node.next
	}
	if l.tail.key == node.key {
		l.tail = node.prev
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	node.prev = nil
	node.next = nil
}

func (l *lru) add(node *lruNode) {
	if len(l.cache) == l.capacity {
		l.remove(l.tail)
	}
	if len(l.cache) == 0 {
		l.head = node
		l.tail = node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}
	l.cache[node.key] = node
}

type lruNodeSerializer struct {
	Key string
	Val []byte
}

type lruSerializer struct {
	Capacity int
	Cache    []lruNodeSerializer
}

func SerializeLRU(w io.Writer, l *lru) error {
	ls := lruSerializer{Capacity: l.capacity, Cache: make([]lruNodeSerializer, len(l.cache))}

	i := 0
	currentNode := l.head

	for currentNode != nil {
		ls.Cache[i] = lruNodeSerializer{
			Key: currentNode.key,
			Val: currentNode.val,
		}
		i++
		currentNode = currentNode.next
	}

	return json.NewEncoder(w).Encode(ls)
}

func DeserializeLRU(r io.Reader) (*lru, error) {
	ls := lruSerializer{}
	err := json.NewDecoder(r).Decode(&ls)
	if err != nil {
		return nil, err
	}
	l := &lru{capacity: ls.Capacity, cache: make(map[string]*lruNode)}
	if len(ls.Cache) > 0 {
		var prev *lruNode
		for i, item := range ls.Cache {
			node := &lruNode{key: item.Key, val: item.Val}
			l.cache[node.key] = node
			if i == 0 {
				l.head = node
			}
			if i == len(ls.Cache)-1 {
				l.tail = node
			}
			if prev != nil {
				node.prev = prev
				prev.next = node
			}
			prev = node
		}
	}
	return l, nil
}
