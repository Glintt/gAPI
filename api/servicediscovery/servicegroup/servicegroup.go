package servicegroup

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/users"
	"errors"
	routing "github.com/qiangxue/fasthttp-routing"
)


type ServiceGroup struct {
	Id           bson.ObjectId `bson:"_id" json:"Id"`
	Name         string
	IsReachable  bool
	HostsEnabled []string
	Services     []bson.ObjectId
}

type ServiceGroupService struct {
	User          users.User
	ServiceGroupRepos ServiceGroupRepository
}

// NewServiceGroupService create service group service
func NewServiceGroupService(c *routing.Context) (ServiceGroupService, error) {
	user := authentication.GetAuthenticatedUser(c)
	appGroupServ := ServiceGroupService{User: user}
	err := appGroupServ.createRepositoryAndBeginTransaction()
	return appGroupServ, err
}

// NewServiceGroupServiceWithUser create service group service
func NewServiceGroupServiceWithUser(user users.User) (ServiceGroupService, error) {
	appGroupServ := ServiceGroupService{User: user}
	err := appGroupServ.createRepositoryAndBeginTransaction()
	return appGroupServ, err
}

func releaseConnection(spgs *ServiceGroupService) {
	spgs.ServiceGroupRepos.CommitTransaction()
	spgs.ServiceGroupRepos.Release()	
}

func (spgs *ServiceGroupService) createRepositoryAndBeginTransaction() error {
	spgs.ServiceGroupRepos = NewServiceGroupRepository(spgs.User)
	if spgs.ServiceGroupRepos == nil {
		return errors.New("Could not get application group repository")
	}
	spgs.ServiceGroupRepos.OpenTransaction()
	return nil
}


func (sg *ServiceGroup) Contains(serviceId bson.ObjectId) bool {
	for _, v := range sg.Services {
		if serviceId == v {
			return true
		}
	}
	return false
}

// GetServiceGroups get list of servcie groups
func (spgs *ServiceGroupService) GetServiceGroups() ([]ServiceGroup, error) {
	serviceGroups, err := spgs.ServiceGroupRepos.GetServiceGroups()
	releaseConnection(spgs)
	return serviceGroups, err
}
// GetServiceGroupById get service group by id
func (spgs *ServiceGroupService) GetServiceGroupById(serviceGroupID string) (ServiceGroup, error) {
	serviceGroup, err := spgs.ServiceGroupRepos.GetServiceGroupById(serviceGroupID)
	releaseConnection(spgs)
	return serviceGroup, err
}
// AddServiceToGroup add servcie to an existing service group
func (spgs *ServiceGroupService) AddServiceToGroup(serviceGroupId string, serviceId string) error {
	err := spgs.ServiceGroupRepos.AddServiceToGroup(serviceGroupId, serviceId)
	releaseConnection(spgs)
	return err

}
// RemoveServiceFromGroup remove service from an existing service group
func (spgs *ServiceGroupService) RemoveServiceFromGroup(serviceGroupId string, serviceId string) error {
	err := spgs.ServiceGroupRepos.RemoveServiceFromGroup(serviceGroupId, serviceId)
	releaseConnection(spgs)
	return err
}
// CreateServiceGroup create new service group
func (spgs *ServiceGroupService) CreateServiceGroup(serviceGroup ServiceGroup) error {
	serviceGroup.Id = bson.NewObjectId()
	err := spgs.ServiceGroupRepos.CreateServiceGroup(serviceGroup)
	releaseConnection(spgs)
	return err
}
// UpdateServiceGroup update an already existing service group
func (spgs *ServiceGroupService) UpdateServiceGroup(serviceGroupId string, serviceGroup ServiceGroup) error {
	err := spgs.ServiceGroupRepos.UpdateServiceGroup(serviceGroupId, serviceGroup)
	releaseConnection(spgs)
	return err

}
// DeleteServiceGroup delete an already existing service group
func (spgs *ServiceGroupService) DeleteServiceGroup(serviceGroupId string) error {
	err := spgs.ServiceGroupRepos.DeleteServiceGroup(serviceGroupId)
	releaseConnection(spgs)
	return err

}
