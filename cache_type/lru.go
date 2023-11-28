package cache_type

import (
	"container/list"
)

type LRU struct {
	capacity int64 // in bytes
	size     int64
	dict     map[string]*list.Element
	list     *list.List
	onEvit   func(key string, value Value)
}

func NewLRU(capacity int64, onEvicted func(string, Value)) *LRU {
	return &LRU{capacity: capacity,
		size:   0,
		dict:   make(map[string]*list.Element),
		list:   list.New(),
		onEvit: onEvicted,
	}
}

type Value interface {
	Len() int
}

// Len the number of cache entries
func (c *LRU) Len() int {
	return c.list.Len()
}

type Node struct {
	key string
	val Value
}

func (cache *LRU) Get(key string) (Value, bool) {
	if node, flag := cache.dict[key]; flag {
		cache.list.MoveToBack(node)
		return node.Value.(*Node).val, true
	}
	return nil, false
}

func (cache *LRU) Add(key string, val Value) {
	if node, flag := cache.dict[key]; flag {
		// key exist
		cache.size -= int64(node.Value.(*Node).val.Len())
		node.Value.(*Node).val = val
		cache.size += int64(val.Len())
		cache.list.MoveToBack(node)
	} else {
		// key not exist
		cache.size += int64(val.Len()) + int64(len(key))
		newNode := &Node{key: key, val: val}
		cache.dict[key] = cache.list.PushBack(newNode)
	}
	// pop least recent if over capacity
	for cache.size > cache.capacity {
		removedNode := cache.list.Front()
		if removedNode != nil {
			kv := removedNode.Value.(*Node)
			delete(cache.dict, kv.key)
			cache.list.Remove(removedNode)
			cache.size -= int64(kv.val.Len()) + int64(len(kv.key))
			if cache.onEvit != nil {
				cache.onEvit(kv.key, kv.val)
			}
		}
	}
}
