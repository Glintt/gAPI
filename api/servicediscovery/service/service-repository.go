package service

import (
	"github.com/Glintt/gAPI/api/users"
	"github.com/Glintt/gAPI/api/database"
)

type ServiceRepositoryInterface interface{
	Update(service Service, serviceExists Service) (string, int)
	CreateService(s Service) (string, int) 
	ListServices(page int, filterQuery string) []Service
	DeleteService(s Service) (string, int)
	Find(s Service) (Service, error)
	ListAllAvailableHosts() ([]string, error)
	NormalizeServices() error
}

func GetServicesRepository(user users.User) ServiceRepositoryInterface{
	if database.SD_TYPE == "mongo" {
		return &ServiceMongoRepository{
			User: user,
		}
	}
	if database.SD_TYPE == "oracle" {
		return &ServiceOracleRepository{
			User: user,
		}
	}
	return nil
}