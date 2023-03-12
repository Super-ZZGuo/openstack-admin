package models

import (
	"fmt"
	"go-admin/common/models"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/remoteconsoles"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type SysRange struct {
	RangeId          int         `json:"rangeId" gorm:"primaryKey;autoIncrement;comment:rangeid"`
	TenantName       string      `json:"tenantName" gorm:"type:varchar(10);comment:TenantName"`
	RangeName        string      `json:"rangeName" gorm:"type:varchar(255);comment:RangeName"`
	Status           string      `json:"status" gorm:"type:varchar(10);comment:Status"`
	Image            string      `json:"image" gorm:"type:varchar(100);comment:Image"`
	Flavor           string      `json:"flavor" gorm:"type:varchar(100);comment:Flavor"`
	RangeOpenstackId string      `json:"rangeOpenstackId" gorm:"type:varchar(100);comment:RangeOpenstackID"`
	ProjectId        int         `json:"projectId" gorm:"type:bigint(20);comment:ProjectId"`
	ProjectName      string      `json:"projectName" gorm:"type:varchar(100);comment:ProjectName"`
	RangeConsole     string      `json:"rangeConsole" gorm:"-"`
	Project          *SysProject `json:"project"`
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
func CreateComputeProvider(TenantName string) *gophercloud.ProviderClient {
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

func CreateComputeClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("openstack create compute client error:%s \r\n", err)
		return nil
	}
	return client
}

func RemoteConsole(computeClient *gophercloud.ServiceClient, serverID string) string {
	computeClient.Microversion = "2.6"
	createOpts := remoteconsoles.CreateOpts{
		Protocol: remoteconsoles.ConsoleProtocolVNC,
		Type:     remoteconsoles.ConsoleTypeNoVNC,
	}

	remtoteConsole, err := remoteconsoles.Create(computeClient, serverID, createOpts).Extract()
	if err != nil {
		fmt.Printf("openstack get remote console error:%s \r\n", err)
		return ""
	}

	return remtoteConsole.URL
}

func ServerList(client *gophercloud.ServiceClient, name string) []servers.Server {
	opts := servers.ListOpts{
		Name: name,
	}

	allPage, err := servers.List(client, opts).AllPages()
	if err != nil {
		fmt.Printf("openstack get server list error:%s \r\n", err)
		return nil
	}

	allServes, err := servers.ExtractServers(allPage)
	if err != nil {
		fmt.Printf("openstack get server list error:%s \r\n", err)
		return nil
	}
	return allServes
}

func UpateServer(client *gophercloud.ServiceClient, name string, serverID string, ImageRef string) error {
	rebuildOpts := servers.RebuildOpts{
		Name:     name,
		ImageRef: ImageRef,
	}

	_, err := servers.Rebuild(client, serverID, rebuildOpts).Extract()
	if err != nil {
		fmt.Printf("openstack update server error:%s \r\n", err)
		return err
	}
	return nil
}
