# GO实现的分布式缓存
仿照https://geektutu.com/post/geecache.html进行的书写,对于lru和consistent hash进行了修改。
## LRU
在原本的lru仿照mysql中的设计加入了新老队列的区别
## consistent hash
Consistent hash中参照[跳跃一致性哈希法]("https://writings.sh/post/consistent-hashing-algorithms-part-3-jump-consistent-hash")，进行了设计