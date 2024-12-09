package initialize

import (
	"fmt"

	"crm/global"

	"github.com/go-redis/redis/v9"
)

// 初始化Redis数据库
func Redis() {
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
		Password: r.Password,
		DB:       r.Database,
	})
	global.Rdb = rdb
}

/*

全局只使用一个 redis 连接问题


### 连接池
- **默认情况下**，`redis.NewClient` 会自动创建并管理一个连接池。因此，通常情况下不需要手动设置连接池。
- **连接池的作用**：连接池可以复用连接，减少每次请求时建立和断开连接的开销，提高性能。

### 全局使用一个 Redis 客户端的风险
- **资源竞争**：多个 goroutine 同时访问同一个 Redis 客户端时，可能会导致资源竞争问题，但 `redis.Client` 是线程安全的，所以一般不会有问题。
- **单点故障**：如果全局使用的 Redis 客户端出现故障（如网络中断、Redis 服务不可用等），所有依赖该客户端的代码都会受到影响。
- **配置变更**：如果需要动态调整连接池大小或其他配置，全局使用一个客户端可能会比较困难。

### 重连机制
- **默认重连机制**：`redis.Client` 内部已经实现了重连机制。当连接断开时，客户端会尝试重新连接。
- **自定义重连逻辑**：如果需要更复杂的重连逻辑（如重试次数、间隔时间等），可以通过设置 `redis.Options` 中的 `Dialer` 来实现。

### 建议
- **使用默认连接池**：通常情况下，默认的连接池配置已经足够使用。
- **监控和日志**：建议添加监控和日志记录，以便在出现问题时能够及时发现并处理。
- **容错处理**：在关键业务逻辑中，可以添加容错处理，如超时重试、降级策略等。

总结来说，当前的代码已经基本满足了大多数场景的需求，但可以根据具体业务需求进一步优化和增强。
*/
