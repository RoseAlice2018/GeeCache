package lru

import "container/list"

//Cache is a LRU cache.It is not safe for concurrent access
type Cache struct{
	maxBytes int64 //允许使用的最大内存
	nbytes  int64 //当前已经使用的内存
	ll	 	*list.List //Go语言标准库实现的双向链表
	cache map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数
}
/*
	键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
 */
type entry struct{
	key string
	value Value
}
/*
为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
 */
//Value use len to count how many bytes it takes
type Value interface {
	Len() int
}

//New is the Constructor of Cache
func New(maxBytes int64,onEvicted func(string,Value))*Cache{
	return &Cache{
		maxBytes:maxBytes,
		ll: list.New(),
		cache : make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
/*
查找功能：
第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾。
c.ll.MoveToFront(ele)，即将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）
 */
func (c *Cache)Get(key string)(value Value,ok bool){
	if ele,ok := c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value,true
	}
	return
}
/*
删除 也就是缓存淘汰
RemoveOldest removes the oldest item
 */
func (c *Cache)RemoveOldest(){
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key,kv.value)
		}
	}
}
/*
	新增/修改
	Add a new value to the cache
 */
func (c *Cache)Add(key string,value Value){
	if ele,ok := c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}else{
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes !=0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//Len the number of cache entries
func (c *Cache)Len() int{
	return c.ll.Len()
}