package model

type Softlink struct {
	SoftlinkID          string `json:"softlinkID"`
	ApplicationID       string `json:"appID"`
	DataSpaceItemPath   string `json:"path"`
	StoredProcedurePath string `json:"storedProcedurePath"`
	JsonParameters      string `json:"jsonParameters"`
}
