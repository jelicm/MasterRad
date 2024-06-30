package model

type Store interface {
	GetApp(namespaceID, appID string) (*Application, error)
	GetDataSpace(applicationID string, dataSpaceId string) (*DataSpace, error)
	PutApp(application *Application) error
	PutDataSpace(applicationID string, dataSpace *DataSpace) error
	PutDataSpaceItem(dataSpaceItem *DataSpaceItem) error
	PutHardlink(hardlink *Hardlink) error
	PutSoftlink(softlink *Softlink) error
	DeleteAllSoftlinksForDataSpace(dataSpaceId string) error
}
