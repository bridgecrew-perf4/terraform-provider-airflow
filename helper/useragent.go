package helper

import (
	"context"
	"net/http"
)

type UserAgentProvider struct {
	userAgent string
}

func NewUserAgentProvider(userAgent string) (*UserAgentProvider, error) {
	return &UserAgentProvider{
		userAgent: userAgent,
	}, nil
}

func (uap *UserAgentProvider) Intercept(
	ctx context.Context,
	req *http.Request,
) error {
	req.Header.Set("User-Agent", uap.userAgent)
	return nil
}
