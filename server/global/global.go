package global

import (
	"crm/config"

	"github.com/go-pay/gopay/alipay"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var (
	Config config.Config
	Db     *gorm.DB
	Rdb    *redis.Client
	Alipay *alipay.Client
)

/*
全局只使用一个 MySQL 连接的问题

使用一个全局的 *gorm.DB 指针来操作数据库在某些场景下是可行的，尤其是在单例模式或服务层需要统一管理数据库连接的情况下。Gorm 本身具有连接池管理功能，它能够处理多个并发请求，而不需要每个请求都创建一个新的数据库连接。
当使用全局的 *gorm.DB 实例时，Gorm 会在内部维护一个连接池，这个连接池会管理一定数量的数据库连接，以便在需要时复用这些连接。这有助于减少数据库连接的创建和销毁带来的性能开销。
如果在使用过程中数据库连接发生中断，Gorm 通常会尝试重新建立连接。根据 Gorm 的文档，当执行查询或事务时，如果连接失败，Gorm 会尝试从连接池中获取新的连接。如果所有连接都无法使用，Gorm 会尝试创建新的连接。如果配置了重试机制，它还可能在失败后重试操作。
但是，需要注意的是，如果数据库长时间不可用，连接池中的所有连接都可能失效，这时 Gorm 将无法自动恢复，除非重启应用或者手动调用相关方法重新初始化数据库连接。
为了确保应用程序的健壮性，可以采取以下措施：
监控与健康检查：定期检查数据库连接的状态，确保连接可用。
错误处理：在关键的数据库操作处添加错误处理逻辑，比如重试机制或记录错误日志。
连接超时与重试策略：在 Gorm 配置中设置合理的连接超时时间和重试次数。
应用重启策略：在容器编排系统（如 Kubernetes）中，可以配置健康检查探针（Health Checks），当检测到应用不健康时自动重启。
使用连接状态监控工具：使用如 Prometheus 和 Grafana 这样的工具来监控数据库连接的状态，及时发现并解决问题。
总之，虽然使用全局的 *gorm.DB 指针是可行的，但必须考虑到连接中断的情况，并实施适当的预防和恢复措施。


*/
