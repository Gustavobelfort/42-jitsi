package oauth

import (
	"net/url"
	"path"
)

func joinURL(src *url.URL, elem ...string) *url.URL {
	newURL := new(url.URL)
	*newURL = *src

	newURL.Path = path.Join(newURL.Path, path.Join(elem...))

	return newURL
}
