package acloudapi

import (
	"github.com/go-resty/resty/v2"
)

type Authenticator interface {
	Authenticate(c *resty.Client, r *resty.Request) error
}

type personalAccessTokenAuthenticator struct {
	token string
}

func NewPersonalAccessTokenAuthenticator(token string) Authenticator {
	return &personalAccessTokenAuthenticator{
		token: token,
	}
}

func (m *personalAccessTokenAuthenticator) Authenticate(c *resty.Client, r *resty.Request) error {
	r.SetAuthScheme("Token")
	r.SetAuthToken(m.token)
	return nil
}
