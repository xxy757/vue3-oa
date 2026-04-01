package main

import (
	"fmt"
	"log"

	"oa-saas/internal/config"
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/cache"
	"oa-saas/internal/router"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	var c cache.Cache
	if cfg.Redis.Enabled {
		rc, err := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			log.Printf("Redis connection failed: %v, using memory cache", err)
			c = cache.NewMemoryCache()
		} else {
			c = rc
		}
	} else {
		c = cache.NewMemoryCache()
		log.Println("Redis disabled, using memory cache")
	}

	r := router.Setup(db, c, cfg)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.Plan{},
		&model.Tenant{},
		&model.Invoice{},
		&model.TenantLog{},
		&model.User{},
		&model.Department{},
		&model.Role{},
		&model.ApprovalFlow{},
		&model.Approval{},
		&model.ApprovalNode{},
		&model.Notice{},
		&model.NoticeRead{},
		&model.Schedule{},
		&model.ScheduleParticipant{},
	)

	model.SeedData(db)

	return db, nil
}
