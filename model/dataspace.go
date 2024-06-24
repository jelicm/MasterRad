package model

type State int

const (
	Open State = iota
	Closed
)

type DataSpaceItem struct {
	DataSpaceItemId string
	SizeKB          int
	Children        []DataSpaceItem
	Parent          *DataSpaceItem
}

type DataSpace struct {
	SizeKB int
	State  State
	Root   DataSpaceItem
}
