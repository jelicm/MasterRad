package model

type Namsespace struct {
	ParentNamspaceId string `json:"parentNamspaceId"`
	Id               string `json:"id"`
	MaxSpaceKB       int    `json:"maxSpaceKB"`
}
