package service

import (
	"fmt"
	"projekat/model"
)

type ApplicationService struct {
	store model.Store
}

func NewApplicationService(store model.Store) *ApplicationService {
	return &ApplicationService{
		store: store,
	}
}

func (service *ApplicationService) RunApplication(applicationId, parentNamespaceId string, sizeKB int) {
	app := model.Application{
		ApplicationId:     applicationId,
		ParentNamespaceId: parentNamespaceId,
	}

	service.store.PutApp(&app)
	fmt.Printf("ApplicationId: %s, ParentNamespaceId: %s\n", app.ApplicationId, app.ParentNamespaceId)

	ds := model.DataSpace{
		DataSpaceId: app.ApplicationId,
		SizeKB:      sizeKB,
		State:       model.Open,
		Root:        model.DataSpaceItem{},
	}

	service.store.PutDataSpace(app.ApplicationId, &ds)
	fmt.Printf("DataSpace ds: %d; State %d\n", ds.SizeKB, ds.State)
}
