package model

import (
	"context"

	"gorm.io/gorm"
)

type tenantCtxKey struct{}

func ContextWithTenantID(ctx context.Context, tenantID uint) context.Context {
	return context.WithValue(ctx, tenantCtxKey{}, tenantID)
}

func TenantIDFromContext(ctx context.Context) uint {
	if v, ok := ctx.Value(tenantCtxKey{}).(uint); ok {
		return v
	}
	return 0
}

func TenantScope(tenantID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}
