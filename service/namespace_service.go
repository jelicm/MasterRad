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

	var dataItems []string

	for _, app := range apps {
		for _, dsID := range app.DataSpaceId {
			ds, err := service.store.GetDataSpace(app.ApplicationId, dsID)
			evaluateError(err)

			if ds.State == model.Open {
				itemsPaths, err := service.store.GetAllDataSpaceItemsForDataSpace(ds.DataSpaceId)
				evaluateError(err)
				dataItems = append(dataItems, itemsPaths...)
			}
		}

	}

	return dataItems
}

func (service *NamespaceService) DeleteAppDefault(nsId, appId string) {
	// default brisanje - brišemo i app i ds zajedno
	app, err := service.store.GetApp(nsId, appId)
	evaluateError(err)
	err = service.store.DeleteAppDefault(app)
	evaluateError(err)
}

func (service *NamespaceService) DeleteAppSaveDataSpace(nsId, app1Id, app2Id string) {
	// sa app1 se prebacuje na app2
	//TODO
	app1, err := service.store.GetApp(nsId, app1Id)
	evaluateError(err)
	app2, err := service.store.GetApp(nsId, app1Id)
	evaluateError(err)

	ns1, err := service.store.GetNamespace(app1.ParentNamespaceId)
	evaluateError(err)
	ns2, err := service.store.GetNamespace(app1.ParentNamespaceId)
	evaluateError(err)

	sumSpace := 0
	for _, dsId := range app1.DataSpaceId {
		ds, err := service.store.GetDataSpace(app1.ApplicationId, dsId)
		evaluateError(err)
		sumSpace += ds.SizeKB
	}

	if app1.ParentNamespaceId != app2.ParentNamespaceId && ns1.ParentNamspaceId != ns2.ParentNamspaceId {
		log.Fatal("You do not have permissions to access this data")
	}

	if app2.FreeSpaceKB > sumSpace {
		moveDataSpace(app1, app2)
		return
	}

	if ns1.FreeSpaceKB > sumSpace-app2.FreeSpaceKB {
		giveDiskToChildApp(ns2, app2)
		moveDataSpace(app1, app2)
		return
	}

	if ns1.ParentNamspaceId != ns2.ParentNamspaceId {
		log.Fatal("Parent namespace does not have enough resources!")
	}

	nsParent, err := service.store.GetNamespace(ns1.ParentNamspaceId)
	evaluateError(err)
	if nsParent.FreeSpaceKB > sumSpace-app2.FreeSpaceKB-ns2.FreeSpaceKB {
		giveDiskToChildNamespace(nsParent, ns2)
		giveDiskToChildApp(ns2, app2)
		moveDataSpace(app1, app2)
		return
	}

	if ns1.FreeSpaceKB > sumSpace-app2.FreeSpaceKB-ns2.FreeSpaceKB-nsParent.FreeSpaceKB {
		giveDiskToParentNamespace(ns1)
		giveDiskToChildNamespace(nsParent, ns2)
		giveDiskToChildApp(ns2, app2)
		moveDataSpace(app1, app2)
		return
	}

	log.Fatal("Not enough resources available")
}

func evaluateError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func moveDataSpace(app1, app2 *model.Application) {
	fmt.Printf("moving dataspace from %s to %s \n", app1.ApplicationId, app2.ApplicationId)
}

func giveDiskToChildApp(ns *model.Namsespace, app *model.Application) {
	fmt.Printf("give disk to child app: %s to %s \n", ns.Id, app.ApplicationId)
}

func giveDiskToChildNamespace(ns1, ns2 *model.Namsespace) {
	fmt.Printf("give disk to child namespace: %s to %s \n", ns1.Id, ns2.Id)
}

func giveDiskToParentNamespace(ns *model.Namsespace) {
	fmt.Printf("give disk to parent namespace: %s \n", ns.Id)
}
