package apis

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type SysRange struct {
	api.Api
}

// GetPage 获取SysRange列表
// @Summary 获取SysRange列表
// @Description 获取SysRange列表
// @Tags SysRange
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysRange}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-range [get]
// @Security Bearer
func (e SysRange) GetPage(c *gin.Context) {
	req := dto.SysRangeGetPageReq{}
	s := service.SysRange{}
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
	list := make([]models.SysRange, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	client, _ := openstack.NewIdentityV3(models.CreateComputeProvider("admin"), gophercloud.EndpointOpts{})
	for i := 0; i < len(list); i++ {
		fmt.Println(list[i].ProjectName)
		pList, err := models.GetProjectList(client, list[i].ProjectName)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("openstack get project失败，\r\n失败信息 %s", err.Error()))
			return
		}
		if len(pList) == 0 {
			err = s.Remove(&dto.SysRangeDeleteReq{
				Ids:         []int{list[i].RangeId},
				ProjectName: list[i].ProjectName,
				RangeNames:  []string{list[i].RangeName},
			}, p)
			if err != nil {
				e.Error(500, err, fmt.Sprintf("获取SysNetwork失败，\r\n失败信息 %s", err.Error()))
				return
			}
		}
	}

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysRange
// @Summary 获取SysRange
// @Description 获取SysRange
// @Tags SysRange
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysRange} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-range/{id} [get]
// @Security Bearer
func (e SysRange) Get(c *gin.Context) {
	req := dto.SysRangeGetReq{}
	s := service.SysRange{}
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
	var object models.SysRange

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}
	client := models.CreateComputeClient(models.CreateComputeProvider(req.ProjectName))
	rangeOpenstackId := models.GetSserverInfo(client, object.RangeName).ID
	console := models.RemoteConsole(client, rangeOpenstackId)
	console = strings.Replace(console, "controller:6080", "192.168.10.11:6080", -1)

	object.RangeConsole = console

	e.OK(object, "查询成功")
}

// Insert 创建SysRange
// @Summary 创建SysRange
// @Description 创建SysRange
// @Tags SysRange
// @Accept application/json
// @Product application/json
// @Param data body dto.SysRangeInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-range [post]
// @Security Bearer
func (e SysRange) Insert(c *gin.Context) {
	req := dto.SysRangeInsertReq{}
	s := service.SysRange{}
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

	computeClient := models.CreateComputeClient(models.CreateComputeProvider(req.ProjectName))
	imageClient := models.CreateImageClient(models.CreateImageProvider(req.ProjectName))
	imageId := models.GetImageId(imageClient, req.Image)
	flavorClient := models.CreateFlavorClient(models.CreateFlavorProvider(req.ProjectName))
	flavorId := models.GetFlavorId(flavorClient, req.Flavor)
	networkClient := models.CreateNetworkClient(models.CreateNetworkProvider(req.ProjectName))
	var networks []servers.Network

	if flavorId == "" || imageId == "" {
		e.Error(500, err, err.Error())
		return
	}

	for _, network := range req.Network {
		networkId, _ := models.GetNetworkId(networkClient, network.NetworkName)
		new := servers.Network{
			UUID:    networkId,
			FixedIP: network.Ipadress,
		}
		networks = append(networks, new)
	}

	for _, network := range req.Network {
		req.Ipadress += network.NetworkName + ": " + network.Ipadress + "\\n"
	}

	createOpts := servers.CreateOpts{
		Name:      req.RangeName,
		ImageRef:  imageId,
		FlavorRef: flavorId,
		Networks:  networks,
	}
	_, err = servers.Create(computeClient, createOpts).Extract()
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}
	req.RangeOpenstackId = models.GetSserverInfo(computeClient, req.RangeName).ID
	req.RangeConsole = "正在创建......"
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysRange
// @Summary 修改SysRange
// @Description 修改SysRange
// @Tags SysRange
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysRangeUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-range/{id} [put]
// @Security Bearer
func (e SysRange) Update(c *gin.Context) {
	req := dto.SysRangeUpdateReq{}
	new := dto.SysRangeStatusUpadteReq{}
	s := service.SysRange{}
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

	computeClient := models.CreateComputeClient(models.CreateComputeProvider(req.ProjectName))
	serverId := models.GetSserverInfo(computeClient, req.RangeName).ID

	switch req.Option {
	case "rebuild":
		imageClient := models.CreateImageClient(models.CreateImageProvider(req.ProjectName))
		imageId := models.GetImageId(imageClient, req.Image)
		err = models.RebuildServer(computeClient, req.RangeName, serverId, imageId)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("rebuild Range失败，\r\n失败信息 %s", err.Error()))
			return
		}

	case "reboot":
		err = models.RebootServer(computeClient, serverId)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("reboot Range失败，\r\n失败信息 %s", err.Error()))
			return
		}
	case "start":
		err = models.StartServer(computeClient, serverId, "start")
		if err != nil {
			e.Error(500, err, fmt.Sprintf("start Range失败，\r\n失败信息 %s", err.Error()))
			return
		}
	case "stop":
		err = models.StartServer(computeClient, serverId, "stop")
		if err != nil {
			e.Error(500, err, fmt.Sprintf("stop Range失败，\r\n失败信息 %s", err.Error()))
			return
		}
	}
	new.RangeId = req.RangeId
	new.Status = models.GetSserverInfo(computeClient, req.RangeName).Status
	console := models.RemoteConsole(computeClient, serverId)
	console = strings.Replace(console, "controller:6080", "192.168.10.11:6080", -1)
	new.RangeConsole = console

	err = s.Update(&new, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "修改成功")
}

// Delete 删除SysRange
// @Summary 删除SysRange
// @Description 删除SysRange
// @Tags SysRange
// @Param data body dto.SysRangeDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-range [delete]
// @Security Bearer
func (e SysRange) Delete(c *gin.Context) {
	s := service.SysRange{}
	req := dto.SysRangeDeleteReq{}
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
	client := models.CreateComputeClient(models.CreateComputeProvider(req.ProjectName))
	for _, name := range req.RangeNames {
		serverID := models.GetSserverInfo(client, name).ID
		err = servers.Delete(client, serverID).ExtractErr()
		if err != nil {
			e.Logger.Error(err)
			e.Error(500, err, err.Error())
			return
		}
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
