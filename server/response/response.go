package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Page struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

// 响应结果
func Result(code int, data interface{}, c *gin.Context) {
	message := msg[code]
	c.JSON(http.StatusOK, Response{code, message, data})
}

// PageResult 生成分页响应结果
// code: 响应状态码，对应预定义的消息
// data: 响应的具体数据，类型为interface{}，可灵活适应各种数据类型
// rows: 数据总行数，用于分页信息的展示
// c: gin.Context，用于向客户端返回响应
// 响应分页结果
func PageResult(code int, data interface{}, rows int64, c *gin.Context) {
	// 根据code从msg映射中获取对应的消息文本
	message := msg[code]

	// 创建Page对象，包含总行数和数据列表
	page := &Page{Total: rows, List: data}

	// 构造Response对象，包含状态码、消息和分页信息
	// 最后通过c.JSON将Response对象以JSON格式返回给客户端
	c.JSON(http.StatusOK, Response{code, message, page})
}
