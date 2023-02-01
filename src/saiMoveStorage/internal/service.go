package internal

import (
	"github.com/saiset-co/saiService"
	"github.com/webmakom-com/saiMoveStorage/utils"
	"go.uber.org/zap"
)

type InternalService struct {
	Context     *saiService.Context
	Logger      *zap.Logger
	StorageFrom *utils.Database
	StorageTo   *utils.Database
}

func Service(service *saiService.Service) *InternalService {
	return &InternalService{
		Context: service.Context,
		Logger:  service.Logger,
		StorageFrom: utils.Storage(
			service.Context.GetConfig("service.storage_from.url", "").(string),
			service.Context.GetConfig("service.storage_from.token", "").(string),
		),
		StorageTo: utils.Storage(
			service.Context.GetConfig("service.storage_to.url", "").(string),
			service.Context.GetConfig("service.storage_to.token", "").(string),
		),
	}
}
