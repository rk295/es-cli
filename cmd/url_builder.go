package cmd

import (
	"net/url"
	"path"
)

type queryParams map[string]string

func buildURL(URLPath string, query map[string]string) string {
	u, err := url.Parse(esURL)
	if err != nil {
		// Better return something than nothing
		return esURL
	}
	q := u.Query()
	q.Set("format", "json")
	q.Set("pretty", "")
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	if URLPath != "" {
		u.Path = path.Join(u.Path, URLPath)
	}
	return u.String()
}
