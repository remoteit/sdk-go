package contracts

type AuthResponse struct {
	AuthHash string `json:"service_authhash"`
	Guid     string `json:"guid"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Code     string `json:"code"`
	Token    string `json:"token"`
}
