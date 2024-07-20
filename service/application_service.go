package service

import (
	"fmt"
	"log"
	"projekat/model"
	"strings"

	"github.com/nats-io/nats.go"
)

type ApplicationService struct {
	store model.Store
	conn  *nats.Conn
}

func NewApplicationService(store model.Store, conn *nats.Conn) *ApplicationService {
	return &ApplicationService{
		store: store,
		conn:  conn,
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

	service.CreateDataItem(app.ApplicationId, &root, "", true)

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

func (service *ApplicationService) CreateDataItem(appID string, dsi *model.DataSpaceItem, scheme string, root bool) (*model.DataSpaceItem, error) {
	//treba neka validacija za root name, tj ili zabraniti da bude name root ako nije root, ili neka drugačija provera

	if !root {
		ds, err := service.store.GetDataSpace(appID, strings.Split(dsi.Path, "/")[0])
		if err != nil {
			return nil, err
		}

		dsiParent, err := service.store.GetDataSpaceItem(dsi.Path)
		if err != nil {
			return nil, err
		}

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
		err = service.store.PutDataSpace(appID, ds)
		if err != nil {
			return nil, err
		}

		if dsiParent.IsLeaf {
			dsiParent.IsLeaf = false
			err = service.store.PutDataSpaceItem(dsiParent)
			if err != nil {
				return nil, err
			}
		}
	}
	//transakcija?

	if dsi.Scheme {
		err := service.store.PutScheme(dsi.Path+"/"+dsi.Name, scheme)
		if err != nil {
			return nil, err
		}
	}
	err := service.store.PutDataSpaceItem(dsi)
	if err != nil {
		return nil, err
	}

	return dsi, nil

}

// znamo od koje aplikacije uzimamo, a ne znamo direktno id od namespace-a
func (service *ApplicationService) CreateSoftlink(app1, app2 *model.Application, dataSpaceItemPath string) (*string, error) {
	dsi, err := service.store.GetDataSpaceItem(dataSpaceItemPath)
	if err != nil {
		return nil, err
	}
	if app1.ParentNamespaceId != app2.ParentNamespaceId && dsi.Permissions[8] != 's' {
		//nije isti rns i others nema prava pristupa
		return nil, fmt.Errorf("no privilages for this data - others")
	}

	if app1.ParentNamespaceId == app2.ParentNamespaceId && dsi.Permissions[5] != 's' {
		//group nema dobre privilegije
		return nil, fmt.Errorf("no privilages for this data - group")
	}

	softlink := model.Softlink{
		SoftlinkID:        app2.ApplicationId + "+" + dataSpaceItemPath,
		ApplicationID:     app2.ApplicationId,
		DataSpaceItemPath: dataSpaceItemPath,
	}
	err = service.store.PutSoftlink(&softlink)
	if err != nil {
		return nil, err
	}

	_, err = service.conn.QueueSubscribe(softlink.SoftlinkID, "softlinks", func(message *nats.Msg) {
		fmt.Printf("RECEIVED MESSAGE: %s\n", string(message.Data))
	})
	if err != nil {
		return nil, err
	}
	return &softlink.SoftlinkID, nil
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
