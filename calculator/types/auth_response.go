package types

type AuthResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
