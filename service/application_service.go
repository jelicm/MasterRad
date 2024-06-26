package service

import (
	"fmt"
	"log"
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
		DataSpaceId:       applicationId,
	}

	err := service.store.PutApp(&app)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ApplicationId: %s, ParentNamespaceId: %s\n", app.ApplicationId, app.ParentNamespaceId)

	ds := model.DataSpace{
		DataSpaceId: app.ApplicationId,
		SizeKB:      sizeKB,
		State:       model.Open,
		Root:        model.DataSpaceItem{},
	}

	err = service.store.PutDataSpace(app.ApplicationId, &ds)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DataSpace ds: %d; State %d\n", ds.SizeKB, ds.State)
}

func (service *ApplicationService) CreateDataItem(app *model.Application, dsi *model.DataSpaceItem) {
	err := service.store.PutDataSpaceItem(dsi)
	if err != nil {
		log.Fatal(err)
	}

	//uspesno dodat idem, sad treba da se kreira hl
	hardlink := model.Hardlink{
		ApplicationID:   app.ApplicationId,
		DataSpaceItemID: dsi.DataSpaceItemId,
	}
	err = service.store.PutHardlink(&hardlink)
	if err != nil {
		log.Fatal(err)
	}

}
