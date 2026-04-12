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

type ProjectRequest struct {
	Owner  int
	Host   string
	RepoID int
}

type Project struct {
	ID         int
	Host       int
	Repository int
	Owner      int
	CreatedAt  string
}

type Product struct {
	ID         int
	ExternalID int
	Name       string
	GroupName  string
	Project    int
	CreatedAt  string
}

type ProductRequest struct {
	ExternalID int
	Name       string
	GroupName  string
	Project    int
}

type ProjectWithProducts struct {
	Project  Project
	Products []Product
}
