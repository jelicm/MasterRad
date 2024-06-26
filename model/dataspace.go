package model

type State int

const (
	Open State = iota
	Closed
)

type DataSpaceItem struct {
	DataSpaceItemId string `json:"id"`
	Path            string `json:"path"`
	SizeKB          int    `json:"sizeKb"`
	Name            string `json:"name"`
	DataSpaceId     string `json:"appID"`
	IsFile          bool   `json:"isFile"`
}

type DataSpace struct {
	DataSpaceId string        `json:"dataSpaceId"`
	SizeKB      int           `json:"sizeKB"`
	State       State         `json:"state"`
	Root        DataSpaceItem `json:"root"`
}
