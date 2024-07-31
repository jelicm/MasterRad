package service

import (
	"fmt"
	"log"
	"projekat/model"
	"slices"
	"strings"
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

func (service *NamespaceService) ChangeDSIState(appId string, dataSpaceItemPath string, state model.State) error {
	dsi, err := service.store.GetDataSpaceItem(dataSpaceItemPath)
	if err != nil {
		return err
	}

	if dsi.State == state {
		return nil
	}

	dsiParent, err := service.store.GetDataSpaceItem(dsi.Path)
	if err != nil {
		return err
	}

	if dsiParent.State != model.Custom {
		return fmt.Errorf("cannot change state because parent is not custom")
	}

	dsId := strings.Split(dsi.Path, "/")[0]
	ds, err := service.store.GetDataSpace(appId, dsId)
	if err != nil {
		return err
	}

	if state == model.Open {
		children, err := service.store.ChangeStateForAllChildren(dsi.GetFullPath(), state)
		if err != nil {
			return err
		}
		ds.OpenItems = append(ds.OpenItems, children...)
		service.store.PutDataSpace(appId, ds)
	} else if state == model.Closed {
		_, err := service.store.ChangeStateForAllChildren(dsi.GetFullPath(), state)
		if err != nil {
			return err
		}
		var newOpenItems []string
		for _, item := range ds.OpenItems {
			if !strings.HasPrefix(item, dsi.GetFullPath()) {
				newOpenItems = append(newOpenItems, item)
			}
		}
		ds.OpenItems = newOpenItems
		service.store.PutDataSpace(appId, ds)

	} else {
		dsi.State = state
		err = service.store.PutDataSpaceItem(dsi)
		if err != nil {
			return err
		}

		for indx, item := range ds.OpenItems {
			if strings.HasPrefix(item, dsi.GetFullPath()) {
				ds.OpenItems = slices.Delete(ds.OpenItems, indx, indx)
			}
		}
	}
	return nil
}

func evaluateError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
