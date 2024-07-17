package service

import (
	"fmt"
	"log"
	"projekat/model"
	"strings"
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
		FreeSpaceKB:       sizeKB / 2,
	}

	err := service.store.PutApp(&app)

	if err != nil {
		return nil, err
	}

	fmt.Printf("ApplicationId: %s, ParentNamespaceId: %s\n", app.ApplicationId, app.ParentNamespaceId)
	root := model.DataSpaceItem{Name: "Root", Path: app.ApplicationId, SizeKB: 1, IsLeaf: true, State: model.Mix, Scheme: false}
	ds := model.DataSpace{
		DataSpaceId: app.ApplicationId,
		SizeKB:      sizeKB / 2,
		UsedKB:      0,
		Root:        root.Path + "/" + root.Name,
		OpenItems:   []string{},
	}

	service.CreateDataItem(&app, &root, "", true)

	err = service.store.PutDataSpace(app.ApplicationId, &ds)
	if err != nil {
		return nil, err
	}
	fmt.Printf("DataSpace ds: %d;\n", ds.SizeKB)

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

func (service *ApplicationService) CreateDataItem(app *model.Application, dsi *model.DataSpaceItem, scheme string, root bool) {
	//treba neka validacija za root name, tj ili zabraniti da bude name root ako nije root, ili neka drugačija provera

	if !root {
		ds, err := service.store.GetDataSpace(app.ApplicationId, strings.Split(dsi.Path, "/")[0])
		evaluateError(err)

		dsiParent, err := service.store.GetDataSpaceItem(dsi.Path)
		evaluateError(err)

		//ignorisemo state ako je poslao korisnik jer roditelj ima vece privilegije
		if dsiParent.State != model.Mix {
			dsi.State = dsiParent.State
		}

		//ovo videti da li je validno i da li treba neka validacija da vrati 400 ako nema seme za open
		dsi.Scheme = scheme != ""

		if dsi.State == model.Open && dsi.Scheme {
			ds.OpenItems = append(ds.OpenItems, dsi.Path+"/"+dsi.Name)
		}

		if ds.UsedKB+dsi.SizeKB > ds.SizeKB {
			log.Fatal("cannot add dataitem - no available resources")
		}

		ds.UsedKB += dsi.SizeKB
		err = service.store.PutDataSpace(app.ApplicationId, ds)
		evaluateError(err)

		if dsiParent.IsLeaf {
			dsiParent.IsLeaf = false
			err = service.store.PutDataSpaceItem(dsiParent)
			evaluateError(err)
		}
	}
	//transakcija?

	if dsi.Scheme {
		err := service.store.PutScheme(dsi.Path+"/"+dsi.Name, scheme)
		evaluateError(err)
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

	softlink := model.Softlink{
		ApplicationID: app2.ApplicationId,
		DataSpaceID:   ds.DataSpaceId,
	}
	err = service.store.PutSoftlink(&softlink)
	evaluateError(err)

}

//da li mi uopste treba struktura hl?

/*func (service *ApplicationService) ChangeDateSpaceState(app model.Application, state model.State) {

	for _, dsID := range app.DataSpaceId {
		ds, _ := service.store.GetDataSpace(app.ApplicationId, dsID)
		ds.State = state

		err := service.store.PutDataSpace(app.ApplicationId, ds)
		evaluateError(err)

		fmt.Printf("DataSpace ds: %d; State %d\n", ds.SizeKB, ds.State)

		//ako se promeni na closed, treba obrisati sve softlinkove

		if state == model.Closed {
			service.store.DeleteAllSoftlinksForDataSpace(ds.DataSpaceId)
		}
	}

}*/
