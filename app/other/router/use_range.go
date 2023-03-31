package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/other/apis"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerTest)
}

// registerTest
func registerTest(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.UseRange{}
	r := v1.Group("/use-range").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.GET("/:id", actions.PermissionAction(), api.Get)
	}
}
