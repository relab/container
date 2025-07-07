package build

import (
	"net/url"
)

// ImageBuildOptions holds the information necessary to build images.
//
// This is a simplified version of the Docker API's [ImageBuildOptions].
// See the full version in case we want to extend this in the future.
//
// [ImageBuildOptions]: https://pkg.go.dev/github.com/docker/docker/api/types/build#ImageBuildOptions
type ImageBuildOptions struct {
	Tags       []string
	Dockerfile string
}

func (o ImageBuildOptions) URL() string {
	query := url.Values{}
	if len(o.Tags) > 0 {
		query["t"] = o.Tags
	}
	if o.Dockerfile != "" {
		query.Set("dockerfile", o.Dockerfile)
	}
	u := url.URL{Scheme: "http", Host: "localhost", Path: "/build", RawQuery: query.Encode()}
	return u.String()
}
