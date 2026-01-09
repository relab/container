package container

import "github.com/relab/container/network"

// The struct definitions in this file are largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the originals [CreateRequest] and [CreateResponse].
//
// [CreateRequest]: https://github.com/moby/moby/blob/master/api/types/container/create_request.go
// [CreateResponse]: https://github.com/moby/moby/blob/master/api/types/container/create_response.go

// CreateRequest is the request message sent to the server for container
// create calls. It is a config wrapper that holds the container [Config]
// (portable) and the corresponding [HostConfig] (non-portable) and
// [network.NetworkingConfig].
//
// See the original [CreateRequest].
//
// [CreateRequest]: https://pkg.go.dev/github.com/docker/docker/api/types/container#CreateRequest
type CreateRequest struct {
	*Config
	HostConfig       *HostConfig               `json:"HostConfig,omitempty"`
	NetworkingConfig *network.NetworkingConfig `json:"NetworkingConfig,omitempty"`
}

// CreateResponse is the response message returned from the server.
//
// OK response to ContainerCreate operation
// swagger:model CreateResponse
type CreateResponse struct {
	// The ID of the created container
	// Required: true
	ID string `json:"Id"`

	// Warnings encountered when creating the container
	// Required: true
	Warnings []string `json:"Warnings"`
}
