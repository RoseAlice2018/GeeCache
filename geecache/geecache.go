package geecache

import (
	"fmt"
	"log"
	"sync"
	"./singleflight"
)

//A Getter loads data for a key
type Getter interface {
	Get(key string)([]byte,error)
}

//A GetterFunc implements Getter with a function
type GetterFunc func(key string)([]byte,error)

//Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
/*
	1.name 一个Group对应的name，唯一标志Group
	2.getter getter 缓存未命中时获取源数据的回调（callback）
	3.mainCache 一开始实现的并发缓存
 */
type Group struct {
	name string
	getter Getter
	mainCache cache
	peers PeerPicker
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
/*
	实例化Group
	并且将Group存储到全局变量中
 */
func NewGroup(name string,cacheBytes int64,getter Getter)*Group{
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return  g
}


// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
/*
	GetGroup 用来特定名称的 Group，这里使用了只读锁 RLock()，因为不涉及任何冲突变量的写操作。
 */
func GetGroup(name string)*Group{
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}


//Get value for a key from cache
/*
	从mainCache中查找缓存，如果存在则返回内存值
	如果不存在，则调用Load方法
 */
func (g *Group)Get(key string)(ByteView,error){
	if key == "" {
		return ByteView{},fmt.Errorf("key is required")
	}

	if v,ok := g.mainCache.get(key);ok {
		log.Println("[GeeCache] hit")
		return v,nil
	}
	return g.load(key)
}
/*
	调用用户回调函数，获取源数据，并且将源数据添加到缓存mainCache中（通过populateCache方法）
 */
func (g *Group)getLocally(key string)(ByteView,error){
	bytes,err := g.getter.Get(key)
	if err != nil{
		return ByteView{}, err
	}
	value := ByteView{b:cloneBytes(bytes)}
	g.populateCache(key,value)
	return value,nil
}

func (g *Group)populateCache(key string,value ByteView){
	g.mainCache.add(key,value)
}
//RegisterPeers register a PeerPicker for choosing remote peer
func (g *Group)RegisterPeers(peers PeerPicker){
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}
/*
	load调用getlocally从本地调用源数据
*/
func (g *Group) load(key string) (value ByteView, err error) {
	// each key is only fetched once (either locally or remotely)
	// regardless of the number of concurrent callers.
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}




























