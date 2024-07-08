package service

import (
	"log"
	"projekat/model"
)

type NamespaceService struct {
	store model.Store
}

func NewNamespaceService(store model.Store) *NamespaceService {
	return &NamespaceService{
		store: store,
	}
}

func (service *NamespaceService) RunDataDiscovery(namespaceId string) []string {
	apps, err := service.store.GetAllAppsForNamespace(namespaceId)
	evaluateError(err)

	var dataItems []string
	for _, app := range apps {
		ds, err := service.store.GetDataSpace(app.ApplicationId, app.DataSpaceId)
		evaluateError(err)
		if ds.State == model.Open {
			itemsPaths, err := service.store.GetAllDataSpaceItemsForDataSpace(ds.DataSpaceId)
			evaluateError(err)
			dataItems = append(dataItems, itemsPaths...)
		}
	}
	return dataItems
}

func evaluateError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
