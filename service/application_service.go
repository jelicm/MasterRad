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
	root := model.DataSpaceItem{Name: "Root", Path: app.ApplicationId, IsLeaf: true}
	ds := model.DataSpace{
		DataSpaceId: app.ApplicationId,
		SizeKB:      sizeKB,
		State:       model.Open,
		Root:        root.Path + "/" + root.Name,
	}

	service.CreateDataItem(&app, &root)

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
	//treba neka validacija za root name, tj ili zabraniti da bude name root ako nije root, ili neka drugačija provera
	if dsi.Name != "Root" {
		dsiParent, err := service.store.GetDataSpaceItem(dsi.Path)
		evaluateError(err)

		if dsiParent.IsLeaf {
			dsiParent.IsLeaf = false
			err = service.store.PutDataSpaceItem(dsiParent)
			evaluateError(err)
		}
	}

	err := service.store.PutDataSpaceItem(dsi)
	evaluateError(err)
}

// znamo od koje aplikacije uzimamo, a ne znamo direktno id od namespace-a
func (service *ApplicationService) CreateSoftlink(app1, app2 *model.Application) {

	app1, err := service.store.GetApp(app1.ParentNamespaceId, app1.ApplicationId)
	evaluateError(err)

	ds, err := service.store.GetDataSpace(app1.ApplicationId, app1.DataSpaceId)
	evaluateError(err)

	//ako je root list, znači nema podataka u tom dataspace-u
	root, err := service.store.GetDataSpaceItem(ds.Root)
	evaluateError(err)

	if root.IsLeaf {
		log.Fatal("There is no available data!")
	}

	if ds.State == model.Closed {
		log.Fatal("you cannot access closed data!")
	}

	softlink := model.Softlink{
		ApplicationID: app2.ApplicationId,
		DataSpaceID:   ds.DataSpaceId,
	}
	err = service.store.PutSoftlink(&softlink)
	evaluateError(err)

}

//da li mi uopste treba struktura hl?

func (service *ApplicationService) ChangeDateSpaceState(app model.Application, state model.State) {

	ds, _ := service.store.GetDataSpace(app.ApplicationId, app.DataSpaceId)
	ds.State = state

	err := service.store.PutDataSpace(app.ApplicationId, ds)
	evaluateError(err)

	fmt.Printf("DataSpace ds: %d; State %d\n", ds.SizeKB, ds.State)

	//ako se promeni na closed, treba obrisati sve softlinkove

	if state == model.Closed {
		service.store.DeleteAllSoftlinksForDataSpace(ds.DataSpaceId)
	}

}
