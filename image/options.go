package image

import "net/url"

// The struct definitions in this file are largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [RemoveOptions].
//
// [RemoveOptions]: https://github.com/moby/moby/blob/master/api/types/image/opts.go#L87

// RemoveOptions holds parameters to remove images.
type RemoveOptions struct {
	Force         bool
	PruneChildren bool
}

func (o RemoveOptions) URL(imageID string) string {
	query := url.Values{}
	if o.Force {
		query.Set("force", "1")
	}
	if o.PruneChildren {
		query.Set("prune_children", "1")
	}
	u := url.URL{Scheme: "http", Host: "localhost", Path: "/images/" + imageID, RawQuery: query.Encode()}
	return u.String()
}
