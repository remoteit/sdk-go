package contracts

type CertificateRequest struct {
	MachineID string `json:"machineId"`
	ServiceID string `json:"serviceId"`
	Name      string `json:"name,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type CertificateResponse struct {
	CN          string `json:"cn"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}
