
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 表示所有 API 接口使用的标准 JSON 响应封装。
	Code int `json:"code"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
// 响应封装包含状态码 200 和提供的数据载荷。
//
}
// 响应封装包含相同的状态码和描述性错误消息。
//
// 参数：
