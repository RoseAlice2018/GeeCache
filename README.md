# GeeCache

##Day 2
1. 实现了LRU缓存的并发控制
2. 实现了GeeCache核心数据结构Group，缓存不存在时，调用回调函数
获取源数据。
##Day 3
1. 利用Go语言标准库http搭建HTTP Server
2. 用main函数启动HTTP Server测试API
3. 实现一致性hash代码
##Day 4
1. 注册节点，借助一致性Hash算法选择节点
2. 实现HTTP客户端，与远程节点的服务端通信
##Day 5
1. 防止缓存击穿
（实现方法：在一瞬间有大量请求get(key)，而且key未被缓存或者未被缓存在当前节点 如果不用singleflight，那么这些请求都会发送远端节点或者从本地数据库读取，会造成远端节点或本地数据库压力猛增。使用singleflight，第一个get(key)请求到来时，singleflight会记录当前key正在被处理，后续的请求只需要等待第一个请求处理完成，取返回值即可。）
## To-do List

- [ ] 单机缓存和基于HTTP的分布式缓存
- [ ] 最近最少访问（Least Recently Used,LRU)缓存策略
- [ ] 使用Go锁机制防止缓存击穿
- [ ] 使用一致性hash选择节点，实现负载均衡