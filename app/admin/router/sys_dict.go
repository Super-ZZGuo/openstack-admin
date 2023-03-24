package router

import (
	"go-admin/app/admin/apis"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerDictRouter)
}

func registerDictRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	dictApi := apis.SysDictType{}
	dataApi := apis.SysDictData{}
	dicts := v1.Group("/dict").Use(authMiddleware.MiddlewareFunc())
	{

		dicts.GET("/data", dataApi.GetPage)
		dicts.GET("/data/:dictCode", dataApi.Get)
		dicts.POST("/data", dataApi.Insert)
		dicts.PUT("/data/:dictCode", dataApi.Update)
		dicts.DELETE("/data", dataApi.Delete)

		dicts.GET("/type-option-select", dictApi.GetAll)
		dicts.GET("/type", dictApi.GetPage)
		dicts.GET("/type/:id", dictApi.Get)
		dicts.POST("/type", dictApi.Insert)
		dicts.PUT("/type/:id", dictApi.Update)
		dicts.DELETE("/type", dictApi.Delete)
	}
	opSelect := v1.Group("/dict-data").Use(authMiddleware.MiddlewareFunc())
	{
		opSelect.GET("/option-select", dataApi.GetAll)
	}
}
