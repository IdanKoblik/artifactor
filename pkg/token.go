package pkg

type ApiToken struct {
	Token    string         `json:"token" bson:"_id"`
	Admin    bool           `json:"admin" bson:"admin"`
	Products []TokenProduct `json:"products" bson:"products"`
}

type TokenProduct struct {
	Name        string             `json:"name" bson:"name"`
	Permissions ProductPermissions `json:"permissions" bson:"permissions"`
}

type ProductPermissions struct {
	Download bool `json:"download" bson:"download"`
	Upload   bool `json:"upload" bson:"upload"`
	Delete   bool `json:"delete" bson:"delete"`
}
