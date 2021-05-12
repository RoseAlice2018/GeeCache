# GeeCache

##Day 2
1. 实现了LRU缓存的并发控制
2. 实现了GeeCache核心数据结构Group，缓存不存在时，调用回调函数
获取源数据。

## To-do List

- [ ] 单机缓存和基于HTTP的分布式缓存
- [ ] 最近最少访问（Least Recently Used,LRU)缓存策略
- [ ] 使用Go锁机制防止缓存击穿
- [ ] 使用一致性hash选择节点，实现负载均衡