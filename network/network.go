package network

// CreateOptions holds options to create a network.
//
// This is a simplified version of the Docker API's [CreateOptions].
// See the full version in case we want to extend this in the future.
//
// [CreateOptions]: https://pkg.go.dev/github.com/docker/docker/api/types/network#CreateOptions
type CreateOptions struct {
	Name   string // Name is the requested name of the network.
	Driver string // Driver is the driver-name used to create the network (e.g. `bridge`, `overlay`)
}

// ConnectOptions represents the data to be used to connect a container to the
// network.
//
// This is a simplified version of the Docker API's [ConnectOptions].
// See the full version in case we want to extend this in the future.
//
// [ConnectOptions]: https://pkg.go.dev/github.com/docker/docker/api/types/network#ConnectOptions
type ConnectOptions struct {
	Container      string
	EndpointConfig *EndpointSettings `json:",omitempty"`
}

// DisconnectOptions represents the data to be used to disconnect a container
// from the network.
type DisconnectOptions struct {
	Container string
	Force     bool
}

// CreateResponse NetworkCreateResponse
//
// OK response to NetworkCreate operation
// swagger:model CreateResponse
type CreateResponse struct {
	// The ID of the created network.
	// Required: true
	ID string `json:"Id"`

	// Warnings encountered when creating the container
	// Required: true
	Warning string `json:"Warning"`
}

// NetworkingConfig represents the container's networking configuration for each of its interfaces
// Carries the networking configs specified in the `docker run` and `docker network connect` commands
type NetworkingConfig struct {
	EndpointsConfig map[string]*EndpointSettings // Endpoint configs for each connecting network
}
