package jira

import (
	"net/http"
	"net/url"

	cftoken "github.com/cloudflare/cloudflared/cmd/cloudflared/token"

	"golang.org/x/sync/syncmap"
)

const cookieheader = "cookie"

type CloudflareAccessRoundTripper struct {
	tokens *syncmap.Map

	t http.RoundTripper
}

func NewCloudflareAccessTransport(t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = http.DefaultTransport
	}
	return &CloudflareAccessRoundTripper{t: t, tokens: new(syncmap.Map)}
}

func (c *CloudflareAccessRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	u := url.URL{
		Host:   r.URL.Host,
		Scheme: "https",
	}
	token, ok := c.tokens.Load(u.Host)
	if !ok {
		if tok, err := cftoken.GetTokenIfExists(&u); err == nil {
			c.tokens.Store(u.Host, tok)
			token = tok
		}
		if tok, err := cftoken.FetchToken(&u); err == nil {
			c.tokens.Store(u.Host, tok)
			token = tok
		} else {
			return nil, err
		}
	}

	r.Header.Add("cf-access-token", token.(string))
	return c.t.RoundTrip(r)
}
