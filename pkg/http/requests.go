package http

import (
	"artifactor/pkg"
)

type RegisterRequest struct {
	Admin    bool               `json:"admin"`
	Products []pkg.TokenProduct `json:"products"`
}
