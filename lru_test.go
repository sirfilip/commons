package commons

import (
	"bytes"
	"errors"
	"testing"
)

func assertEqualBytes(t *testing.T, expected, given []byte) {
	t.Helper()
	if expected == nil && given != nil {
		t.Errorf("expected bytes missmatch %#v => %#v", expected, given)
	}
	if !bytes.Equal(expected, given) {
		t.Errorf("expected bytes missmatch %#v => %#v", expected, given)
	}
}

func TestLRUCache(t *testing.T) {
	tests := map[string]struct {
		lru   LRU
		cache []struct {
			key string
			val []byte
		}
		expect []struct {
			key string
			val []byte
			err error
		}
	}{
		"negative capacity must not keep any items": {
			lru: NewLRU(-1),
			cache: []struct {
				key string
				val []byte
			}{
				{
					key: "one",
					val: []byte("1"),
				},
				{
					key: "two",
					val: []byte("2"),
				},
			},
			expect: []struct {
				key string
				val []byte
				err error
			}{
				{
					key: "one",
					val: nil,
					err: ErrNotFound,
				},
			},
		},
		"zero capacity must not keep any items": {
			lru: NewLRU(0),
			cache: []struct {
				key string
				val []byte
			}{
				{
					key: "one",
					val: []byte("1"),
				},
				{
					key: "two",
					val: []byte("2"),
				},
			},
			expect: []struct {
				key string
				val []byte
				err error
			}{
				{
					key: "one",
					val: nil,
					err: ErrNotFound,
				},
			},
		},
		"must respect capacity": {
			lru: NewLRU(1),
			cache: []struct {
				key string
				val []byte
			}{
				{
					key: "one",
					val: []byte("1"),
				},
				{
					key: "two",
					val: []byte("2"),
				},
			},
			expect: []struct {
				key string
				val []byte
				err error
			}{
				{
					key: "one",
					val: nil,
					err: ErrNotFound,
				},
				{
					key: "two",
					val: []byte("2"),
					err: nil,
				},
			},
		},
		"must respect lru": {
			lru: NewLRU(2),
			cache: []struct {
				key string
				val []byte
			}{
				{
					key: "one",
					val: []byte("1"),
				},
				{
					key: "two",
					val: []byte("2"),
				},
				{
					key: "one",
					val: []byte("1"),
				},
				{
					key: "three",
					val: []byte("3"),
				},
			},
			expect: []struct {
				key string
				val []byte
				err error
			}{
				{
					key: "one",
					val: []byte("1"),
					err: nil,
				},
				{
					key: "two",
					val: nil,
					err: ErrNotFound,
				},
				{
					key: "three",
					val: []byte("3"),
					err: nil,
				},
			},
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			for _, item := range test.cache {
				test.lru.Set(item.key, item.val)
			}
			for _, expect := range test.expect {
				val, err := test.lru.Get(expect.key)
				if !errors.Is(expect.err, err) {
					t.Errorf("expected errors missmatch: %v => %v", expect.err, err)
				}
				assertEqualBytes(t, expect.val, val)
			}
		})
	}
}
