package models

import (
	"fmt"
	"go-admin/common/models"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
)

type SysProject struct {
	ProjectId          int    `json:"projectId" gorm:"primaryKey;autoIncrement;comment:projectId"`
	ProjectName        string `json:"projectName" gorm:"type:varchar(100);comment:ProjectName"`
	Status             string `json:"status" gorm:"type:varchar(10);comment:status"`
	ProjectOpenstackId string `json:"projectOpenstackId" gorm:"type:varchar(100);comment:ProjectOpenstackId"`
	Tag                string `json:"tag" gorm:"type:varchar(100);comment:tag"`
	models.ModelTime
	models.ControlBy
}

func (SysProject) TableName() string {
	return "sys_project"
}

func (e *SysProject) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysProject) GetId() interface{} {
	return e.ProjectId
}

func CreateIdentityProvider(TenantName string) *gophercloud.ProviderClient {
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

func CreateIdentityClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Printf("openstack create identity client error:%s \r\n", err)
		return nil
	}
	return client
}

func CreateProject(identityClient *gophercloud.ServiceClient, project_name string, Project_Description string) (project *projects.Project, err error) {
	createOpts := projects.CreateOpts{
		Name:        project_name,
		Description: Project_Description,
	}

	project, err = projects.Create(identityClient, createOpts).Extract()
	if err != nil {
		fmt.Printf("openstack create project error:%s \r\n", err)
		return nil, err
	}

	projectID := GetProjectId(identityClient, project_name)
	userID := "68cabe91a1904b3da72959341823d3e0"
	roleID := "d8c478fc612242a284f4b3007adafb68"

	err = roles.Assign(identityClient, roleID, roles.AssignOpts{
		UserID:    userID,
		ProjectID: projectID,
	}).ExtractErr()

	if err != nil {
		fmt.Printf("openstack add admin to the project error:%s \r\n", err)
		return nil, err
	}

	return
}

func UpateProject(identityClient *gophercloud.ServiceClient, newname string, oldname string) (project *projects.Project, err error) {
	projectID := GetProjectId(identityClient, oldname)

	updateOpts := projects.UpdateOpts{
		Name: newname,
	}

	project, err = projects.Update(identityClient, projectID, updateOpts).Extract()
	if err != nil {
		fmt.Printf("openstack update project error:%s \r\n", err)
		return nil, err
	}
	return
}

func DelteProject(identityClient *gophercloud.ServiceClient, name string) (err error) {
	projectID := GetProjectId(identityClient, name)
	err = projects.Delete(identityClient, projectID).ExtractErr()
	if err != nil {
		fmt.Printf("openstack delete project error:%s \r\n", err)
		return err
	}
	return
}

func GetProjectId(identityClient *gophercloud.ServiceClient, project_name string) (projectID string) {
	//get tenantid
	listOpts := projects.ListOpts{
		Name: project_name,
	}

	allPages, err := projects.List(identityClient, listOpts).AllPages()
	if err != nil {
		fmt.Printf("openstack get project id error:%s \r\n", err)
		return ""
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		fmt.Printf("openstack get project id error:%s \r\n", err)
		return ""
	}

	projectID = allProjects[0].ID
	return
}

func GetProjectList(identityClient *gophercloud.ServiceClient, name string) (projectList []projects.Project, err error) {
	//get tenantid
	listOpts := projects.ListOpts{
		Name: name,
	}

	allPages, err := projects.List(identityClient, listOpts).AllPages()
	if err != nil {
		fmt.Printf("openstack get project id error:%s \r\n", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		fmt.Printf("openstack get project id error:%s \r\n", err)
	}

	return allProjects, err
}
