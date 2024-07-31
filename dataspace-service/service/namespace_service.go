package service

import (
	"fmt"
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
	fmt.Println(len(apps))
	var dataItems []string

	for _, app := range apps {

		ds, err := service.store.GetDataSpace(app.ApplicationId, app.DataSpaceId)
		evaluateError(err)
		fmt.Println(len(ds.OpenItems))
		itemsPaths, err := service.store.GetAllSchemes(ds.OpenItems)
		evaluateError(err)
		dataItems = append(dataItems, itemsPaths...)

	}

	return dataItems
}

func (service *NamespaceService) DeleteAppDefault(nsId, appId string) error {
	// default brisanje - brišemo i app i ds zajedno
	app, err := service.store.GetApp(nsId, appId)
	if err != nil {
		return err
	}
	ds, err := service.store.GetDataSpace(app.ApplicationId, app.DataSpaceId)
	if err != nil {
		return err
	}

	err = service.store.DeleteAppDefault(app)
	if err != nil {
		return err
	}

	//brisemo softlinkove za taj ds, proveravamo samo ono sto je open
	for _, item := range ds.OpenItems {
		err = service.store.DeleteAllSoftlinksForDataSpaceItem(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func evaluateError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
