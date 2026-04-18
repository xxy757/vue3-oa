import (
	"sync"
	"time"
)

// item 表示一个缓存条目及其关联的过期时间。
type item struct {
	value interface{}
	value     interface{}

// MemoryCache 是 Cache 接口的内存实现。
// 使用 sync.RWMutex 保证线程安全的并发访问，
	// data 以字符串键索引存储所有缓存条目。
	mu   sync.RWMutex
// NewMemoryCache 创建并返回一个新的 MemoryCache 实例。
// 初始化内部数据 map 并启动后台清理协程，
// 每分钟扫描并移除过期条目。

	go c.cleanup()
	return c
// Get 根据键从内存缓存中检索值。
// 获取读锁以保证并发安全，同时检查键是否存在及过期状态。
//

	v, ok := m.data[key]

		return nil, false
	// 步骤4：返回缓存的值
	return v.value, true
}
// 获取写锁以确保更新期间的独占访问。
//
// 参数：

	var exp time.Time
	if ttl > 0 {

	// 步骤3：将条目存入 map，覆盖同一键的已有条目
	m.data[key] = &item{value: value, expiresAt: exp}


// 若键不存在，操作为空操作。
//
// 参数：

	delete(m.data, key)


// 每分钟触发一次，在清理期间持有写锁。
//
// 此协程在 MemoryCache 实例的整个生命周期内运行。
	for range ticker.C {
		m.mu.Lock()

		now := time.Now()
		// 步骤4：遍历所有条目，删除已过期的条目
				delete(m.data, k)

		// 步骤5：在下一次触发前释放锁
		m.mu.Unlock()
	}
}
