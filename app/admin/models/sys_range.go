package models

import (
	"fmt"
	"go-admin/common/models"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

type SysRange struct {
	RangeId          int    `json:"rangeId" gorm:"primaryKey;autoIncrement;comment:rangeid"`
	TenantName       string `json:"tenantName" gorm:"type:varchar(10);comment:TenantName"`
	RangeName        string `json:"rangeName" gorm:"type:varchar(255);comment:RangeName"`
	Status           string `json:"status" gorm:"type:varchar(10);comment:Status"`
	Image            string `json:"image" gorm:"type:varchar(100);comment:Image"`
	Flavor           string `json:"flavor" gorm:"type:varchar(100);comment:Flavor"`
	RangeOpenstackID string `json:"rangeOpenstackID" gorm:"type:varchar(100);comment:RangeOpenstackID"`
	models.ModelTime
	models.ControlBy
}

func (SysRange) TableName() string {
	return "sys_range"
}

func (e *SysRange) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysRange) GetId() interface{} {
	return e.RangeId
}

//create a new openstack cilent
func createProvider() *gophercloud.ProviderClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://controller:5000/v3/",
		Username:         "admin",
		Password:         "admin",
		DomainName:       "default",
		TenantID:         "64335e8f232f445f8c9d5bd4402f83df",
		TenantName:       "admin",
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("openstack create provider client error:%s \r\n", err)
	}

	return provider
}

func CreateComputeClient() *gophercloud.ServiceClient {
	provider := createProvider()
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("openstack create compute client error:%s \r\n", err)
	}
	return client
}
