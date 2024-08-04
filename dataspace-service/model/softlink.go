package model

type SLType int

const (
	Group SLType = iota
	Others
)

type Softlink struct {
	SoftlinkID          string `json:"softlinkID"`
	ApplicationID       string `json:"appID"`
	DataSpaceItemPath   string `json:"path"`
	StoredProcedurePath string `json:"storedProcedurePath"`
	JsonParameters      string `json:"jsonParameters"`
	Type                SLType `json:"slType"`
}
