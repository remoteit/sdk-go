package contracts

type ServiceCredentials struct {
	UID    string `json:"uid"`
	Secret string `json:"secret"`
}

type ServiceConfigResponse struct {
	Hostname string
	Port     int
	Type     int
	Disabled bool
}
