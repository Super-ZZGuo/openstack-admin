package models

import (
	"errors"
	"fmt"
	"go-admin/common/models"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

type SysImage struct {
	ImageId     int    `json:"imageId" gorm:"primaryKey;autoIncrement;comment:imageid"`
	ImageName   string `json:"imageName" gorm:"type:varchar(255);comment:ImageName"`
	OpenstackId string `json:"openstackId" gorm:"type:varchar(100);comment:OpenstackId"`
	Type        string `json:"type" gorm:"type:varchar(100);comment:OpenstackId"`
	Tag         string `json:"tag" gorm:"type:varchar(100);comment:Tag"`
	Path        string `json:"-" gorm:"type:varchar(100);comment:Tag"`
	models.ModelTime
	models.ControlBy
}

func (SysImage) TableName() string {
	return "sys_image"
}

func (e *SysImage) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysImage) GetId() interface{} {
	return e.ImageId
}

func CreateImageProvider(TenantName string) *gophercloud.ProviderClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://controller:5000/v3/",
		Username:         "admin",
		Password:         "admin",
		DomainName:       "default",
		// TenantID:         "64335e8f232f445f8c9d5bd4402f83df",
		TenantName: TenantName,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("openstack create provider client error:%s \r\n", err)
		return nil
	}

	return provider
}

func CreateImageClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("openstack create image client error:%s \r\n", err)
		return nil
	}
	return client
}

func CreateImage(imageClient *gophercloud.ServiceClient, name string, diskformat string, path string) error {
	createOpts := images.CreateOpts{
		Name:            name,
		DiskFormat:      diskformat,
		ContainerFormat: "bare",
	}

	_, err := images.Create(imageClient, createOpts).Extract()
	if err != nil {
		fmt.Printf("openstack create image error:%s \r\n", err)
		return err
	}

	//上传镜像
	imageID := GetImageId(imageClient, name)
	if imageID == "" {
		fmt.Printf("openstack get imageid error:%s \r\n", err)
		return errors.New("get image id error")
	}

	imageData, err := os.Open(path)
	if err != nil {
		fmt.Printf("openstack upload image to openstack error:%s \r\n", err)
		return err
	}
	defer imageData.Close()

	err = imagedata.Upload(imageClient, imageID, imageData).ExtractErr()
	if err != nil {
		fmt.Printf("openstack upload image to openstack error:%s \r\n", err)
		return err
	}
	return nil
}

func UpadteImage(imageClient *gophercloud.ServiceClient, newName string, imageID string) error {
	updateOpts := images.UpdateOpts{
		images.ReplaceImageName{
			NewName: newName,
		},
	}

	_, err := images.Update(imageClient, imageID, updateOpts).Extract()
	if err != nil {
		fmt.Printf("openstack upload image to openstack error:%s \r\n", err)
		return err
	}
	return nil
}

func DeleteImage(imageClient *gophercloud.ServiceClient, imageID string) error {
	err := images.Delete(imageClient, imageID).ExtractErr()
	if err != nil {
		fmt.Printf("openstack delelte image error:%s \r\n", err)
		return err
	}
	return nil
}

func GetImageId(imagesClient *gophercloud.ServiceClient, name string) (Id string) {
	listOpts := images.ListOpts{
		Name: name,
	}

	allPages, err := images.List(imagesClient, listOpts).AllPages()
	if err != nil {
		fmt.Printf("openstack get image id error:%s \r\n", err)
		return ""
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		fmt.Printf("openstack get image id error:%s \r\n", err)
		return ""
	}
	return allImages[0].ID
}
