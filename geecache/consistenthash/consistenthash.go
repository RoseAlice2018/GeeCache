package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//Hash maps bytes to uint32
type Hash func(data []byte)uint32

//Map contains all hashed keys
/*
	Map是一致性哈希算法的主数据结构，包含4个成员变量：
	hash:Hash函数
	replicas:虚拟节点倍数
	keys:	hash环
	hashmap: 虚拟节点和真是节点的映射表 键值是虚拟节点的hash，值是真实节点
 */
type Map struct{
	hash Hash
	replicas int
	keys []int //Sorted
	hashMap map[int]string
}

//New creates a Map instance
/*
	Hash采取依赖注入的方式，允许用于替换成自定义的Hash函数 默认crc32.checksumIEEE
 */
func New(replicas int,fn Hash)*Map{
	m := &Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}
/*
	Add 往机器中添加节点

 */
func (m *Map)Add(keys...string){
	for _,key := range keys{
		for i := 0; i < m.replicas ;i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i)+key)))
			m.keys = append(m.keys,hash)
			m.hashMap[hash] = key
		}
	}
}

/*
	选择节点
	Get gets the closest item in the hash to the provided key
 */
func (m *Map)Get(key string)string{
	if len(m.keys)==0{
		return ""
	}

	hash := int(m.hash([]byte(key)))
	//Binary search for appropriate replica.
	idx := sort.Search(len(m.keys),func(i int)bool{
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}





























