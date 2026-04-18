
import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 是基于 Redis 的 Cache 接口实现。
	client *redis.Client

// NewRedisCache 通过连接到指定的 Redis 服务器创建新的 RedisCache。
// 返回前执行一次 PING 操作以验证连接可用性。
		Addr:     fmt.Sprintf("%s:%d", host, port),
		DB:       db,
	})

	// 步骤2：发送 PING 命令验证连接，超时时间为 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		return nil, err
	}

	// 步骤3：返回已初始化的 RedisCache
	return &RedisCache{client: client}, nil
// Get 根据键从 Redis 中检索值。
//
// 参数：
	// 步骤2：执行 GET 命令
	if err != nil {

	// 步骤3：返回获取到的字符串值
	return val, true
}
// 值通过 go-redis 以字符串形式存储。
//
// 参数：
	// 步骤2：使用指定的键、值和 TTL 执行 SET
}
// Redis 将其视为空操作。
//
// 参数：
	// 步骤2：执行 DEL 删除键
}
