package ktntwitchgo

type AuthEvent struct {
	AccessToken		string		`json:"access_token"`
	RefreshToken	string		`json:"refresh_token"`
	TokenType		*string		`json:"token_type"`
	ExpiresIn		int			`json:"expires_in"`
	Scope			[]string	`json:"scope"`
}
