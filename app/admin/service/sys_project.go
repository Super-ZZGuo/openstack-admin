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

type SysProject struct {
	service.Service
}

// GetPage 获取SysProject列表
func (e *SysProject) GetPage(c *dto.SysProjectGetPageReq, p *actions.DataPermission, list *[]models.SysProject, count *int64) error {
	var err error
	var data models.SysProject

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysProjectService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysProject对象
func (e *SysProject) Get(d *dto.SysProjectGetReq, p *actions.DataPermission, model *models.SysProject) error {
	var data models.SysProject

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysProject error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysProject对象
func (e *SysProject) Insert(c *dto.SysProjectInsertReq) error {
    var err error
    var data models.SysProject
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysProjectService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改SysProject对象
func (e *SysProject) Update(c *dto.SysProjectUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.SysProject{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("SysProjectService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除SysProject
func (e *SysProject) Remove(d *dto.SysProjectDeleteReq, p *actions.DataPermission) error {
	var data models.SysProject

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveSysProject error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
