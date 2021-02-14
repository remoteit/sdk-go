package contracts

type Authentication struct {
	AuthHash string `json:"authHash"`
	User     string `json:"user"`
	UserID   string `json:"userID"`
	Token    string `json:"token"`
}
