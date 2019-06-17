package service

type ServiceRepositoryInterface interface{
	Update(service Service, serviceExists Service) (string, int)
	CreateService(s Service) (string, int) 
	ListServices(page int, filterQuery string) []Service
	DeleteService(s Service) (string, int)
	Find(s Service) (Service, error)
	ListAllAvailableHosts() ([]string, error)
	NormalizeServices() error
}