package contracts

type Auth struct {
	AuthHash string `json:"authhash"`
	Username string `json:"username"`
	Guid     string `json:"guid"`
}
