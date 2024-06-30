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

func (service *ApplicationService) RunApplication(applicationId, parentNamespaceId string, sizeKB int) (*model.Application, error) {
	app := model.Application{
		ApplicationId:     applicationId,
		ParentNamespaceId: parentNamespaceId,
		DataSpaceId:       applicationId,
	}

	err := service.store.PutApp(&app)

	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	fmt.Printf("ApplicationId: %s, ParentNamespaceId: %s\n", app.ApplicationId, app.ParentNamespaceId)

	ds := model.DataSpace{
		DataSpaceId: app.ApplicationId,
		SizeKB:      sizeKB,
		State:       model.Open,
		Root:        model.DataSpaceItem{Path: ""},
	}

	err = service.store.PutDataSpace(app.ApplicationId, &ds)
	if err != nil {
		return nil, err
	}
	fmt.Printf("DataSpace ds: %d; State %d\n", ds.SizeKB, ds.State)

	//kad se kreira ds, odmah se kreira i hl
	hardlink := model.Hardlink{
		ApplicationID: app.ApplicationId,
		DataSpaceID:   ds.DataSpaceId,
	}
	err = service.store.PutHardlink(&hardlink)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (service *ApplicationService) CreateDataItem(app *model.Application, dsi *model.DataSpaceItem) {
	err := service.store.PutDataSpaceItem(dsi)
	if err != nil {
		log.Fatal(err)
	}
}

// znamo od koje aplikacije uzimamo, a ne znamo direktno id od namespace-a
func (service *ApplicationService) CreateSoftlink(app1, app2 *model.Application) {

	app1, err := service.store.GetApp(app1.ParentNamespaceId, app1.ApplicationId)
	if err != nil {
		log.Fatal(err)
	}

	ds, err := service.store.GetDataSpace(app1.ApplicationId, app1.DataSpaceId)

	if err != nil {
		log.Fatal(err)
	}
	//ovu proveru vidi kako obeleziti root
	//if ds.Root.Path == "" {
	//	log.Fatal("There is no available data!")
	//}

	if ds.State == model.Closed {
		log.Fatal("you cannot access closed data!")
	}

	softlink := model.Softlink{
		ApplicationID: app2.ApplicationId,
		DataSpaceID:   ds.DataSpaceId,
	}
	err = service.store.PutSoftlink(&softlink)
	if err != nil {
		log.Fatal(err)
	}

}

//da li mi uopste treba struktura hl?

func (service *ApplicationService) ChangeDateSpaceState(app model.Application, state model.State) {

	ds, _ := service.store.GetDataSpace(app.ApplicationId, app.DataSpaceId)
	ds.State = state

	err := service.store.PutDataSpace(app.ApplicationId, ds)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DataSpace ds: %d; State %d\n", ds.SizeKB, ds.State)

	//ako se promeni na closed, treba obrisati sve softlinkove

	if state == model.Closed {
		service.store.DeleteAllSoftlinksForDataSpace(ds.DataSpaceId)
	}

}
