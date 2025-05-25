package deployflow

import (
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (*Router) SetupRouters(router *gin.RouterGroup) *gin.RouterGroup {
	// 获取部署列表，GET /api/v1/namespaces/{namespace}/deployflows
	router.GET("/namespaces/:namespace/deployflows", GetDeploys)
	// 获取单个部署，GET /api/v1/namespaces/{namespace}/deployflows/{name}
	router.GET("/namespaces/:namespace/deployflows/:name", GetDeploy)
	// 创建部署（需要JSON体），POST /api/v1/namespaces/{namespace}/deployflows
	router.POST("/namespaces/:namespace/deployflows", CreateDeploy)
	// 部分更新部署（需要JSON体），PUT /api/v1/namespaces/{namespace}/deployflows/{name}
	router.PATCH("/namespaces/:namespace/deployflows/:name", PatchDeploy)
	// 删除部署，DELETE /api/v1/namespaces/{namespace}/deployflows/{name}
	router.DELETE("/namespaces/:namespace/deployflows/:name", DeleteDeploy)
	// 回滚操作，POST /api/v1/namespaces/{namespace}/instances/{name}/rollbacks
	router.POST("/namespaces/:namespace/instances/:name/rollbacks", CreateRollback)
	// 重启实例，POST /api/v1/namespaces/{namespace}/instances/{name}/restarts
	router.POST("/namespaces/:namespace/instances/:name/restarts", CreateRestart)
	// 扩缩容操作，POST /api/v1/namespaces/{namespace}/instances/{name}/scales
	router.POST("/namespaces/:namespace/instances/:name/scales", CreateScale)

	return router
}
