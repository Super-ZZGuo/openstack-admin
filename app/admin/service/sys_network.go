package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysNetwork struct {
	service.Service
}

// GetPage 获取SysNetwork列表
func (e *SysNetwork) GetPage(c *dto.SysNetworkGetPageReq, p *actions.DataPermission, list *[]models.SysNetwork, count *int64) error {
	var err error
	var data models.SysNetwork

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysNetworkService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysNetwork对象
func (e *SysNetwork) Get(d *dto.SysNetworkGetReq, p *actions.DataPermission, model *models.SysNetwork) error {
	var data models.SysNetwork

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysNetwork error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysNetwork对象
func (e *SysNetwork) Insert(c *dto.SysNetworkInsertReq) error {
	var err error
	var data models.SysNetwork
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysNetworkService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改SysNetwork对象
func (e *SysNetwork) Update(c *dto.SysNetworkPutUpdateReq, p *actions.DataPermission) error {
	var err error
	var model models.SysNetwork
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysNetwork path error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("network_id = ? ", c.NetworkId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysNetwork path error: %s", err)
		return err
	}
	return nil
}

// Remove 删除SysNetwork
func (e *SysNetwork) Remove(d *dto.SysNetworkDeleteReq, p *actions.DataPermission) error {
	var data models.SysNetwork

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysNetwork error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
