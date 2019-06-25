package appgroups

import (
	"errors"

	"github.com/Glintt/gAPI/api/authentication"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	userModels "github.com/Glintt/gAPI/api/users/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

type ApplicationGroup struct {
	Id       bson.ObjectId `bson:"_id" json:"Id"`
	Name     string
	Services []bson.ObjectId
}

type ApplicationGroupService struct {
	User          userModels.User
	AppGroupRepos AppGroupRepository
}

// NewApplicationGroupService create application group service
func NewApplicationGroupService(c *routing.Context) (ApplicationGroupService, error) {
	user := authentication.GetAuthenticatedUser(c)
	appGroupServ := ApplicationGroupService{User: user}
	err := appGroupServ.createRepository()
	return appGroupServ, err
}

// NewApplicationGroupServiceWithUser create application group service
func NewApplicationGroupServiceWithUser(user userModels.User) (ApplicationGroupService, error) {
	appGroupServ := ApplicationGroupService{User: user}
	err := appGroupServ.createRepository()
	return appGroupServ, err
}

func releaseConnection(apgs *ApplicationGroupService) {
	apgs.AppGroupRepos.CommitTransaction()
	apgs.AppGroupRepos.Release()
}

func (apgs *ApplicationGroupService) createRepository() error {
	apgs.AppGroupRepos = NewAppGroupRepository(apgs.User)
	if apgs.AppGroupRepos == nil {
		return errors.New("Could not get application group repository")
	}
	return nil
}

// GetApplicationGroupsPermissions Get user application groups
func (apgs *ApplicationGroupService) GetApplicationGroupsPermissions(userID string) ([]ApplicationGroup, error) {
	apgs.AppGroupRepos.OpenTransaction()
	appGroups ,err := apgs.AppGroupRepos.GetApplicationGroupsForUser(userID)
	releaseConnection(apgs)
	return appGroups, err
}

// UpdateApplicationGroup update application group
func (apgs *ApplicationGroupService) UpdateApplicationGroup(appGroupID string, newGroup ApplicationGroup) error {
	apgs.AppGroupRepos.OpenTransaction()
	err := apgs.AppGroupRepos.UpdateApplicationGroup(appGroupID, newGroup)
	releaseConnection(apgs)
	return err
}

// FindServiceApplicationGroup find application group for a service
func (apgs *ApplicationGroupService) FindServiceApplicationGroup(serviceID string) ApplicationGroup {
	apgs.AppGroupRepos.OpenTransaction()
	appGroup := apgs.AppGroupRepos.FindServiceApplicationGroup(serviceID)
	releaseConnection(apgs)
	return appGroup
}

// CreateApplicationGroup create application group
func (apgs *ApplicationGroupService) CreateApplicationGroup(bodyMap ApplicationGroup) error {
	apgs.AppGroupRepos.OpenTransaction()
	err := apgs.AppGroupRepos.CreateApplicationGroup(bodyMap)
	releaseConnection(apgs)
	return err
}

// GetApplicationGroups get list of application groups
func (apgs *ApplicationGroupService) GetApplicationGroups(page int, nameFilter string) []ApplicationGroup {
	apgs.AppGroupRepos.OpenTransaction()
	appGroups := apgs.AppGroupRepos.GetApplicationGroups(page, nameFilter)
	releaseConnection(apgs)
	return appGroups
}

// GetApplicationGroupByID get application group by id
func (apgs *ApplicationGroupService) GetApplicationGroupByID(appGroupID string) (ApplicationGroup, error) {
	apgs.AppGroupRepos.OpenTransaction()
	appGroup, err := apgs.AppGroupRepos.GetApplicationGroupByID(appGroupID)
	releaseConnection(apgs)
	return appGroup, err

}

// GetServicesForApplicationGroup get application group's services
func (apgs *ApplicationGroupService) GetServicesForApplicationGroup(appGroup ApplicationGroup) ([]service.Service, error) {
	apgs.AppGroupRepos.OpenTransaction()
	services, err := apgs.AppGroupRepos.GetServicesForApplicationGroup(appGroup)
	releaseConnection(apgs)
	return services, err
}

// DeleteApplicationGroup delete applicaiton group
func (apgs *ApplicationGroupService) DeleteApplicationGroup(appGroupID string) error {
	apgs.AppGroupRepos.OpenTransaction()
	err := apgs.AppGroupRepos.DeleteApplicationGroup(appGroupID)
	releaseConnection(apgs)
	return err
}

// AddServiceToGroup add srevice to group
func (apgs *ApplicationGroupService) AddServiceToGroup(appGroupID string, serviceID string) error {
	apgs.AppGroupRepos.OpenTransaction()
	err := apgs.AppGroupRepos.AddServiceToGroup(appGroupID, serviceID)
	releaseConnection(apgs)
	return err
}

// RemoveServiceFromGroup get application group by id
func (apgs *ApplicationGroupService) RemoveServiceFromGroup(appGroupID string, serviceID string) error {
	apgs.AppGroupRepos.OpenTransaction()
	err := apgs.AppGroupRepos.RemoveServiceFromGroup(appGroupID, serviceID)
	releaseConnection(apgs)
	return err
}

// UngroupedServices get application group by id
func (apgs *ApplicationGroupService) UngroupedServices() []service.Service {
	apgs.AppGroupRepos.OpenTransaction()
	services := apgs.AppGroupRepos.UngroupedServices()
	releaseConnection(apgs)
	return services
}
