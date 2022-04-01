package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nbytes   int64
	//实际上保存的是节点的访问信息
	ll *list.List
	//这才是保存的真正的节点信息
	cache map[string]*list.Element
	// 可选的，当一个entry被清除的时候会执行
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value 接口类型，用Len来计算它携带了多少byte的数据
type Value interface {
	Len() int
}

// New 构建一个Cache实例
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Get 查找某个键对应的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除掉队尾的节点
func (c *Cache) RemoveOldest() {
	//获取应该移除的节点（即队尾节点）
	ele := c.ll.Back()
	if ele != nil {
		//首先从队列中移除节点
		c.ll.Remove(ele)
		//然后从获取该节点对应的entry
		kv := ele.Value.(*entry)
		//从cache中删除掉这个节点
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// Add 增加键值对映射
func (c *Cache) Add(key string, value Value) {
	// 如果已经存在该键
	if ele, ok := c.cache[key]; ok {
		// 将键移到队首
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 更新键对应的值
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 如果该键不存在，需要在list和cache中都增加一个新的键值对
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
