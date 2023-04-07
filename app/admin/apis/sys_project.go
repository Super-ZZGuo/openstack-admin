package apis

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type SysProject struct {
	api.Api
}

// GetPage 获取SysProject列表
// @Summary 获取SysProject列表
// @Description 获取SysProject列表
// @Tags SysProject
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysProject}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-project [get]
// @Security Bearer
func (e SysProject) GetPage(c *gin.Context) {
	req := dto.SysProjectGetPageReq{}
	s := service.SysProject{}
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
	list := make([]models.SysProject, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysProject
// @Summary 获取SysProject
// @Description 获取SysProject
// @Tags SysProject
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysProject} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-project/{id} [get]
// @Security Bearer
func (e SysProject) Get(c *gin.Context) {
	req := dto.SysProjectGetReq{}
	s := service.SysProject{}
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
	var object models.SysProject

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SysProject
// @Summary 创建SysProject
// @Description 创建SysProject
// @Tags SysProject
// @Accept application/json
// @Product application/json
// @Param data body dto.SysProjectInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-project [post]
// @Security Bearer
func (e SysProject) Insert(c *gin.Context) {
	req := dto.SysProjectInsertReq{}
	s := service.SysProject{}
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

	checkReq := dto.SysProjectGetPageReq{
		Status: "2",
	}

	p := actions.GetPermissionFromContext(c)
	openList := make([]models.SysProject, 0)
	var count int64

	err = s.GetPage(&checkReq, p, &openList, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}
	if len(openList) > 0 {
		e.Error(500, errors.New("请先关闭开启的项目再进行创建"), fmt.Sprintf("请先关闭开启的项目再进行创建"))
		return
	}

	client := models.CreateIdentityClient(models.CreateIdentityProvider("admin"))
	_, err = models.CreateProject(client, req.ProjectName, req.Tag)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ProjectOpenstackId = models.GetProjectId(client, req.ProjectName)
	req.Status = "2"

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysProject
// @Summary 修改SysProject
// @Description 修改SysProject
// @Tags SysProject
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysProjectUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-project/{id} [put]
// @Security Bearer
func (e SysProject) Update(c *gin.Context) {
	req := dto.SysProjectUpdateReq{}
	s := service.SysProject{}
	new := dto.SysProjectPutUpdate{}
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
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	new.ProjectId = req.ProjectId

	switch req.Option {
	case "cName":
		client := models.CreateIdentityClient(models.CreateIdentityProvider(req.OldProjectName))
		new.ProjectOpenstackId = models.GetProjectId(client, req.OldProjectName)
		_, err = models.UpateProject(client, req.NewProjectName, req.OldProjectName)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("修改SysProject失败，\r\n失败信息 %s", err.Error()))
			return
		}
		new.ProjectName = req.NewProjectName
		new.Tag = req.OldTag
		new.Status = req.Status
	case "cTag":
		new.ProjectName = req.OldProjectName
		new.Tag = req.NewTag
		new.Status = req.Status
	case "cTagAndName":
		client := models.CreateIdentityClient(models.CreateIdentityProvider(req.OldProjectName))
		new.ProjectOpenstackId = models.GetProjectId(client, req.OldProjectName)
		_, err = models.UpateProject(client, req.NewProjectName, req.OldProjectName)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("修改SysProject失败，\r\n失败信息 %s", err.Error()))
			return
		}
		new.ProjectName = req.NewProjectName
		new.Tag = req.NewTag
		new.Status = req.Status
	case "cStatusToOpen":
		if req.Status == "2" { // 2为开启状态 1为关闭状态
			e.OK(req.GetId(), "靶场已经开启，无法重复开启")
			return
		}
		pReq := dto.SysProjectGetPageReq{
			Status: "2",
		}
		pList := make([]models.SysProject, 0)
		var countP int64
		err = s.GetPage(&pReq, p, &pList, &countP)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
			return
		}
		if countP > 0 {
			e.OK(req.GetId(), "靶场开启失败，最多允许同时开启一个靶场")
			return
		}
		rReq := dto.SysRangeGetPageReq{
			ProjectName: req.OldProjectName,
		}
		rList := make([]models.SysRange, 0)
		var countR int64
		err = s.GetRangePage(&rReq, p, &rList, &countR)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
			return
		}
		computeClient := models.CreateComputeClient(models.CreateComputeProvider(req.OldProjectName))
		for _, pRange := range rList {
			rangeId := models.GetSserverInfo(computeClient, pRange.RangeName).ID
			err = models.StartServer(computeClient, rangeId, "start")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("开启靶机失败，\r\n失败信息 %s", err.Error()))
				return
			}
		}
		new.ProjectName = req.OldProjectName
		new.Tag = req.OldTag
		new.Status = "2"
	case "cStatusToClose":
		if req.Status == "1" { // 2为开启状态 1为关闭状态
			e.OK(req.GetId(), "靶场已经关闭，无法重复关闭")
			return
		}
		rReq := dto.SysRangeGetPageReq{
			ProjectName: req.OldProjectName,
		}
		rList := make([]models.SysRange, 0)
		var count int64
		err = s.GetRangePage(&rReq, p, &rList, &count)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
			return
		}
		computeClient := models.CreateComputeClient(models.CreateComputeProvider(req.OldProjectName))
		for _, pRange := range rList {
			rangeId := models.GetSserverInfo(computeClient, pRange.RangeName).ID
			err = models.StartServer(computeClient, rangeId, "stop")
			if err != nil {
				e.Error(500, err, fmt.Sprintf("关闭靶机失败，\r\n失败信息 %s", err.Error()))
				return
			}
		}
		new.ProjectName = req.OldProjectName
		new.Tag = req.OldTag
		new.Status = "1"
	default:
		new.ProjectName = req.OldProjectName
		new.Tag = req.OldTag
		new.Status = req.Status
	}

	err = s.Update(&new, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "修改成功")
}

// Delete 删除SysProject
// @Summary 删除SysProject
// @Description 删除SysProject
// @Tags SysProject
// @Param data body dto.SysProjectDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-project [delete]
// @Security Bearer
func (e SysProject) Delete(c *gin.Context) {
	s := service.SysProject{}
	req := dto.SysProjectDeleteReq{}
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

	for _, name := range req.ProjectNames {
		client := models.CreateIdentityClient(models.CreateIdentityProvider(name))
		projectId := models.GetProjectId(client, name)
		computeClient := models.CreateComputeClient(models.CreateComputeProvider(name))
		networkClient := models.CreateNetworkClient(models.CreateNetworkProvider(name))
		serverList := models.GetSserverList(computeClient)
		for _, server := range serverList {
			err = servers.Delete(computeClient, server.ID).ExtractErr()
			if err != nil {
				e.Logger.Error(err)
				e.Error(500, err, err.Error())
				return
			}
		}
		networkList, err := models.GetNetworkList(networkClient, projectId)
		flag := true
		for flag { //等待openstack删除结束，因为删除实例需要时间，如果直接执行删除网络的话会显示端口占用
			slist := models.GetSserverList(computeClient)
			if len(slist) == 0 {
				flag = false
			}
		}
		if err != nil {
			e.Error(500, err, fmt.Sprintf("删除SysProject失败，\r\n失败信息 %s", err.Error()))
			return
		}
		for _, network := range networkList {
			err = models.DeleteNetwork(networkClient, network.Name)
			if err != nil {
				e.Logger.Error(err)
				e.Error(500, err, err.Error())
				return
			}
		}
		err = models.DelteProject(client, name)
		if err != nil {
			return
		}
	}

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
