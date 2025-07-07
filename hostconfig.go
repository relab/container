package container

import "github.com/relab/container/mount"

// The struct definitions in this file is largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [HostConfig].
//
// [HostConfig]: https://github.com/moby/moby/blob/master/api/types/container/hostconfig.go#L424

// HostConfig the non-portable Config structure of a container.
//
// This is a simplified version of the Docker API's [HostConfig].
// See the full version in case we want to extend this in the future.
//
// [HostConfig]: https://pkg.go.dev/github.com/docker/docker/api/types/container#HostConfig
type HostConfig struct {
	// Applicable to all platforms
	PortBindings PortMap // Port mapping between the exposed port (container) and the host
	AutoRemove   bool    // Automatically remove container when it exits

	// Mounts specs used by the container
	Mounts []mount.Mount `json:",omitempty"`
}
