package commons

import (
	"bufio"
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

func assertLRUEqual(t *testing.T, expected, given *lru) {
	t.Helper()
	if expected.capacity != given.capacity {
		t.Errorf("LRUs are not equal %#v => %#v", expected, given)
	}
	if len(expected.cache) != len(given.cache) {
		t.Errorf("LRUs are not equal %#v => %#v", expected, given)
	}
	expectedCurr := expected.head
	givenCurr := given.head
	for expectedCurr != nil {
		if !bytes.Equal(expectedCurr.val, givenCurr.val) {
			t.Errorf("LRUs are not equal %#v => %#v", expected, given)
		}
		expectedCurr = expectedCurr.next
		givenCurr = givenCurr.next
	}
}

func TestLRUCacheOps(t *testing.T) {
	for title, test := range map[string]struct {
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
	} {
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

func TestLRUCacheSERDE(t *testing.T) {
	for title, test := range map[string]struct {
		lru   LRU
		cache []struct {
			key string
			val []byte
		}
	}{
		"negative capacity SERDE": {
			lru: NewLRU(-1),
			cache: []struct {
				key string
				val []byte
			}{},
		},
		"zero capacity SERDE": {
			lru: NewLRU(0),
			cache: []struct {
				key string
				val []byte
			}{},
		},
		"lru full": {
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
		},
		"lru not filled": {
			lru: NewLRU(2),
			cache: []struct {
				key string
				val []byte
			}{
				{
					key: "one",
					val: []byte("1"),
				},
			},
		},
	} {
		t.Run(title, func(t *testing.T) {
			for _, item := range test.cache {
				test.lru.Set(item.key, item.val)
			}
			var buff bytes.Buffer

			err := SerializeLRU(&buff, test.lru.(*lru))
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}
			r := bufio.NewReader(&buff)
			des, err := DeserializeLRU(r)
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}
			assertLRUEqual(t, test.lru.(*lru), des)
		})
	}
}
