package gitlab

type OauthToken struct {
	Token string `json:"access_token"`
}

type GitlabUser struct {
	ID	  	  int    `json:"id"`
	Username  string `json:"username"`
	Token 	  string `json:"-"`
	Host 	  string `json:"-"`
}
