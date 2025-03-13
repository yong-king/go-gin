package utills

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/youngking/gin-blog/pkg/setting"
)

/*
`gin.Context` 主要用于：

1. **获取请求信息**
  - URL 参数、查询参数、POST 表单数据、JSON 数据等。

2. **设置和返回响应**
  - 设置 HTTP 状态码、返回 JSON、字符串或文件。

3. **存储和共享数据**
  - 在中间件和处理函数之间共享数据（`c.Set()` / `c.Get()`）。

4. **控制请求流程**
  - 提前终止请求（`c.Abort()`）或重定向。
*/
func GetPage(c *gin.Context) int {
	result := 0
	/*
		com.StrTo 作用：

		安全转换 字符串 → 整数、浮点数、布尔值
		简化代码，避免手动 strconv.Atoi
		避免错误，可以使用 .MustInt() 直接报错
		在 GetPage(c *gin.Context) 中，它用于获取分页参数 page 并转换为整数，用于数据库查询的偏移量计算！
	*/
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		/*
			只有当 page > 0 时才计算偏移量：
			page=1 时，result = (1-1) * PageSize = 0（第一页从索引 0 开始）。
			page=2 时，result = (2-1) * PageSize = PageSize（第二页从 PageSize 开始）。
			page=3 时，result = (3-1) * PageSize = 2 * PageSize（第三页从 2*PageSize 开始）。

		*/
		result = (page - 1) * setting.AppSetting.PageSize
	}
	return result

}
