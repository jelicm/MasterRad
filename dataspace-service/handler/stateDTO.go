package handler

type StateDTO struct {
	ApplicationId     string `json:"applicationID"`
	DataSpaceItemPath string `json:"dsiPath"`
	State             int    `json:"state"`
	Scheme            string `json:"scheme"`
}
