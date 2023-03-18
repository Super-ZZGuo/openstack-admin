package models

import (
	"fmt"
	"go-admin/common/models"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

type SysFlavor struct {
	FlavorId   int    `json:"flavorId" gorm:"primaryKey;autoIncrement;comment:imageid"`
	FlavorName string `json:"flavorName" gorm:"type:varchar(10);comment:FlavorName"`
	Disk       int    `json:"disk" gorm:"type:bigint(20);comment:Disk"`
	Vcpu       int    `json:"vcpu" gorm:"type:bigint(20);comment:Vcpu"`
	Ram        int    `json:"ram" gorm:"type:bigint(20);comment:Ram"`
	Tag        string `json:"tag" gorm:"type:varchar(255);comment:Tag"`
	models.ModelTime
	models.ControlBy
}

func (SysFlavor) TableName() string {
	return "sys_flavor"
}

func (e *SysFlavor) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysFlavor) GetId() interface{} {
	return e.FlavorId
}

//create a new openstack cilent
func CreateFlavorProvider(TenantName string) *gophercloud.ProviderClient {
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

func CreateFlavorClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("openstack create flavor client error:%s \r\n", err)
		return nil
	}
	return client
}

func CreateFlavor(computeClient *gophercloud.ServiceClient, ID int, Name string, Disk int, RAM int, VCPUs int) error {
	createOpts := flavors.CreateOpts{
		ID:         strconv.Itoa(ID),
		Name:       Name,
		Disk:       gophercloud.IntToPointer(Disk),
		RAM:        RAM,
		VCPUs:      VCPUs,
		RxTxFactor: 1.0,
	}

	_, err := flavors.Create(computeClient, createOpts).Extract()
	if err != nil {
		fmt.Printf("openstack create flavor error:%s \r\n", err)
		return err
	}
	return nil
}

func DeleteFlavor(computeClient *gophercloud.ServiceClient, ID int) error {
	err := flavors.Delete(computeClient, strconv.Itoa(ID)).ExtractErr()
	if err != nil {
		fmt.Printf("openstack delete flavor error:%s \r\n", err)
		return err
	}
	return nil
}

func GetFlavorId(computeClient *gophercloud.ServiceClient, name string) string {
	listOpts := flavors.ListOpts{}

	allPages, err := flavors.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		fmt.Printf("openstack get flavor id error:%s \r\n", err)
		return ""
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		fmt.Printf("openstack get flavor id error:%s \r\n", err)
		return ""
	}
	for _, flavor := range allFlavors {
		if flavor.Name == name {
			return flavor.ID
		}
	}
	return ""
}
