package types

type AuthRequest struct {
	Username string
	SsoId string
	Host string
}

type Account struct {
	DisplayName string
	AuthData    Auth
}

type Auth struct {
	Account  int
	Username string
	SsoId    int
	Host     int
}
