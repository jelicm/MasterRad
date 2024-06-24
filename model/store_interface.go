package model

type Store interface {
	GetApp(id string) (*Application, error)
	PutApp(application *Application) error
	PutDataSpace(dataSPace *DataSpace) error
}
