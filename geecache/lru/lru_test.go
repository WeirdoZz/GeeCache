package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("weirdo", String("1234"))
	if v, ok := lru.Get("weirdo"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit weirdo=1234 failed")
	}
	if _, ok := lru.Get("peach"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "weirdo", "peach", "bye"
	v1, v2, v3 := "23", "22", "0"
	cap := len(k1 + k2 + v1 + v2)

	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("weirdo"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest weirdo failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("call onEnvicted failed,expect keys equals to %s", expect)
	}
}
