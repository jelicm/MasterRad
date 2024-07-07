package model

type State int

const (
	Open State = iota
	Closed
)

type DataSpaceItem struct {
	Path        string `json:"path"`
	SizeKB      int    `json:"sizeKb"`
	Name        string `json:"name"`
	DataSpaceId string `json:"appID"`
	IsLeaf      bool   `json:"isLeaf"`
}

type DataSpace struct {
	DataSpaceId string `json:"dataSpaceId"`
	SizeKB      int    `json:"sizeKB"`
	State       State  `json:"state"`
	Root        string `json:"root"`
}
