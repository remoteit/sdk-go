package contracts

import (
	"encoding/json"
)

type Protocol int

const (
	TCP Protocol = iota
	UDP
)

func GetServiceProtocol(serviceType int) Protocol {
	if (serviceType & 0x8000) != 0 {
		return UDP
	}

	return TCP
}

type Service struct {
	CreatedTimestamp int64  `json:"createdtimestamp"`
	Disabled         bool   `json:"disabled"`
	HardwareID       string `json:"hardwareid"`
	Hostname         string `json:"hostname"`
	Overload         int    `json:"overload"`
	Port             int    `json:"port"`
	Secret           string `json:"secret"`
	Type             int    `json:"type"`
	UID              string `json:"uid"`
	TemplateID       string `json:"templateID"`
	ManufactureID    int    `json:"manufactureid"`
}

func (s Service) IsMultiPort() bool {
	return s.Type == MultiPortServiceID || s.Overload == MultiPortServiceID
}

func (s Service) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(data)
}

type ServiceRegistrationInfo struct {
	Name             string
	ServiceType      string
	ServiceTypeAsInt int
	HardwareID       string
}

type DefinedService struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
