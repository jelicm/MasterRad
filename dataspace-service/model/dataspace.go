package model

type State int

const (
	Open State = iota
	Closed
	Mix
)

type DataSpaceItem struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	SizeKB      int    `json:"sizeKb"`
	State       State  `json:"state"`
	Scheme      bool   `json:"scheme"`
	Permissions string `json:"permissions"`
	IsLeaf      bool   `json:"isLeaf"`
}

type DataSpace struct {
	DataSpaceId string   `json:"dataSpaceId"`
	SizeKB      int      `json:"sizeKB"`
	UsedKB      int      `json:"usedKB"`
	Root        string   `json:"root"`
	OpenItems   []string `json:"openItems"`
}
