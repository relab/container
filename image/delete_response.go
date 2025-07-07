package image

// The struct definitions in this file are largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [DeleteResponse].
//
// [DeleteResponse]: https://github.com/moby/moby/blob/master/api/types/image/delete_response.go

// DeleteResponse delete response
// swagger:model DeleteResponse
type DeleteResponse struct {
	// The image ID of an image that was deleted
	Deleted string `json:"Deleted,omitempty"`

	// The image ID of an image that was untagged
	Untagged string `json:"Untagged,omitempty"`
}
