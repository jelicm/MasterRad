package handler

type DataItemDTO struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	SizeKB      int    `json:"sizeKb"`
	State       int    `json:"state"`
	HasScheme   bool   `json:"hasSchema"`
	Permissions string `json:"permissions"`
	AppID       string `json:"appID"`
	Scheme      string `json:"scheme"`
}
