package model

type Store interface {
	GetApp(namespaceID, appID string) (*Application, error)
	PutApp(application *Application) error
	PutDataSpace(applicationID string, dataSpace *DataSpace) error
}
