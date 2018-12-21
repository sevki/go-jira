package jira

import (
	"net/http"
)

const cookieheader = "cookie"

func NewCloudflareAccessTransport(token string, t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = http.DefaultTransport
	}
	return &CloudflareAccessRoundTripper{token: token, t: t}
}

type CloudflareAccessRoundTripper struct {
	token string

	t http.RoundTripper
}

func (c *CloudflareAccessRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("cf-access-token", c.token)
	return c.t.RoundTrip(r)
}
