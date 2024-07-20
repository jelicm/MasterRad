package model

type Store interface {
	GetApp(namespaceID, appID string) (*Application, error)
	GetDataSpace(applicationID string, dataSpaceId string) (*DataSpace, error)
	GetDataSpaceItem(path string) (*DataSpaceItem, error)
	PutApp(application *Application) error
	PutDataSpace(applicationID string, dataSpace *DataSpace) error
	PutDataSpaceItem(dataSpaceItem *DataSpaceItem) error
	PutHardlink(hardlink *Hardlink) error
	PutSoftlink(softlink *Softlink) error
	DeleteAllSoftlinksForDataSpace(dataSpaceId string) error
	GetAllAppsForNamespace(namespaceId string) ([]Application, error)
	GetAllDataSpaceItemsForDataSpace(dataSpaceId string) ([]string, error)
	DeleteAppDefault(app *Application) error
	GetNamespace(namespaceID string) (*Namsespace, error)
	PutScheme(path string, scheme string) error
	GetAllSchemes(schemes []string) ([]string, error)
}
