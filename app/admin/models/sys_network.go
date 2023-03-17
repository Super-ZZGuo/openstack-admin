package models

import (
	"fmt"
	"go-admin/common/models"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

type SysNetwork struct {
	NetworkId   int    `json:"networkId" gorm:"primaryKey;autoIncrement;comment:imageid"`
	NetworkName string `json:"networkName" gorm:"type:varchar(10);comment:NetworkName"`
	Cidr        string `json:"cidr" gorm:"type:varchar(20);comment:Cidr"`
	ProjectName string `json:"peojectName" gorm:"type:varchar(20);comment:"`
	PoolStart   string `json:"poolStart" gorm:"type:varchar(20);comment:PoolStart"`
	PoolEnd     string `json:"poolEnd" gorm:"type:varchar(20);comment:PoolEnd"`
	Tag         string `json:"tag" gorm:"type:varchar(255);comment:Tag"`
	models.ModelTime
	models.ControlBy
}

func (SysNetwork) TableName() string {
	return "sys_network"
}

func (e *SysNetwork) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysNetwork) GetId() interface{} {
	return e.NetworkId
}

//create a new openstack cilent
func CreateNetworkProvider(TenantName string) *gophercloud.ProviderClient {
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

func CreateNetworkClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("openstack create network client error:%s \r\n", err)
		return nil
	}
	return client
}

func CreateNetwork(networkClient *gophercloud.ServiceClient, name string, CIDR string, poolStart string, poolEnd string) error {
	iTrue := true
	createOpts := networks.CreateOpts{
		Name:         name,
		AdminStateUp: &iTrue,
	}

	_, err := networks.Create(networkClient, createOpts).Extract()
	if err != nil {
		fmt.Printf("openstack create network error:%s \r\n", err)
		return err
	}

	err = createSubNetwork(networkClient, name, CIDR, poolStart, poolEnd)
	if err != nil {
		fmt.Printf("openstack create network error:%s \r\n", err)
		return err
	}
	return nil
}

func GetNetworkId(networkClient *gophercloud.ServiceClient, name string) (string, error) {
	listOpts := networks.ListOpts{
		Name: name,
	}

	allPages, err := networks.List(networkClient, listOpts).AllPages()
	if err != nil {
		fmt.Printf("openstack create network error:%s \r\n", err)
		return "", err
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		fmt.Printf("openstack create network error:%s \r\n", err)
		return "", err
	}

	return allNetworks[0].ID, nil
}

func DeleteNetwork(networkClient *gophercloud.ServiceClient, name string) error {
	networkID, err := GetNetworkId(networkClient, name)
	if err != nil {
		fmt.Printf("openstack get network id error:%s \r\n", err)
		return err
	}
	err = networks.Delete(networkClient, networkID).ExtractErr()
	if err != nil {
		fmt.Printf("openstack delete network error:%s \r\n", err)
		return err
	}
	return nil
}

func UpadteNetwork(networkClient *gophercloud.ServiceClient, newname string, oldname string) error {
	networkID, err := GetNetworkId(networkClient, oldname)
	if err != nil {
		fmt.Printf("openstack get network id error:%s \r\n", err)
		return err
	}
	updateOpts := networks.UpdateOpts{
		Name: &newname,
	}

	_, err = networks.Update(networkClient, networkID, updateOpts).Extract()
	if err != nil {
		fmt.Printf("openstack update network id error:%s \r\n", err)
		return err
	}
	return nil
}

func createSubNetwork(networkClient *gophercloud.ServiceClient, name string, CIDR string, poolStart string, poolEnd string) error {
	NetworkID, err := GetNetworkId(networkClient, name)
	if err != nil {
		return err
	}
	createOpts := subnets.CreateOpts{
		Name:      name + "_subnet",
		NetworkID: NetworkID,
		IPVersion: 4,
		CIDR:      CIDR,
		AllocationPools: []subnets.AllocationPool{
			{
				Start: poolStart,
				End:   poolEnd,
			},
		},
	}

	_, err = subnets.Create(networkClient, createOpts).Extract()
	if err != nil {
		fmt.Printf("openstack create sub network error:%s \r\n", err)
		return err
	}
	return nil
}
