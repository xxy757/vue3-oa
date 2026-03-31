// package main 表示这是一个可执行程序的入口包
// Go 程序必须有一个 main 包，且包含 main 函数
package main

// import 导入依赖包
// Go 使用包（package）来组织代码，导入后才能使用
import (
	"fmt" // 标准库：格式化输入输出，如 Printf、Sprintf

	"log" // 标准库：简单的日志记录，如 Print、Fatal

	"oa-saas/internal/config" // 项目内部包：配置文件加载和管理

	"oa-saas/internal/model" // 项目内部包：数据库模型定义（User、Department 等）

	"oa-saas/internal/pkg/cache" // 项目内部包：缓存组件（Redis、内存缓存）

	"oa-saas/internal/router" // 项目内部包：HTTP 路由配置

	"gorm.io/driver/mysql" // GORM 的 MySQL 驱动，用于连接 MySQL 数据库

	"gorm.io/gorm" // GORM：Go 最流行的 ORM 框架，用于操作数据库
)

// main 函数是程序的唯一入口
// Go 程序启动时会自动调用 main 函数
func main() {
	// 加载配置文件
	// := 是 Go 的短变量声明，自动推断变量类型
	// config.Load() 返回两个值：配置对象和错误（Go 的多返回值特性）
	cfg, err := config.Load()
	if err != nil {
		// log.Fatalf 打印错误信息并退出程序（退出码为 1）
		// %v 是格式化占位符，可以打印任意类型的值
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	// 调用 initDB 函数（在本文件后面定义），传入配置对象
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 初始化缓存组件
	// var 声明变量但不赋值，后面会根据条件赋值
	// cache.Cache 是一个接口类型，可以是 Redis 或内存缓存
	var c cache.Cache

	// 根据配置决定使用哪种缓存
	if cfg.Redis.Enabled {
		// 配置启用了 Redis，尝试连接
		rc, err := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			// Redis 连接失败，打印警告并降级使用内存缓存
			// log.Printf 只打印日志，不会退出程序
			log.Printf("Redis connection failed: %v, using memory cache", err)
			c = cache.NewMemoryCache()
		} else {
			// Redis 连接成功，使用 Redis 缓存
			c = rc
		}
	} else {
		// 配置未启用 Redis，直接使用内存缓存
		c = cache.NewMemoryCache()
		log.Println("Redis disabled, using memory cache")
	}

	// 初始化 HTTP 路由
	// router.Setup 配置所有的 API 路由和中间件
	r := router.Setup(db, c, cfg)

	// 构建服务器监听地址
	// fmt.Sprintf 格式化字符串，%d 是整数占位符
	// 结果如 ":8080"，表示监听 8080 端口
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	// 启动 HTTP 服务器
	// r.Run(addr) 会阻塞程序，持续监听请求
	// 直到服务器停止或出错才会返回
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDB 初始化数据库连接
// 参数：cfg *config.Config - 配置对象的指针（指针可以避免复制大对象）
// 返回：*gorm.DB - 数据库连接对象指针，error - 错误信息（nil 表示成功）
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 构建数据库连接字符串（DSN = Data Source Name）
	// 格式：用户名:密码@tcp(主机:端口)/数据库名?参数
	// charset=utf8mb4：支持完整的 UTF-8 字符（包括 emoji）
	// parseTime=True：自动将数据库时间类型转为 Go 的 time.Time
	// loc=Local：使用本地时区
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// 使用 GORM 打开数据库连接
	// mysql.Open(dsn) 创建 MySQL 连接器
	// &gorm.Config{} 是 GORM 的配置选项（空配置使用默认值）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 连接失败，返回 nil 和错误信息
		return nil, err
	}

	// AutoMigrate 自动迁移数据库表结构
	// 根据 Go 结构体自动创建或更新数据库表
	// 传入结构体的指针 &model.User{} 等
	// 不需要手动写 SQL 建表语句
	db.AutoMigrate(
		&model.User{},               // 用户表
		&model.Department{},         // 部门表
		&model.Role{},               // 角色表
		&model.ApprovalFlow{},       // 审批流程表
		&model.Approval{},           // 审批记录表
		&model.ApprovalNode{},       // 审批节点表
		&model.Notice{},             // 公告表
		&model.NoticeRead{},         // 公告阅读记录表
		&model.Schedule{},           // 日程表
		&model.ScheduleParticipant{}, // 日程参与者表
	)

	// SeedData 填充初始数据
	// 如创建默认管理员账号、默认部门等
	model.SeedData(db)

	// 返回数据库连接对象，nil 表示没有错误
	return db, nil
}
