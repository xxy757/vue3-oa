
import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 扩展标准 JWT 注册声明，增加多租户用户认证所需的
	UserID uint `json:"user_id"`
	UserID   uint `json:"user_id"`
	jwt.RegisteredClaims

// GenerateToken 创建一个新的已签名 JWT 令牌，包含给定的用户和租户信息。
// 使用 HMAC-SHA256 签名，并根据指定的有效期设置过期声明。
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// 根据配置的时长设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			// 记录令牌签发时间，用于审计
		},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
	// 步骤2：使用 HS256 签名方法创建包含声明的新 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
// ParseToken 验证并解析 JWT 令牌字符串，提取自定义 Claims 载荷。
// 使用提供的密钥验证签名，并检查令牌是否已过期。
//
	})
		return nil, err

	// 步骤2：类型断言解析后的声明并验证令牌有效性
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
