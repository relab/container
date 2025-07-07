package mount

// The struct definitions in this file are largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [Mount].
//
// [Mount]: https://github.com/moby/moby/blob/master/api/types/mount/mount.go

// Type represents the type of a mount.
type Type string

// Type constants
const (
	// TypeBind is the type for mounting host dir
	TypeBind Type = "bind"
	// TypeVolume is the type for remote storage volumes
	TypeVolume Type = "volume"
	// TypeTmpfs is the type for mounting tmpfs
	TypeTmpfs Type = "tmpfs"
	// TypeNamedPipe is the type for mounting Windows named pipes
	TypeNamedPipe Type = "npipe"
	// TypeCluster is the type for Swarm Cluster Volumes.
	TypeCluster Type = "cluster"
	// TypeImage is the type for mounting another image's filesystem
	TypeImage Type = "image"
)

// Mount represents a mount (volume).
type Mount struct {
	Type Type `json:",omitempty"`
	// Source specifies the name of the mount. Depending on mount type, this
	// may be a volume name or a host path, or even ignored.
	// Source is not supported for tmpfs (must be an empty value)
	Source string `json:",omitempty"`
	Target string `json:",omitempty"`
}
