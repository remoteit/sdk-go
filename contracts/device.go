package contracts

type Device struct {
	DeviceAddress string `json:"deviceaddress,omitempty"`
	DeviceType    string `json:"devicetype,omitempty"`
	DeviceAlias   string `json:"devicealias,omitempty"`
	OwnerUserName string `json:"ownerusername,omitempty"`
	Scripting     bool   `json:"scripting,omitempty"`
}

type DeviceListAllResponse struct {
	Status  string   `json:"status,omitempty"`
	Reason  string   `json:"reason"`
	Devices []Device `json:"devices,omitempty"`
}

type DefinedDevice struct {
	ID       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	Services []DefinedService `json:"services,omitempty"`
}
