package api

//api meta

type API_META_VERSION struct {
	Meta_version int `json:"meta_version"`
}

//assign version
func (apim *API_META_VERSION) MetaVersion(version int) {
	apim.Meta_version = version
}

type API_META_STATUS struct {
	Meta_status  int    `json:"meta_status"`
	Meta_message string `json:"meta_message"`
}

//assign status
func (apim *API_META_STATUS) MetaStatus(status int, message string) {
	apim.Meta_message = message
	apim.Meta_status = status
}
