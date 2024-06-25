package model

type State int

const (
	Open State = iota
	Closed
)

type DataSpaceItem struct {
	DataSpaceItemId string          `json:"dataSpaceItemId"`
	SizeKB          int             `json:"sizeKb"`
	Children        []DataSpaceItem `json:"children"`
	ParentId        string          `json:"parentId"`
}

type DataSpace struct {
	DataSpaceId string        `json:"dataSpaceId"`
	SizeKB      int           `json:"sizeKB"`
	State       State         `json:"state"`
	Root        DataSpaceItem `json:"root"`
}
