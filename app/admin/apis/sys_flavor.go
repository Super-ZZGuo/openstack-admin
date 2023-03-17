package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type SysFlavor struct {
	api.Api
}

// GetPage 获取SysFlavor列表
// @Summary 获取SysFlavor列表
// @Description 获取SysFlavor列表
// @Tags SysFlavor
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysFlavor}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-flavor [get]
// @Security Bearer
func (e SysFlavor) GetPage(c *gin.Context) {
	req := dto.SysFlavorGetPageReq{}
	s := service.SysFlavor{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.SysFlavor, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysFlavor失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysFlavor
// @Summary 获取SysFlavor
// @Description 获取SysFlavor
// @Tags SysFlavor
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysFlavor} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-flavor/{id} [get]
// @Security Bearer
func (e SysFlavor) Get(c *gin.Context) {
	req := dto.SysFlavorGetReq{}
	s := service.SysFlavor{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.SysFlavor

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysFlavor失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SysFlavor
// @Summary 创建SysFlavor
// @Description 创建SysFlavor
// @Tags SysFlavor
// @Accept application/json
// @Product application/json
// @Param data body dto.SysFlavorInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-flavor [post]
// @Security Bearer
func (e SysFlavor) Insert(c *gin.Context) {
	req := dto.SysFlavorInsertReq{}
	s := service.SysFlavor{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	client := models.CreateFlavorClient(models.CreateFlavorProvider("admin"))
	err = models.CreateFlavor(client, req.FlavorId, req.FlavorName, req.Disk, req.Ram, req.Vcpu)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysFlavor失败，\r\n失败信息 %s", err.Error()))
		return
	}

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysFlavor失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysFlavor
// @Summary 修改SysFlavor
// @Description 修改SysFlavor
// @Tags SysFlavor
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysFlavorUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-flavor/{id} [put]
// @Security Bearer
// func (e SysFlavor) Update(c *gin.Context) {
//     req := dto.SysFlavorUpdateReq{}
//     s := service.SysFlavor{}
//     err := e.MakeContext(c).
//         MakeOrm().
//         Bind(&req).
//         MakeService(&s.Service).
//         Errors
//     if err != nil {
//         e.Logger.Error(err)
//         e.Error(500, err, err.Error())
//         return
//     }
// 	req.SetUpdateBy(user.GetUserId(c))
// 	p := actions.GetPermissionFromContext(c)

// 	err = s.Update(&req, p)
// 	if err != nil {
// 		e.Error(500, err, fmt.Sprintf("修改SysFlavor失败，\r\n失败信息 %s", err.Error()))
//         return
// 	}
// 	e.OK( req.GetId(), "修改成功")
// }

// Delete 删除SysFlavor
// @Summary 删除SysFlavor
// @Description 删除SysFlavor
// @Tags SysFlavor
// @Param data body dto.SysFlavorDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-flavor [delete]
// @Security Bearer
func (e SysFlavor) Delete(c *gin.Context) {
	s := service.SysFlavor{}
	req := dto.SysFlavorDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	client := models.CreateFlavorClient(models.CreateFlavorProvider("admin"))
	for _, id := range req.Ids {
		err = models.DeleteFlavor(client, id)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("删除SysFlavor失败，\r\n失败信息 %s", err.Error()))
			return
		}
	}

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysFlavor失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
