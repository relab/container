package container

// The struct definitions in this file is largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [Config].
//
// [Config]: https://github.com/moby/moby/blob/master/api/types/container/config.go

// Config contains the configuration data about a container.
// It should hold only portable information about the container.
// Here, "portable" means "independent from the host we are running on".
// Non-portable information *should* appear in HostConfig.
//
// This is a simplified version of the Docker API's [Config].
// See the full version in case we want to extend this in the future.
//
// [Config]: https://pkg.go.dev/github.com/docker/docker/api/types/container#Config
type Config struct {
	User         string   // User that will run the command(s) inside the container, also support user:group
	ExposedPorts PortSet  `json:",omitempty"` // List of exposed ports
	Env          []string // List of environment variable to set in the container
	Cmd          []string // Command to run when starting the container
	Image        string   // Name of the image as it was passed by the operator (e.g. could be symbolic)
}

// PortBinding represents a binding between a Host IP address and a Host Port
type PortBinding struct {
	// HostIP is the host IP Address
	HostIP string `json:"HostIp"`
	// HostPort is the host port number
	HostPort string
}

// PortMap is a collection of PortBinding indexed by Port
type PortMap map[Port][]PortBinding

// PortSet is a collection of structs indexed by Port
type PortSet map[Port]struct{}

// Port is a string containing port number and protocol in the format "80/tcp"
type Port string
