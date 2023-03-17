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

type SysFlavor struct {
	service.Service
}

// GetPage 获取SysFlavor列表
func (e *SysFlavor) GetPage(c *dto.SysFlavorGetPageReq, p *actions.DataPermission, list *[]models.SysFlavor, count *int64) error {
	var err error
	var data models.SysFlavor

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysFlavorService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysFlavor对象
func (e *SysFlavor) Get(d *dto.SysFlavorGetReq, p *actions.DataPermission, model *models.SysFlavor) error {
	var data models.SysFlavor

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysFlavor error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysFlavor对象
func (e *SysFlavor) Insert(c *dto.SysFlavorInsertReq) error {
	var err error
	var data models.SysFlavor
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysFlavorService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改SysFlavor对象
// func (e *SysFlavor) Update(c *dto.SysFlavorUpdateReq, p *actions.DataPermission) error {
//     var err error
//     var data = models.SysFlavor{}
//     e.Orm.Scopes(
//             actions.Permission(data.TableName(), p),
//         ).First(&data, c.GetId())
//     c.Generate(&data)

//     db := e.Orm.Save(&data)
//     if err = db.Error; err != nil {
//         e.Log.Errorf("SysFlavorService Save error:%s \r\n", err)
//         return err
//     }
//     if db.RowsAffected == 0 {
//         return errors.New("无权更新该数据")
//     }
//     return nil
// }

// Remove 删除SysFlavor
func (e *SysFlavor) Remove(d *dto.SysFlavorDeleteReq, p *actions.DataPermission) error {
	var data models.SysFlavor

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysFlavor error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
