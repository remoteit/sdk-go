package contracts

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

const (
	InvalidApplicationType = -1
)

type ApplicationType struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Port        int    `json:"port,omitempty"`
	Proxy       bool   `json:"proxy,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
}

func NewCustomApplicationType(applicationType int) ApplicationType {
	return ApplicationType{
		ID: applicationType,
	}
}

const (
	TCPServiceID       int = 1
	BulkServiceID      int = 35
	MultiPortServiceID int = 40
)

type DefinedDevice struct {
	ID       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	Services []DefinedService `json:"services,omitempty"`
}

type DefinedService struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

const DefaultServiceType = "00:00:00:00:00:01:00:00:04:60:00:00"

func GetServiceType(deviceType string, applicationType int, manufacturer int, platform int) string {
	deviceTypeAsBytes := parseHex(deviceType)

	// buffer.writeUInt16BE(applicationType, 0)
	applicationTypeAsHex := intToHex(applicationType)
	deviceTypeAsBytes = writeBytesToBuffer(deviceTypeAsBytes, applicationTypeAsHex, 0)

	// buffer.writeUInt16BE(applicationType === 35 ? 40 : 0, 10)
	val1AsHex := intToHex(0)
	deviceTypeAsBytes = writeBytesToBuffer(deviceTypeAsBytes, val1AsHex, 10)

	if applicationType == BulkServiceID {
		multiPortServiceIDAsHex := intToHex(MultiPortServiceID)
		deviceTypeAsBytes = writeBytesToBuffer(deviceTypeAsBytes, multiPortServiceIDAsHex, 10)
	}

	// if (manufacturer) buffer.writeUInt16BE(manufacturer, 2)
	manufacturerAsHex := intToHex(manufacturer)
	deviceTypeAsBytes = writeBytesToBuffer(deviceTypeAsBytes, manufacturerAsHex, 2)

	// if (platform) buffer.writeUInt16BE(platform, 8)
	platformAsHex := intToHex(platform)
	deviceTypeAsBytes = writeBytesToBuffer(deviceTypeAsBytes, platformAsHex, 8)

	return formatHex(deviceTypeAsBytes, 2)
}

func intToHex(n int) []byte {
	return parseHex(fmt.Sprintf("%04x", n))
}

func parseHex(bytesAsString string) []byte {

	values := []byte{}
	for _, byteAsString := range strings.Split(bytesAsString, ":") {
		b, _ := hex.DecodeString(byteAsString)
		values = append(values, b...)
	}

	return values
}

func formatHex(buffer []byte, split int) string {
	if len(buffer) == 0 {
		return ""
	}

	if split == 0 {
		split = 2
	}

	bufferAsString := hex.EncodeToString(buffer)
	bufferAsString = strings.ToUpper(bufferAsString)

	re := regexp.MustCompile(fmt.Sprintf(`(\S{%d})`, split))
	return strings.Join(re.FindAllString(bufferAsString, -1), ":")
}

func writeBytesToBuffer(dst []byte, src []byte, startIndex int) []byte {
	for i := 0; i < len(src); i++ {
		dst[startIndex] = src[i]
		startIndex++
	}

	return dst
}
